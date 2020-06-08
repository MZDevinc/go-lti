package lti

import "net/http"

// GetRequesterValues gets the issuer and client_id of a request before parsing it
func GetRequesterValues(req *http.Request) (string, string) {
	issuer := req.FormValue("iss")
	clientID := req.FormValue("client_id")

	return issuer, clientID
}
