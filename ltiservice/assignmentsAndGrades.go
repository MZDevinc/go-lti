package ltiservice

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/MZDevinc/go-lti/lti"
	"github.com/pkg/errors"
)

// AGService An instance of an Assignment and Grade services connection
type AGService struct {
	ltis *LTIService
	// Scopes provided by the launch message from which the AGS is created
	Scopes       []string
	LineItemURL  *string
	LineItemsURL *string
}

// GetAGService get an upgraded service object that can handle server-to-server calls
// based on the Assignments and Grade Services specification
// This is a separate service because it is optionally supported and requires additional configuration
func (ltis *LTIService) GetAGService(msg lti.LaunchMessage) (*AGService, error) {
	ags := AGService{
		ltis: ltis,
	}

	if msg.Endpoint == nil {
		return nil, fmt.Errorf("Message had no assignment and grade services endpoint")
	}

	ags.Scopes = msg.Endpoint.Scope
	if msg.Endpoint.LineItem != "" {
		ags.LineItemURL = &msg.Endpoint.LineItem
	}
	if msg.Endpoint.LineItems != "" {
		ags.LineItemsURL = &msg.Endpoint.LineItems
	}

	return &ags, nil
}

// HasScope check whether the AGService has a given scope
func (ags *AGService) HasScope(scope string) bool {
	for _, str := range ags.Scopes {
		if str == scope {
			return true
		}
	}

	return false
}

// FindOrCreateLineItem returns an existing line item, or creates one if it doesn't exist
// Existing line items are matched based on the "tag" field, which must be unique among the line items in a context
// Even when specifying an existing LineItem, the platform can and will modify the results (primarily to add a platform
// ID for the line item, but also to correct attributes such as maximum score)
//
// Will fail if the LTIService underlying the AGService cannot provide a signing key, the AGService doesn't have
// sufficient information for the platform's line items URL, or doesn't have sufficient scope to access line items
func (ags *AGService) FindOrCreateLineItem(lineItem lti.LineItem) (lti.LineItem, error) {
	result := lti.LineItem{}

	ags.ltis.debug("findOrCreateLineItem: %+v", lineItem)
	inscope := ags.HasScope(lti.ScopeLineItem)
	if !inscope {
		fmt.Println("NOT IN SCOPE")
		return result, fmt.Errorf("missing necessary scope: %q", lti.ScopeLineItem)
	}

	if ags.LineItemsURL == nil {
		return result, fmt.Errorf("missing line item url")
	}

	// lineitemsURL := (*s.svcData)["lineitems"].(string)
	ags.ltis.debug("calling GET on lineitems url: %q", *ags.LineItemsURL)
	res, err := ags.ltis.DoServiceRequest(ags.Scopes, *ags.LineItemsURL, "", "", "", "application/vnd.ims.lis.v2.lineitemcontainer+json")
	if err != nil {
		return result, errors.Wrap(err, "Failure fetching existing lineitems")
	}
	log.Printf("lineitems initial lookup result: %+v", res)

	// find lineitem in existing list from provider,
	// if it exists, return it (tag should equal lineItem.Tag if it's the same)
	var existingLineitems []lti.LineItem
	if err := json.Unmarshal([]byte(res.Body), &existingLineitems); err != nil {
		return result, errors.Wrap(err, "Failed to process lineitems")
	}
	for _, li := range existingLineitems {
		if li.Tag == lineItem.Tag {
			log.Printf("Found lineitem amongst existing, returning: %+v", li)
			return li, nil
		}
	}

	// since we didn't find one, create it and return it
	bodyBytes, err := json.Marshal(lineItem)
	if err != nil {
		return result, fmt.Errorf("Failed to serialize lineitem for sending")
	}
	ags.ltis.debug("calling POST on lineitems url: %q with body: %q", *ags.LineItemsURL, string(bodyBytes))

	res, err = ags.ltis.DoServiceRequest(ags.Scopes, *ags.LineItemsURL, "POST", string(bodyBytes), "application/vnd.ims.lis.v2.lineitem+json", "application/vnd.ims.lis.v2.lineitem+json")
	if err != nil {
		return result, errors.Wrap(err, "Failed to create new line item")
	}
	ags.ltis.debug("result from lineitem post: %+v", res)
	if err := json.Unmarshal([]byte(res.Body), &result); err != nil {
		log.Printf("Error during unmarshall: %+v", err)
		return result, errors.Wrap(err, "Failed to parse new line item")
	}
	return result, nil
}

// PutGrade saves a score to the LTI platform for the given line item
func (ags *AGService) PutGrade(lineItem lti.LineItem, grade lti.Grade) error {
	inscope := ags.HasScope(lti.ScopeScore)
	if !inscope {
		return fmt.Errorf("missing necessary scope: %q", lti.ScopeScore)
	}

	if lineItem.ID == "" {
		return fmt.Errorf("line item is missing id/endpoint")
	}

	scoreURL := fmt.Sprintf("%s/scores", lineItem.ID)
	ags.ltis.debug("Final score url: %s", scoreURL)

	jsonBodyBytes, err := json.Marshal(grade)
	if err != nil {
		return errors.Wrap(err, "Failed to encode JSON")
	}

	res, err := ags.ltis.DoServiceRequest(ags.Scopes, scoreURL, "POST", string(jsonBodyBytes), "application/vnd.ims.lis.v1.score+json", "")
	if err != nil {
		return errors.Wrap(err, "Failed to put grade")
	}
	log.Printf("put grades service request result: %+v", res)

	return nil
}
