package ltiservice

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/pkg/errors"
)

// LTIService An instance of an LTI connection
type LTIService struct {
	Store          sessions.Store
	Config         Config
	routes         []routeDef
	SigningKeyFunc *func() (jwa.SignatureAlgorithm, interface{}, error)
	OutgoingJWTkid string
	debug          func(string, ...interface{})
}

// Config configuration for the Platform/Tool interface
type Config struct {
	AuthLoginURL string // URL on the Platform that handles the Login redirect
	LaunchURL    string // URL on the Tool that handles the Launch request
	ClientID     string // The Platform's Client ID
	KeySetURL    string // URL on the Platform that provides its public keys via JWKS
	AuthTokenURL string // URL to obtain an auth token
	Issuer       string // Issuer URL, for creating initial JWT
}

// NewLTIService Returns an LTIService initialized with given configuration and stores
func NewLTIService(store sessions.Store, config Config) *LTIService {
	debug := func(format string, a ...interface{}) {
		// Production mode, no-op
	}

	return &LTIService{Store: store, Config: config, debug: debug}
}

// NewLTIServiceWithDebug Returns an LTIService initialized with given configuration and stores,
// which will output debug messages using log.Printf
func NewLTIServiceWithDebug(store sessions.Store, config Config) *LTIService {
	debug := func(format string, a ...interface{}) {
		log.Printf(format, a)
	}

	return &LTIService{Store: store, Config: config, debug: debug}
}

// NewLTIServiceWithCustomDebug Returns an LTIService initialized with given configuration and stores,
// which will output debug messages using the given debug handler
func NewLTIServiceWithCustomDebug(store sessions.Store, config Config, debug func(string, ...interface{})) *LTIService {
	return &LTIService{Store: store, Config: config, debug: debug}
}

// SetSigningKeyFunc Define a function that can be used to get a signing key for JWTs
func (ltis *LTIService) SetSigningKeyFunc(handler func() (jwa.SignatureAlgorithm, interface{}, error)) {
	ltis.SigningKeyFunc = &handler
}

// getValidationKey fetches the public key used to validate a JWT token from the platform
// Currently always pulls from external URL (as defined in the service's config object) with no cache
func (ltis *LTIService) getValidationKey(token *jwt.Token) (interface{}, error) {

	//TO-DO
	//Maybe try to get this from cache first

	keyset, err := jwk.Fetch(ltis.Config.KeySetURL)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed fetching keyset from endpoint: %q", ltis.Config.KeySetURL)
	}

	//TO-DO
	//Now that we have the keyset, save it in cache if we implement a cache

	kid, ok := token.Header["kid"]
	if !ok {
		return nil, errors.Wrapf(err, "Failed fetching keyset from endpoint: %q", "no kid")
	}
	kidStr := kid.(string)
	ltis.debug("Looking for token kid: %q", kidStr)

	keys := keyset.LookupKeyID(kidStr)
	if keys == nil || len(keys) < 1 {
		return nil, fmt.Errorf("Token validation key not found for kid: %q", kidStr)
	} else if len(keys) > 1 {
		log.Printf("Multiple validation keys found for kid value: %q (using first one)", kidStr)
	}

	// materializedKey, err := keys[0].Materialize()
	// if err != nil {
	// 	return nil, err
	// }

	ltis.debug("Returning parsed key: %+v\n", keys[0])

	var key interface{}
	if err := keys[0].Raw(&key); err != nil {
		log.Printf("failed to create public key: %s", err)
		return key, err
	}

	return key, nil
}

func (ltis *LTIService) getSigningKey() (jwa.SignatureAlgorithm, interface{}, error) {
	if ltis.SigningKeyFunc == nil {
		return "", nil, fmt.Errorf("No key available")
	}

	handler := *ltis.SigningKeyFunc
	return handler()
}

// ServiceResult is a holder object for the results of a service call
type ServiceResult struct {
	Header http.Header
	Body   string
}

// DoServiceRequest fetches an auth token for a service call, then makes and returns the results of that call
func (ltis *LTIService) DoServiceRequest(scopes []string, url, pMethod, body, pContentType, pAccept string) (*ServiceResult, error) {
	var (
		method      = "GET"
		contentType = "application/json"
		accept      = "application/json"
		req         *http.Request
	)
	if pMethod != "" {
		method = pMethod
	}
	if pContentType != "" {
		contentType = pContentType
	}
	if pAccept != "" {
		accept = pAccept
	}
	accessToken, err := ltis.GetAccessToken(scopes)
	if err != nil {
		return nil, err
	}
	// log.Printf("access token fetched: %s", accessToken)
	client := &http.Client{Timeout: time.Second * 30}
	if method == "POST" || method == "PUT"{
		req, err = http.NewRequest(method, url, strings.NewReader(body))
		if err != nil {
			return nil, errors.Wrapf(err, "DoServiceReq: Error Creating new request for POST to %q", url)
		}
		req.Header.Add("Content-Type", contentType)
	} else { // GET
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, errors.Wrapf(err, "DoServiceReq: Error Creating new request for method: %q to %q", method, url)
		}
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Add("Accept", accept)

	ltis.debug("About to make request for url, request: %+v", req)
	ltis.debug("Request body: %s", body)

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "DoServiceReq: Error Executing new request for method: %q to %q", method, url)
	}

	ltis.debug("Response received for method: %q to %q: %q", method, url, resp.Status)

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "DoServiceReq: Error reading the response body for method: %q to %q", method, url)
	}

	return &ServiceResult{Header: resp.Header, Body: string(bodyBytes)}, nil
}
