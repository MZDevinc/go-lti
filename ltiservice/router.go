package ltiservice

import (
	"regexp"
	"strings"

	"github.com/MZDevinc/go-lti/lti"
	"github.com/gin-gonic/gin"
)

type routeDef struct {
	Path      string
	MatchExp  *regexp.Regexp
	HasParams bool
	Handler   func(*gin.Context, map[string]string)
}

// DefineRoute Match URL path pattern with a handler function
// The path can define parameters similarly to gin (/:parameter/), and the handler will receive those parameter in a map
// when the route is matched
func (ltis *LTIService) DefineRoute(path string, handler func(c *gin.Context, params map[string]string)) error {
	found := false
	for i := range ltis.routes {
		if ltis.routes[i].Path == path {
			ltis.routes[i].Handler = handler
			found = true
		}
	}

	if !found {
		matchExp, hasParams, err := ltis.processRoute(path)
		if err != nil {
			return err
		}
		ltis.routes = append(ltis.routes, routeDef{
			Path:      path,
			MatchExp:  matchExp,
			HasParams: hasParams,
			Handler:   handler,
		})
	}

	return nil
}

// ParseRoute given a definite URL path, tries to match that path with routes that have been defined
func (ltis *LTIService) ParseRoute(c *gin.Context, path string, msg lti.LaunchMessage) {
	ltis.debug("Route path", path)
	for i := range ltis.routes {
		route := ltis.routes[i]
		if route.HasParams {
			match := route.MatchExp.FindStringSubmatch(path)
			if match != nil {
				ltis.debug("Matched route", route.Path)
				paramNames := route.MatchExp.SubexpNames()
				params := make(map[string]string)
				for j := range paramNames {
					params[paramNames[j]] = match[j]
				}

				route.Handler(c, params)
				return
			}
		} else {
			match := route.MatchExp.FindStringIndex(path)
			if match != nil {
				ltis.debug("Matched route", route.Path)
				route.Handler(c, nil)
				return
			}
		}
	}
}

func (ltis *LTIService) processRoute(path string) (*regexp.Regexp, bool, error) {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	parts := strings.Split(path, "/")

	hasParams := false
	pattern := ""
	for i := 0; i < len(parts); i++ {
		part := parts[i]
		if part == "" {
			continue
		}

		if part[0] == ':' {
			hasParams = true
			part = strings.TrimPrefix(part, ":")
			part = regexp.QuoteMeta(part)
			pattern += `\/(?P<` + part + `>[^\/]+)`
		} else {
			part = regexp.QuoteMeta(part)
			pattern += `\/` + part
		}
	}
	pattern += `\/?$`

	re, err := regexp.Compile(pattern)
	return re, hasParams, err
}
