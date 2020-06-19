package ltiservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// GetAccessToken Create a JWT that will request an oauth access token from the platform, send that to the registered
// AuthTokenURL, and then return the response token.
func (ltis *LTIService) GetAccessToken(scopes []string) (string, error) {
	sort.Strings(scopes)
	scopeStr := strings.Join(scopes, " ")

	method, privkey, err := ltis.getSigningKey()
	if method != jwa.RS256 {
		return "", errors.Wrapf(err, "GetAccessToken: Wrong Tool Private Key type for clientID: %q.", ltis.Config.ClientID)
	}
	if err != nil {
		return "", errors.Wrapf(err, "GetAccessToken: Error getting Tool Private Key for clientID: %q.", ltis.Config.ClientID)
	}

	timestamp := int(time.Now().Unix())
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": ltis.Config.Issuer,
		"sub": ltis.Config.ClientID,
		"aud": ltis.Config.AuthTokenURL,
		"iat": timestamp,
		"exp": timestamp + 60,
		"jti": fmt.Sprintf("lti-service-token-%s", uuid.NewV4().String()),
	})
	tokenStr, err := token.SignedString(privkey)
	if err != nil {
		return "", errors.Wrapf(err, "GetAccessToken: Error signing token for clientId: %q.", ltis.Config.ClientID)
	}
	ltis.debug("GetAccessToken generated JWT: %s", tokenStr)

	client := &http.Client{Timeout: time.Second * 30}
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	form.Add("client_assertion", tokenStr)
	form.Add("scope", scopeStr)

	ltis.debug("GetAccessToken fetch URL: %s", ltis.Config.AuthTokenURL)
	ltis.debug("GetAccessToken fetch parameters: %s", form.Encode())

	req, err := http.NewRequest("POST", ltis.Config.AuthTokenURL, strings.NewReader(form.Encode()))

	if err != nil {
		return "", errors.Wrapf(err, "GetAccessToken: Error generating the token request url for clientId: %q.", ltis.Config.ClientID)
	}
	response, err := client.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, "GetAccessToken: Error executing the form POST for clientId: %q.", ltis.Config.ClientID)
	}
	log.Printf("Access token response status: %s", response.Status)

	defer response.Body.Close()
	// log.Printf("returned headers: %v", response.Header)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.Wrapf(err, "GetAccessToken: Error reading body of access token fetch response for clientId: %q.", ltis.Config.ClientID)
	}
	ltis.debug("GetAccessToken response body: %s", body)

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return "", fmt.Errorf("GetAccessToken: Error response from access token fetch (%q)", response.Status)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", errors.Wrapf(err, "GetAccessToken: Failed to parse json from body of access token fetch response for clientId: %q.", ltis.Config.ClientID)
	}
	accessToken := data["access_token"].(string)
	ltis.debug("GatAccessToken received access token: %s", accessToken)

	return accessToken, nil
}
