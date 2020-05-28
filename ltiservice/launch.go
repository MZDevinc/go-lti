package ltiservice

import (
	"fmt"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
)

// The name of the property in the request where the user information
// from the JWT will be stored.
const userProperty = "user"

//GetLaunchHandler Returns a handler for a LaunchMessage
//Once the incoming JWT is decoded and validated, the provided callback function will
//be executed
func (ltis *LTIService) GetLaunchHandler(callback func(claims jwt.MapClaims)) http.Handler {
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ltis.launch(w, req, callback)
	})

	//Wraps the handler with middleware that decodes the incoming JWT
	opts := jwtmiddleware.Options{
		SigningMethod: jwt.SigningMethodRS256,
		UserProperty:  userProperty,
		Extractor: func(r *http.Request) (string, error) {
			return r.FormValue("id_token"), nil
		},
		Debug:               true,
		ValidationKeyGetter: ltis.getValidationKey,
		ErrorHandler:        tokenMWErrorHandler,
	}
	jwtMW := jwtmiddleware.New(opts)

	return jwtMW.Handler(handlerFunc)
}

func (ltis *LTIService) launch(w http.ResponseWriter, req *http.Request, callback func(jwt.MapClaims)) {
	//Extract claims from the JWT
	userToken := req.Context().Value(userProperty)
	tok := userToken.(*jwt.Token)
	claims := tok.Claims.(jwt.MapClaims)
	fmt.Println("CLAIMS")
	fmt.Println(claims)

	callback(claims)
}

//tokenMWErrorHandler provided to the JWT middleware for it to handle errors
func tokenMWErrorHandler(w http.ResponseWriter, r *http.Request, err string) {
	fmt.Println("JWT ERROR HANDLER")
	http.Error(w, fmt.Sprintf("Token issue: %s", err), 401)
}
