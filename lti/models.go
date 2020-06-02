package lti

// This package contains models for objects as defined in the LTI 1.3 specification

// DefinesUserIdentity represents that an object contains fields to identify a user, and provides methods for accessing them
type DefinesUserIdentity interface {
	IsAnonymous() bool
	GetSub() *string
	GetGivenName() *string
	GetFamilyName() *string
	GetName() *string
	GetEmail() *string
}

// DefinesContext represents that an object contains a Context object, and provides a method for accessing it
type DefinesContext interface {
	GetContext() *Context
}

// ResourceLinkMessage represents an LTI Resource Link Launch Request
// See http://www.imsglobal.org/spec/lti/v1p3/#launch-from-a-resource-link-0
// Many claims are optional. Optional claims are represented by pointer fields.
type ResourceLinkMessage struct {
	// MessageType guaranteed to be "LtiResourceLinkRequest" for this message type
	MessageType   string `json:"https://purl.imsglobal.org/spec/lti/claim/message_type"`
	Version       string `json:"https://purl.imsglobal.org/spec/lti/claim/version"`
	DeploymentID  string `json:"https://purl.imsglobal.org/spec/lti/claim/deployment_id"`
	TargetLinkURI string `json:"https://purl.imsglobal.org/spec/lti/claim/target_link_uri"`

	ResourceLink ResourceLink `json:"https://purl.imsglobal.org/spec/lti/claim/resource_link"`

	// User identity claims
	// If Sub is not defined, the request is anonymous
	Sub        *string `json:"sub"`
	GivenName  *string `json:"given_name"`
	FamilyName *string `json:"family_name"`
	Name       *string `json:"name"`
	Email      *string `json:"email"`

	// Roles may be empty
	Roles []string `json:"https://purl.imsglobal.org/spec/lti/claim/roles"`

	// Context represents the context, or class/course, from which the launch occurred
	Context *Context `json:"https://purl.imsglobal.org/spec/lti/claim/context"`

	ToolPlatform *ToolPlatform `json:"https://purl.imsglobal.org/spec/lti/claim/tool_platform"`

	// RoleScopeMentor contains an array of the user ID values which the current, launching user can access as a mentor.
	// The sender of the message MUST NOT include a list of user ID values in this property unless they also provide
	// http://purl.imsglobal.orb/vocab/lis/v2/membership#Mentor as one of the values passed in the roles claim.
	RoleScopeMentor *[]string `json:"https://purlimsglobal.org/spec/lti/claim/role_scope_mentor"`

	// LaunchPresentation contains contextual information about how the launch will be displayed in the platform
	LaunchPresentation *LaunchPresentation `json:"https://purl.imsglobal.org/spec/lti/claim/launch_presentation"`

	// LIS contains additional information about Learning Information Services software associations
	LIS *LIS `json:"https://purl.imsglobal.org/spec/lti/claim/lis"`

	// Additional custom properties
	// See http://www.imsglobal.org/spec/lti/v1p3/#custom-variables-0
	Custom *map[string]string

	// Additional custom properties in the root of the message
	// "Vendors MAY extend the information model for any message type and inject additional properties into the
	// message's JSON object by adding one or more claims. Vendors MUST use a fully-qualified URL as the claim name for
	// any of their extension claims.
	// By best practice, vendors should define custom variables... instead of relying on extension properties."
	Extensions *map[string]string
}

// ResourceLink composes properties for the resource link from which the launch message occurs
type ResourceLink struct {
	ID          string  `json:"id"`
	Description *string `json:"description"`
	Title       *string `json:"title"`
}

// Context represents the context, or class/course, from which the launch occurred
type Context struct {
	ID    string  `json:"id"`
	Type  *string `json:"type"`
	Label *string `json:"label"`
	Title *string `json:"title"`
}

// ToolPlatform represents information about the platform from which the launch occurred
type ToolPlatform struct {
	GUID              string  `json:"guid"`
	ContactEmail      *string `json:"contact_email"`
	Description       *string `json:"description"`
	Name              *string `json:"name"`
	URL               *string `json:"url"`
	ProductFamilyCode *string `json:"product_family_code"`
	Version           *string `json:"version"`
}

// LaunchPresentation contains contextual information about how the launch will be displayed in the platform
type LaunchPresentation struct {
	DocumentTarget *string `json:"document_target"`
	Height         *string `json:"height"`
	Width          *string `json:"width"`
	ReturnURL      *string `json:"return_url"`
	Locale         *string `json:"locale"`
}

// LIS contains additional information about Learning Information Services software associations
type LIS struct {
	CourseOfferingSourcedID *string `json:"course_offering_sourcedid"`
	CourseSectionSourcedID  *string `json:"course_section_sourcedid"`
	OutcomeServiceURL       *string `json:"outcome_service_url"`
	PersonSourcedID         *string `json:"person_sourcedid"`
	ResultSourcedID         *string `json:"result_sourcedid"`
}

const (
	//Core institution roles
	InstitutionRoleAdministrator = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Administrator"
	InstitutionRoleFaculty       = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Faculty"
	InstitutionRoleGuest         = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Guest"
	InstitutionRoleNone          = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#None"
	InstitutionRoleOther         = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Other"
	InstitutionRoleStaff         = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Staff"
	InstitutionRoleStudent       = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Student"

	//Non‑core institution roles
	InstitutionRoleAlumni             = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Alumni"
	InstitutionRoleInstructor         = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Instructor"
	InstitutionRoleLearner            = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Learner"
	InstitutionRoleMember             = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Member"
	InstitutionRoleMentor             = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Mentor"
	InstitutionRoleObserver           = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Observer"
	InstitutionRoleProspectiveStudent = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#ProspectiveStudent"

	//Core context roles
	ContextRoleAdministrator    = "http://purl.imsglobal.org/vocab/lis/v2/membership#Administrator"
	ContextRoleContentDeveloper = "http://purl.imsglobal.org/vocab/lis/v2/membership#ContentDeveloper"
	ContextRoleInstructor       = "http://purl.imsglobal.org/vocab/lis/v2/membership#Instructor"
	ContextRoleLearner          = "http://purl.imsglobal.org/vocab/lis/v2/membership#Learner"
	ContextRoleMentor           = "http://purl.imsglobal.org/vocab/lis/v2/membership#Mentor"

	//Non‑core context roles
	ContextRoleManager = "http://purl.imsglobal.org/vocab/lis/v2/membership#Manager"
	ContextRoleMember  = "http://purl.imsglobal.org/vocab/lis/v2/membership#Member"
	ContextRoleOfficer = "http://purl.imsglobal.org/vocab/lis/v2/membership#Officer"
)

//HasRole check if the message includes the given role
func (rlm ResourceLinkMessage) HasRole(role string) bool {
	for i := range rlm.Roles {
		if rlm.Roles[i] == role {
			return true
		}
	}

	return false
}

//HasAnyRole check if the message includes any of a given list of roles
func (rlm ResourceLinkMessage) HasAnyRole(roles []string) bool {
	for i := range rlm.Roles {
		for j := range roles {
			if rlm.Roles[i] == roles[j] {
				return true
			}
		}
	}

	return false
}

//IsAnonymous is the launch request anonymous
func (rlm ResourceLinkMessage) IsAnonymous() bool {
	return rlm.Sub == nil
}

//GetSub returns claim "sub"
//Not required, may be nil (request is anonymous)
func (rlm ResourceLinkMessage) GetSub() *string {
	return rlm.Sub
}

//GetGivenName returns claim "given_name"
//Not required, may be nil
func (rlm ResourceLinkMessage) GetGivenName() *string {
	return rlm.GivenName
}

//GetFamilyName returns claim "family_name"
//Not required, may be nil
func (rlm ResourceLinkMessage) GetFamilyName() *string {
	return rlm.FamilyName
}

//GetName returns claim "name" (full display name)
//Not required, may be nil
func (rlm ResourceLinkMessage) GetName() *string {
	return rlm.Name
}

//GetEmail returns claim "email"
//Not required, may be nil
func (rlm ResourceLinkMessage) GetEmail() *string {
	return rlm.Email
}

//GetContext returns claim "https://purl.imsglobal.org/spec/lti/claim/context"
//Not required, may be nil
func (rlm ResourceLinkMessage) GetContext() *Context {
	return rlm.Context
}
