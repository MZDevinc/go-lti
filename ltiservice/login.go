package ltiservice

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

//GetLoginHandler Returns a handler for the Login (first) step of the OIDC process
//This receives the Platform's request to log in, and redirects it back to the Platform
//after adding appropriate auth info
func (ltis *LTIService) GetLoginHandler() http.Handler {
	handlerFunc := http.HandlerFunc(ltis.login)
	return handlerFunc
}

func (ltis *LTIService) login(w http.ResponseWriter, req *http.Request) {
	//We're still not really sure what the purpose of storing the session is
	//sess, _ := O.store.Get(req, O.sessionName)

	if ltis.Config.LaunchURL == "" {
		http.Error(w, "launch url is not configured.", 400)
		return
	}

	err := validateLogin(req)
	if err != nil {
		http.Error(w, errors.Wrap(err, "Oidc Login validation failure").Error(), 400)
		return
	}

	state := fmt.Sprintf("state-%s", uuid.NewV4().String())
	setStateCookie(w, state)

	nonce := fmt.Sprintf("nonce-%s", uuid.NewV4().String())
	//For now we will ignore checking that the nonce matches and has been used only once
	//O.cache.PutNonce(req, nonce)

	//Construct a query for the redirect
	redirReq, err := http.NewRequest("GET", ltis.Config.AuthLoginURL, nil)
	if err != nil {
		http.Error(w, errors.Wrap(err, "Failed to construct redirect").Error(), 500)
		return
	}
	q := redirReq.URL.Query()
	q.Add("scope", "openid")
	q.Add("response_type", "id_token")
	q.Add("response_mode", "form_post")
	q.Add("prompt", "none")
	q.Add("client_id", ltis.Config.ClientID)
	q.Add("redirect_uri", ltis.Config.LaunchURL)
	q.Add("state", state)
	q.Add("nonce", nonce)
	q.Add("login_hint", req.FormValue("login_hint"))
	if mh := req.FormValue("lti_message_hint"); mh != "" {
		q.Add("lti_message_hint", mh)
	}
	redirReq.URL.RawQuery = q.Encode()
	redirURL := redirReq.URL.String()
	log.Printf("OIDC Login Redir: %s", redirURL)

	//We're still not really sure what the purpose of storing the session is
	//sess.Save(req, w)
	http.Redirect(w, req, redirURL, 302)
}

func validateLogin(req *http.Request) error {
	iss := req.FormValue("iss")
	// log.Printf("iss: %q", iss)
	if iss == "" {
		return fmt.Errorf("issuer not found")
	}
	loginHint := req.FormValue("login_hint")
	// log.Printf("login_hint: %q", loginHint)
	if loginHint == "" {
		return fmt.Errorf("login hint not found")
	}
	return nil
}

func setStateCookie(w http.ResponseWriter, state string) {
	secs := 3600

	// Cookie with SameSite: none, for Chrome compatibility
	cookie := http.Cookie{
		Name:     fmt.Sprintf("mzdevinc_lti_go_%s", state),
		Value:    state,
		Expires:  time.Now().Add(time.Second * time.Duration(secs)),
		MaxAge:   secs,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, &cookie)

	// Cookie without SameSite: none, for iOS 12 compatibility
	cookie2 := http.Cookie{
		Name:     fmt.Sprintf("mzdevinc_lti_go2_%s", state),
		Value:    state,
		Expires:  time.Now().Add(time.Second * time.Duration(secs)),
		MaxAge:   secs,
		SameSite: http.SameSiteDefaultMode,
	}
	http.SetCookie(w, &cookie2)
}
