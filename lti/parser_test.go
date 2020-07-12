package lti

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	msg1 := LaunchMessage{
		Iss:           "present",
		Aud:           "present",
		Iat:           1000,
		Exp:           1000,
		Nonce:         "present",
		MessageType:   "LtiResourceRequest",
		Version:       "1.3",
		DeploymentID:  "present",
		TargetLinkURI: "present",
		ResourceLink: ResourceLink{
			ID:          "present",
			Description: "present",
		},
	}
	err1 := validateLaunchMessage(msg1)
	assert.NoError(t, err1)

	msg2 := LaunchMessage{
		Iss: "present",
		// Aud:        missing,
		Iat:           1000,
		Exp:           1000,
		Nonce:         "present",
		MessageType:   "LtiResourceRequest",
		Version:       "1.3",
		DeploymentID:  "present",
		TargetLinkURI: "present",
		ResourceLink: ResourceLink{
			ID:          "present",
			Description: "present",
		},
	}
	err2 := validateLaunchMessage(msg2)
	fmt.Println("err2", err2)
	assert.Error(t, err2)

	msg3 := LaunchMessage{
		Iss:           "present",
		Aud:           "present",
		Iat:           1000,
		Exp:           1000,
		Nonce:         "present",
		MessageType:   "LtiResourceRequest",
		Version:       "1.3",
		DeploymentID:  "present",
		TargetLinkURI: "present",
		// ResourceLink: missing
	}
	err3 := validateLaunchMessage(msg3)
	fmt.Println("err3", err3)
	assert.Error(t, err3)

	msg4 := LaunchMessage{
		Iss:           "present",
		Aud:           "present",
		Iat:           1000,
		Exp:           1000,
		Nonce:         "present",
		MessageType:   "LtiResourceRequest",
		Version:       "1.3",
		DeploymentID:  "present",
		TargetLinkURI: "present",
		ResourceLink: ResourceLink{
			// ID:       missing,
			Description: "present",
		},
	}
	err4 := validateLaunchMessage(msg4)
	fmt.Println("err4", err4)
	assert.Error(t, err4)
}
