package ltiservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/MZDevinc/go-lti/lti"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/pkg/errors"
)

// SendDeepLinkingResponse Sends specified content items in a deep linking response to a given launch message.
// Details about where and how the response should be sent are pulled from the launch message.
// Since the Deep Linking response must be sent from within the client, this function creates an HTML stub page with
// Javascript to perform the actual transmission; therefore, a gin context parameter is also required.
func (ltis *LTIService) SendDeepLinkingResponse(c *gin.Context, msg lti.LaunchMessage, items []lti.ContentItem) error {
	token, err := ltis.GetDeepLinkingResponseJWT(msg, items)
	if err != nil {
		return err
	}

	// Create HTML template to transmit response
	respHTML, err := ltis.createDeepLinkingResponseHTML(msg.DeepLinkingSettings.DeepLinkReturnURL, token)
	if err != nil {
		return err
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(respHTML))
	return nil
}

// GetDeepLinkingResponseJWT Takes the specified launch message and content items, forms them into a deep linking
// response, then marshals that struct into a JWT that is signed with a key retrieved from the SigningKeyFunc that is
// registered with the LTIService. If no SigningKeyFunc is defined, will return an error.
func (ltis *LTIService) GetDeepLinkingResponseJWT(msg lti.LaunchMessage, items []lti.ContentItem) (string, error) {
	timestamp := int(time.Now().Unix())

	resp := lti.DeepLinkingResponse{
		Iss:   msg.Aud,
		Aud:   msg.Iss,
		Iat:   timestamp,
		Exp:   timestamp + 3600,
		Nonce: msg.Nonce,

		MessageType:  "LtiDeepLinkingResponse",
		Version:      "1.3.0",
		DeploymentID: msg.DeploymentID,
		Data:         msg.DeepLinkingSettings.Data,

		ContentItems: items,
	}

	// Create response JWT out of response object
	token, err := ltis.createJWT(resp)
	ltis.debug("Deep linking response JWT:", string(token))

	return string(token), err
}

func (ltis *LTIService) createJWT(data interface{}) ([]byte, error) {
	method, key, err := ltis.getSigningKey()
	if err != nil {
		return nil, err
	}

	buf, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, `failed to marshal token`)
	}

	hdr := jws.NewHeaders()
	if hdr.Set(`alg`, method.String()) != nil {
		return nil, errors.Wrap(err, `failed to sign payload`)
	}
	if hdr.Set(`typ`, `JWT`) != nil {
		return nil, errors.Wrap(err, `failed to sign payload`)
	}
	if hdr.Set(`kid`, ltis.OutgoingJWTkid) != nil {
		return nil, errors.Wrap(err, `failed to sign payload`)
	}
	signed, err := jws.Sign(buf, method, key, jws.WithHeaders(hdr))
	if err != nil {
		return nil, errors.Wrap(err, `failed to sign payload`)
	}

	fmt.Println(method, key)
	return signed, nil
}

var respTemplate = `
<html>
	<head></head>
	<body>
		<script>
			const form = document.createElement('form');
			form.method = "post";
			form.action = {{.Path}};
		
			const hiddenField = document.createElement('input');
			hiddenField.type = 'hidden';
			hiddenField.name = 'id_token;
			hiddenField.value = {{.Token}};
	
			form.appendChild(hiddenField);
		
			document.body.appendChild(form);
			form.submit();
		</script>
	</body>
</html>
`

type deepLinkingResponseData struct {
	Path  string
	Token string
}

func (ltis *LTIService) createDeepLinkingResponseHTML(path string, token string) (string, error) {
	tmpl, err := template.New("deepLinkingResponseHTML").Parse(respTemplate)
	if err != nil {
		return "", err
	}

	dat := deepLinkingResponseData{
		Path:  path,
		Token: token,
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, dat)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
