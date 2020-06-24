package ltiservice

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/MZDevinc/go-lti/lti"
	"github.com/pkg/errors"
)

// NRPService An instance of a Names and Roles Provisioning service connection
type NRPService struct {
	ltis          *LTIService
	Scopes        []string
	MembersURL    string
	nextCallRegex *regexp.Regexp
}

// GetNRPService get an upgraded service object that can handle server-to-server calls based on the Names and Roles
// Provisioning Service specification
// This is a separate service because it is optionally supported and requires additional configuration
func (ltis *LTIService) GetNRPService(msg lti.LaunchMessage) (*NRPService, error) {
	nrps := NRPService{
		ltis: ltis,
	}

	if msg.NamesRoleService == nil {
		return nil, fmt.Errorf("Message had no names and roles provisioning service endpoint")
	}

	nrps.Scopes = []string{lti.ScopeContextMembershipReadonly}
	nrps.MembersURL = msg.NamesRoleService.ContextMembershipsURL
	nrps.nextCallRegex = regexp.MustCompile("^?<(.*)>; ?rel=\"next\"$")

	return &nrps, nil
}

// GetMembers uses the Message Launches context and auth token to return a list of users associated with this launch
func (nrps *NRPService) GetMembers() (*lti.MemberResponse, error) {
	ret := &lti.MemberResponse{}

	count := 0
	svcURL := nrps.MembersURL
	for svcURL != "" {
		count++
		res, err := nrps.ltis.DoServiceRequest(nrps.Scopes, svcURL, "GET", "", "", "")
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to fetch member fetch #%d", count)
		}
		nrps.ltis.debug("------ nrps (%s) iteration %d success, body len: %d --------------", svcURL, count, len(res.Body))
		nrps.ltis.debug("  body: %+v", res.Body)
		resp := &lti.MemberResponse{}
		if err := json.Unmarshal([]byte(res.Body), resp); err != nil {
			return nil, errors.Wrapf(err, "failed to parse json, fetch #%d", count)
		}

		fmt.Println("RECEIVED MEMBERS")
		fmt.Println(resp.Members)

		if count == 1 {
			ret.Context = resp.Context
			ret.ID = resp.ID
			ret.Members = []lti.Member{}
		}
		ret.Members = append(ret.Members, resp.Members...)
		svcURL = nrps.getNextPageURL(res)
	}
	return ret, nil
}

func (nrps *NRPService) getNextPageURL(res *ServiceResult) string {
	for k, v := range res.Header {
		hKey := strings.ToLower(k)
		if hKey == "link" {
			for _, link := range v {
				if res := nrps.nextCallRegex.FindSubmatch([]byte(link)); res != nil {
					nextURL := string(res[1])
					nrps.ltis.debug("Next Url determined: %v", nextURL)
					return nextURL
				}
			}
		}
	}
	return ""
}
