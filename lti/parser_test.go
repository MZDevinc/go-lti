package lti

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	r1 := ResourceLink{
		ID:          "present",
		Description: "present",
	}
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
		ResourceLink:  &r1,
	}
	err1 := validateLaunchMessage(msg1)
	assert.NoError(t, err1)

	r2 := ResourceLink{
		ID:          "present",
		Description: "present",
	}
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
		ResourceLink:  &r2,
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
	assert.NoError(t, err3)

	r4 := ResourceLink{
		// ID:       missing,
		Description: "present",
	}
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
		ResourceLink:  &r4,
	}
	err4 := validateLaunchMessage(msg4)
	fmt.Println("err4", err4)
	assert.Error(t, err4)
}
