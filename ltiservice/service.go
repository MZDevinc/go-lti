package ltiservice

import (
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/pkg/errors"
)

//LTIService An instance of an LTI connection
type LTIService struct {
	Store  sessions.Store
	Config Config
	debug  func(string, ...interface{})
}

//Config configuration for the Platform/Tool interface
type Config struct {
	AuthLoginURL string //URL on the Platform that handles the Login redirect
	LaunchURL    string //URL on the Tool that handles the Launch request
	ClientID     string //The Platform's Client ID
	KeySetURL    string //URL on the Platform that provides its public keys via JWKS
}

//NewLTIService Returns an LTIService initialized with given configuration and stores
func NewLTIService(store sessions.Store, config Config) *LTIService {
	debug := func(format string, a ...interface{}) {
		//Production mode, no-op
	}

	return &LTIService{Store: store, Config: config, debug: debug}
}

//NewLTIServiceWithDebug Returns an LTIService initialized with given configuration and stores,
//which will output debug messages using log.Printf
func NewLTIServiceWithDebug(store sessions.Store, config Config) *LTIService {
	debug := func(format string, a ...interface{}) {
		log.Printf(format, a)
	}

	return &LTIService{Store: store, Config: config, debug: debug}
}

//NewLTIServiceWithCustomDebug Returns an LTIService initialized with given configuration and stores,
//which will output debug messages using the given debug handler
func NewLTIServiceWithCustomDebug(store sessions.Store, config Config, debug func(string, ...interface{})) *LTIService {
	return &LTIService{Store: store, Config: config, debug: debug}
}

//getValidationKey fetches the public key used to validate the jwt token
//Currently always pulls from external URL with no cache
func (ltis *LTIService) getValidationKey(token *jwt.Token) (interface{}, error) {

	//TO-DO
	//Maybe try to get this from cache first

	keyset, err := jwk.Fetch(ltis.Config.KeySetURL)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed fetching keyset from endpoint: %q", ltis.Config.KeySetURL)
	}

	//TO-DO
	//Now that we have the keyset, save it in cache if we implement a cache

	kid := token.Header["kid"].(string)
	ltis.debug("Looking for token kid: %q", kid)

	keys := keyset.LookupKeyID(kid)
	if keys == nil || len(keys) < 1 {
		return nil, fmt.Errorf("Token validation key not found for kid: %q", kid)
	} else if len(keys) > 1 {
		log.Printf("Multiple validation keys found for kid value: %q (using first one)", kid)
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
