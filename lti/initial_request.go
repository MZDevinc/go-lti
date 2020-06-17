package lti

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// GetRequesterValuesFromForm gets the issuer and client_id of a request from POST form (login request)
func GetRequesterValuesFromForm(req *http.Request) (string, string) {
	issuer := req.FormValue("iss")
	clientID := req.FormValue("client_id")

	return issuer, clientID
}

// GetRequesterValuesFromJWT gets the issuer and client_id of a request from JWT (launch request)
func GetRequesterValuesFromJWT(req *http.Request) (string, string, error) {
	idToken := req.FormValue("id_token")

	parsed, _ := jwt.Parse(idToken, nil)
	if parsed == nil {
		return "", "", fmt.Errorf("Could not parse JWT")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", fmt.Errorf("Could not get JWT claims")
	}

	issuer := fmt.Sprintf("%s", claims["iss"])
	clientID := fmt.Sprintf("%s", claims["aud"])

	return issuer, clientID, nil
}
