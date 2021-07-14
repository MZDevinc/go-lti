package lti

import "time"

// This package contains models for objects as defined in the LTI 1.3 specification

// DefinesContext represents that an object contains a Context object, and provides a method for accessing it
type DefinesContext interface {
	GetContext() *Context
}

// LaunchMessage represents an LTI Resource Link Launch Request
// See http://www.imsglobal.org/spec/lti/v1p3/#launch-from-a-resource-link-0
// Many claims are optional. Optional claims are represented by pointer fields.
type LaunchMessage struct {
	Iss   string `json:"iss" required:"true"` // Issuer
	Aud   string `json:"aud" required:"true"` // Client ID
	Iat   int    `json:"iat" required:"true"` // Token created timestamp
	Exp   int    `json:"exp" required:"true"` // Token will expire timestamp
	Nonce string `json:"nonce" required:"true"`

	// Message type should be either "LtiResourceLinkRequest" or "LtiDeepLinkingRequest"
	MessageType   string `json:"https://purl.imsglobal.org/spec/lti/claim/message_type" required:"true"`
	Version       string `json:"https://purl.imsglobal.org/spec/lti/claim/version" required:"true"`
	DeploymentID  string `json:"https://purl.imsglobal.org/spec/lti/claim/deployment_id" required:"true"`
	TargetLinkURI string `json:"https://purl.imsglobal.org/spec/lti/claim/target_link_uri"`

	ResourceLink *ResourceLink `json:"https://purl.imsglobal.org/spec/lti/claim/resource_link"`

	// User identity claims
	// If Sub is not defined, the request is anonymous
	Sub        *string `json:"sub"`
	GivenName  *string `json:"given_name"`
	FamilyName *string `json:"family_name"`
	MiddleName *string `json:"middle_name"`
	Name       *string `json:"name"`
	Email      *string `json:"email"`
	Picture    *string `json:"picture"`
	Locale     *string `json:"locale"`

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

	// DeepLinkingSettings contains information about the deep linking callback
	// Only defined when message_type is "LtiDeepLinkingRequest"
	DeepLinkingSettings *DeepLinkingSettings `json:"https://purl.imsglobal.org/spec/lti-dl/claim/deep_linking_settings"`

	// Endpoint contains information about Assignment and Grade Services connected to this message/context
	Endpoint *AGSEndpoint `json:"https://purl.imsglobal.org/spec/lti-ags/claim/endpoint"`

	// NamesRoleService contains information about the Names and Roles Provisioning Service connected to this message/context
	NamesRoleService *NamesRoleService `json:"https://purl.imsglobal.org/spec/lti-nrps/claim/namesroleservice"`

	// Additional custom properties
	// See http://www.imsglobal.org/spec/lti/v1p3/#custom-variables-0
	Custom *map[string]interface{} `json:"https://purl.imsglobal.org/spec/lti/claim/custom"`

	// Additional custom properties in the root of the message
	// "Vendors MAY extend the information model for any message type and inject additional properties into the
	// message's JSON object by adding one or more claims. Vendors MUST use a fully-qualified URL as the claim name for
	// any of their extension claims.
	// By best practice, vendors should define custom variables... instead of relying on extension properties."
	Extensions *map[string]string
}

// ResourceLink composes properties for the resource link from which the launch message occurs
type ResourceLink struct {
	ID          string `json:"id" required:"true"`
	Description string `json:"description"`
	Title       string `json:"title"`
}

// Context represents the context, or class/course, from which the launch occurred
type Context struct {
	ID    string   `json:"id" required:"true"`
	Type  []string `json:"type"`
	Label string   `json:"label"`
	Title string   `json:"title"`
}

// ToolPlatform represents information about the platform from which the launch occurred
type ToolPlatform struct {
	GUID              interface{} `json:"guid" required:"true"`
	ContactEmail      string      `json:"contact_email"`
	Description       string      `json:"description"`
	Name              string      `json:"name"`
	URL               string      `json:"url"`
	ProductFamilyCode string      `json:"product_family_code"`
	Version           string      `json:"version"`
}

// LaunchPresentation contains contextual information about how the launch will be displayed in the platform
type LaunchPresentation struct {
	DocumentTarget string `json:"document_target"`
	Height         int    `json:"height,omitempty"`
	Width          int    `json:"width,omitempty"`
	ReturnURL      string `json:"return_url"`
	Locale         string `json:"locale"`
}

// LIS contains additional information about Learning Information Services software associations
type LIS struct {
	CourseOfferingSourcedID string `json:"course_offering_sourcedid"`
	CourseSectionSourcedID  string `json:"course_section_sourcedid"`
	OutcomeServiceURL       string `json:"outcome_service_url"`
	PersonSourcedID         string `json:"person_sourcedid"`
	ResultSourcedID         string `json:"result_sourcedid"`
}

// DeepLinkingSettings additional information for a Deep Linking request
type DeepLinkingSettings struct {
	DeepLinkReturnURL                 string   `json:"deep_link_return_url" required:"true"`
	AcceptTypes                       []string `json:"accept_types"`
	AcceptPresentationDocumentTargets []string `json:"accept_presentation_document_targets"`
	AcceptMediaTypes                  string   `json:"accept_media_types"`
	AcceptMultiple                    bool     `json:"accept_multiple"`
	AutoCreate                        bool     `json:"auto_create"`
	Title                             string   `json:"title"`
	Text                              string   `json:"text"`
	Data                              string   `json:"data"`
}

// AGSEndpoint information about the platform endpoint for Assignment and Grade Services
type AGSEndpoint struct {
	Scope     []string `json:"scope"`
	LineItem  string   `json:"lineitem"`
	LineItems string   `json:"lineitems"`
}

// NamesRoleService information about the platform endpoint for the Names and Roles Provisioning Service
type NamesRoleService struct {
	ContextMembershipsURL string `json:"context_memberships_url"`
}

const (
	// Core institution roles

	// InstitutionRoleAdministrator http://purl.imsglobal.org/vocab/lis/v2/institution/person#Administrator
	InstitutionRoleAdministrator = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Administrator"
	// InstitutionRoleFaculty http://purl.imsglobal.org/vocab/lis/v2/institution/person#Faculty
	InstitutionRoleFaculty = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Faculty"
	// InstitutionRoleGuest http://purl.imsglobal.org/vocab/lis/v2/institution/person#Guest
	InstitutionRoleGuest = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Guest"
	// InstitutionRoleNone http://purl.imsglobal.org/vocab/lis/v2/institution/person#None
	InstitutionRoleNone = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#None"
	// InstitutionRoleOther http://purl.imsglobal.org/vocab/lis/v2/institution/person#Other
	InstitutionRoleOther = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Other"
	// InstitutionRoleStaff http://purl.imsglobal.org/vocab/lis/v2/institution/person#Staff
	InstitutionRoleStaff = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Staff"
	// InstitutionRoleStudent http://purl.imsglobal.org/vocab/lis/v2/institution/person#Student
	InstitutionRoleStudent = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Student"

	//Non‑core institution roles

	// InstitutionRoleAlumni http://purl.imsglobal.org/vocab/lis/v2/institution/person#Alumni
	InstitutionRoleAlumni = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Alumni"
	// InstitutionRoleInstructor http://purl.imsglobal.org/vocab/lis/v2/institution/person#Instructor
	InstitutionRoleInstructor = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Instructor"
	// InstitutionRoleLearner http://purl.imsglobal.org/vocab/lis/v2/institution/person#Learner
	InstitutionRoleLearner = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Learner"
	// InstitutionRoleMember http://purl.imsglobal.org/vocab/lis/v2/institution/person#Member
	InstitutionRoleMember = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Member"
	// InstitutionRoleMentor http://purl.imsglobal.org/vocab/lis/v2/institution/person#Mentor
	InstitutionRoleMentor = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Mentor"
	// InstitutionRoleObserver http://purl.imsglobal.org/vocab/lis/v2/institution/person#Observer
	InstitutionRoleObserver = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Observer"
	// InstitutionRoleProspectiveStudent http://purl.imsglobal.org/vocab/lis/v2/institution/person#ProspectiveStudent
	InstitutionRoleProspectiveStudent = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#ProspectiveStudent"

	//Core context roles

	// ContextRoleAdministrator http://purl.imsglobal.org/vocab/lis/v2/membership#Administrator
	ContextRoleAdministrator = "http://purl.imsglobal.org/vocab/lis/v2/membership#Administrator"
	// ContextRoleContentDeveloper http://purl.imsglobal.org/vocab/lis/v2/membership#ContentDeveloper
	ContextRoleContentDeveloper = "http://purl.imsglobal.org/vocab/lis/v2/membership#ContentDeveloper"
	// ContextRoleInstructor http://purl.imsglobal.org/vocab/lis/v2/membership#Instructor
	ContextRoleInstructor = "http://purl.imsglobal.org/vocab/lis/v2/membership#Instructor"
	// ContextRoleLearner http://purl.imsglobal.org/vocab/lis/v2/membership#Learner
	ContextRoleLearner = "http://purl.imsglobal.org/vocab/lis/v2/membership#Learner"
	// ContextRoleMentor http://purl.imsglobal.org/vocab/lis/v2/membership#Mentor
	ContextRoleMentor = "http://purl.imsglobal.org/vocab/lis/v2/membership#Mentor"

	//Non‑core context roles

	// ContextRoleManager http://purl.imsglobal.org/vocab/lis/v2/membership#Manager
	ContextRoleManager = "http://purl.imsglobal.org/vocab/lis/v2/membership#Manager"
	// ContextRoleMember http://purl.imsglobal.org/vocab/lis/v2/membership#Member
	ContextRoleMember = "http://purl.imsglobal.org/vocab/lis/v2/membership#Member"
	// ContextRoleOfficer http://purl.imsglobal.org/vocab/lis/v2/membership#Officer
	ContextRoleOfficer = "http://purl.imsglobal.org/vocab/lis/v2/membership#Officer"
)

//HasRole check if the message includes the given role
func (rlm LaunchMessage) HasRole(role string) bool {
	for i := range rlm.Roles {
		if rlm.Roles[i] == role {
			return true
		}
	}

	return false
}

//HasAnyRole check if the message includes any of a given list of roles
func (rlm LaunchMessage) HasAnyRole(roles []string) bool {
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
func (rlm LaunchMessage) IsAnonymous() bool {
	return rlm.Sub == nil
}

//DeepLinkingResponse encompasses the entire response from the tool to the platform after a deep linking session
type DeepLinkingResponse struct {
	Iss   string `json:"iss" required:"true"` // Issuer
	Aud   string `json:"aud" required:"true"` // Client ID
	Iat   int    `json:"iat" required:"true"` // Token created timestamp
	Exp   int    `json:"exp" required:"true"` // Token will expire timestamp
	Nonce string `json:"nonce" required:"true"`

	// Message type should always be "LtiDeepLinkingResponse"
	MessageType  string `json:"https://purl.imsglobal.org/spec/lti/claim/message_type" required:"true"`
	Version      string `json:"https://purl.imsglobal.org/spec/lti/claim/version" required:"true"`
	DeploymentID string `json:"https://purl.imsglobal.org/spec/lti/claim/deployment_id" required:"true"`
	Data         string `json:"https://purl.imsglobal.org/spec/lti-dl/claim/data" required:"true"` // Must match "Data" field of request

	// ContentItems contains the selected links from the deep linking flow
	// Nil if no links were selected
	ContentItems []ContentItem `json:"https://purl.imsglobal.org/spec/lti-dl/claim/content_items"`

	Msg      string `json:"https://purl.imsglobal.org/spec/lti-dl/claim/msg,omitempty"`
	Log      string `json:"https://purl.imsglobal.org/spec/lti-dl/claim/log,omitempty"`
	ErrorMsg string `json:"https://purl.imsglobal.org/spec/lti-dl/claim/errormsg,omitempty"`
	ErrorLog string `json:"https://purl.imsglobal.org/spec/lti-dl/claim/errorlog,omitempty"`
}

// AddContentItem add a content item to the deep linking response
// Prevents any items with a duplicate Type+UniqueContent
func (dlr *DeepLinkingResponse) AddContentItem(add ContentItem) {
	if dlr.ContentItems == nil {
		dlr.ContentItems = []ContentItem{}
	}

	contentItems := dlr.ContentItems

	// Check if item already exists, replace if so (just in case, to update ancillary fields)
	for i, ci := range contentItems {
		if ci.GetType() == add.GetType() && ci.GetUniqueContent() == add.GetUniqueContent() {
			contentItems[i] = add
			return
		}
	}

	// If item didn't already exist, add it
	contentItems = append(contentItems, add)
	dlr.ContentItems = contentItems
}

// RemoveContentItem remove a content item from the deep linking response
func (dlr *DeepLinkingResponse) RemoveContentItem(remove ContentItem) {
	dlr.RemoveContentItemByIdentifiers(remove.GetType(), remove.GetUniqueContent())
}

// RemoveContentItemByIdentifiers remove a content item from the deep linking response, given Type and UniqueContent
func (dlr *DeepLinkingResponse) RemoveContentItemByIdentifiers(contentType, uniqueContent string) {
	if dlr.ContentItems == nil {
		// There is no slice to begin with, so no items, so we don't need to do anything
		return
	}

	newSlice := []ContentItem{}
	for _, ci := range dlr.ContentItems {
		if !(ci.GetType() == contentType && ci.GetUniqueContent() == uniqueContent) {
			newSlice = append(newSlice, ci)
		}
	}

	if len(newSlice) == 0 {
		dlr.ContentItems = nil
	} else {
		dlr.ContentItems = newSlice
	}
}

const (
	// ContentItemTypeLink "link"
	ContentItemTypeLink = "link"
	// ContentItemTypeResourceLink "ltiResourceLink"
	ContentItemTypeResourceLink = "ltiResourceLink"
	// ContentItemTypeFile "file"
	ContentItemTypeFile = "file"
	// ContentItemTypeHTMLFragment "html"
	ContentItemTypeHTMLFragment = "html"
	// ContentItemTypeImage "image"
	ContentItemTypeImage = "image"
)

// ContentItem defines common methods for content items
type ContentItem interface {
	GetType() string
	// GetUniqueContent will return the URL, src, or string representation of the content item, such that two content
	// items of the same type with the same return value should be considered functionally equivalent
	GetUniqueContent() string
}

// Link is a fully qualified URL to a resource hosted on the internet
// The item may include different rendering options (window, iframe, embed). As a best practice, the tool SHOULD return
// all the ones that apply, allowing the platform to use the best option based on the actual rendering context when the
// item is displayed.
type Link struct {
	// Value must be "link"
	Type string `json:"type"`
	// Fully qualified URL of the resource. This link must be navigable to.
	URL string `json:"url"`
	// Plain text to use as the title or heading for content. (optional)
	Title string `json:"title,omitempty"`
	// Plain text description of the content item intended to be displayed to all users who can access the item. (optional)
	Text string `json:"text,omitempty"`
	// Fully qualified URL, height, and width of an icon image to be placed with the file. A platform may not support
	// the display of icons, but where it does, it may choose to use a local copy of the icon rather than linking to the
	// URL provided (which would also allow it to resize the image to suit its needs). (optional)
	Icon *Icon `json:"icon,omitempty"`
	// Fully qualified URL, height, and width of a thumbnail image to be made a hyperlink. This allows the hyperlink to
	// be opened within the platform from text or an image, or from both. (optional)
	Thumbnail *Thumbnail `json:"thumbnail,omitempty"`
	// The embed property has a single required property html that contains the HTML fragment to embed the resource
	// directly inside HTML. It is commonly used as a way to embed a resource in an HTML editor. Platform must make
	// sure to properly sanitize the HTML prior to inclusion. (optional)
	Embed string `json:"embed,omitempty"`
	// The window property indicates how to open the resource in a new window/tab. (optional)
	Window *Window `json:"window,omitempty"`
	// The iframe property indicates the resource can be embedded using an IFrame. (optional)
	IFrame *IFrame `json:"iframe,omitempty"`
}

// GetType get the content item type
func (link Link) GetType() string {
	return link.Type
}

// GetUniqueContent get the URL
func (link Link) GetUniqueContent() string {
	return link.URL
}

// ResourceLinkItem link to an LTI resource, usually delivered by the same tool to which the deep linking request was
// made to. A platform may support links associated to other tools. How this association may happen is not specified.
type ResourceLinkItem struct {
	// Value must be "ltiResourceLink"
	Type string `json:"type"`
	// Fully qualified URL of the resource. If absent, the base LTI URL of the tool must be used for launch. (optional)
	URL string `json:"url,omitempty"`
	// Plain text to use as the title or heading for content. (optional)
	Title string `json:"title,omitempty"`
	// Plain text description of the content item intended to be displayed to all users who can access the item. (optional)
	Text string `json:"text,omitempty"`
	// Fully qualified URL, height, and width of an icon image to be placed with the file. A platform may not support
	// the display of icons, but where it does, it may choose to use a local copy of the icon rather than linking to the
	// URL provided (which would also allow it to resize the image to suit its needs). (optional)
	Icon *Icon `json:"icon,omitempty"`
	// Fully qualified URL, height, and width of a thumbnail image to be made a hyperlink. This allows the hyperlink to
	// be opened within the platform from text or an image, or from both. (optional)
	Thumbnail *Thumbnail `json:"thumbnail,omitempty"`
	// The window property indicates how to open the resource in a new window/tab. (optional)
	Window *Window `json:"window,omitempty"`
	// The iframe property indicates the resource can be embedded using an IFrame. (optional)
	IFrame *IFrame `json:"iframe,omitempty"`
	// A map of key/value custom parameters. Those parameters must be included in the LtiResourceLinkRequest payload.
	// Value may include substitution parameters as defined in the LTI Core Specification (optional)
	Custom *map[string]interface{} `json:"custom,omitempty"`
	// A lineItem object that indicates this activity is expected to receive scores; the platform may automatically
	// create a corresponding line item when the resource link is created, using the maximum score as the default
	// maximum points. A line item created as a result of a Deep Linking interaction must be exposed in a subsequent
	// line item service call, with the resourceLinkId of the associated resource link, as well as the resourceId and
	// tag if present in the line item definition (optional)
	LineItem *LineItem `json:"lineItem,omitempty"`
	// Indicates the initial start and end time this activity should be made available to learners (optional)
	Available *TimeRange `json:"available,omitempty"`
	// Indicates the initial start and end time submissions for this activity can be made by learners (optional)
	Submission *TimeRange `json:"submission,omitempty"`
}

// GetType get the content item type
func (rl ResourceLinkItem) GetType() string {
	return rl.Type
}

// GetUniqueContent get the URL
func (rl ResourceLinkItem) GetUniqueContent() string {
	return rl.URL
}

// File a resource transferred from the tool to stored and/or processed by the platform. The URL to the resource should
// be considered short lived and the platform must process the file within a short time frame (within minutes).
type File struct {
	// Value must be "file"
	Type string `json:"type"`
	// Fully qualified URL of the resource. This link may be short-lived or expire after 1st use.
	URL string `json:"url"`
	// Plain text to use as the title or heading for content. (optional)
	Title string `json:"title,omitempty"`
	// Plain text description of the content item intended to be displayed to all users who can access the item. (optional)
	Text string `json:"text,omitempty"`
	// Fully qualified URL, height, and width of an icon image to be placed with the file. A platform may not support
	// the display of icons, but where it does, it may choose to use a local copy of the icon rather than linking to the
	// URL provided (which would also allow it to resize the image to suit its needs). (optional)
	Icon *Icon `json:"icon,omitempty"`
	// Fully qualified URL, height, and width of a thumbnail image to be made a hyperlink. This allows the hyperlink to
	// be opened within the platform from text or an image, or from both. (optional)
	Thumbnail *Thumbnail `json:"thumbnail,omitempty"`
	// ISO 8601 Date and time. The URL will be available until this time. No guarantees after that. (optional)
	ExpiresAt string `json:"expiresAt,omitempty"`
}

// GetType get the content item type
func (file File) GetType() string {
	return file.Type
}

// GetUniqueContent get the URL
func (file File) GetUniqueContent() string {
	return file.URL
}

// HTMLFragment a fragment to be embedded in an HTML document on the platform. If the HTML fragment renders a a single
// resource which is also addressable directly, the tool SHOULD use the link type with an embed code.
type HTMLFragment struct {
	// Value must be "html"
	Type string `json:"type"`
	// HTML fragment to be embedded. The platform is expected to sanitize it against cross-site scripting attacks.
	HTML string `json:"html"`
	// Plain text to use as the title or heading for content. (optional)
	Title string `json:"title,omitempty"`
	// Plain text description of the content item intended to be displayed to all users who can access the item. (optional)
	Text string `json:"text,omitempty"`
}

// GetType get the content item type
func (html HTMLFragment) GetType() string {
	return html.Type
}

// GetUniqueContent get the HTML string
func (html HTMLFragment) GetUniqueContent() string {
	return html.HTML
}

// Image a URL pointing to an image resource that SHOULD be rendered directly in the browser agent using the HTML img tag.
type Image struct {
	// Value must be "image"
	Type string `json:"type"`
	// Fully qualified URL of the image
	URL string `json:"url"`
	// Plain text to use as the title or heading for content. (optional)
	Title string `json:"title,omitempty"`
	// Plain text description of the content item intended to be displayed to all users who can access the item. (optional)
	Text string `json:"text,omitempty"`
	// Fully qualified URL, height, and width of an icon image to be placed with the file. A platform may not support
	// the display of icons, but where it does, it may choose to use a local copy of the icon rather than linking to the
	// URL provided (which would also allow it to resize the image to suit its needs). (optional)
	Icon *Icon `json:"icon,omitempty"`
	// Fully qualified URL, height, and width of a thumbnail image to be made a hyperlink. This allows the hyperlink to
	// be opened within the platform from text or an image, or from both. (optional)
	Thumbnail *Thumbnail `json:"thumbnail,omitempty"`
	// Integer representing the width in pixels of the image
	Width int `json:"width,omitempty"`
	// Integer representing the height in pixels of the image
	Height int `json:"height,omitempty"`
}

// GetType get the content item type
func (img Image) GetType() string {
	return img.Type
}

// GetUniqueContent get the URL
func (img Image) GetUniqueContent() string {
	return img.URL
}

// Icon represents an icon image
type Icon struct {
	// Fully qualified URL to the image file
	URL string `json:"url"`
	// Integer representing the width in pixels of the image
	Width int `json:"width,omitempty"`
	// Integer representing the height in pixels of the image
	Height int `json:"height,omitempty"`
}

// Thumbnail represents a thumbnail image
type Thumbnail struct {
	// Fully qualified URL to the image file
	URL string `json:"url"`
	// Integer representing the width in pixels of the image
	Width int `json:"width,omitempty"`
	// Integer representing the height in pixels of the image
	Height int `json:"height,omitempty"`
}

// Window indicates how to open the resource in a new window/tab
type Window struct {
	// String identifying the name of the window to open; this allows for a single window to be shared as the target of
	// multiple links, preventing a proliferation of new window/tabs. (optional)
	TargetName string `json:"targetName,omitempty"`
	// Integer representing the width in pixels of the new window (optional)
	Width int `json:"width,omitempty"`
	// Integer representing the height in pixels of the new window (optional)
	Height int `json:"height,omitempty"`
	// Comma-separated list of window features as per the window.open() definition
	// (https://developer.mozilla.org/en-US/docs/Web/API/Window/open) (optional)
	WindowFeatures *string `json:"windowFeatures,omitempty"`
}

// IFrame indicates how to open the resource in an IFrame
type IFrame struct {
	// The URL to use as the src of the IFrame. The src value may differ from the link url
	Src string `json:"src"`
	// Integer representing the width in pixels of the new window (optional)
	Width int `json:"width,omitempty"`
	// Integer representing the height in pixels of the new window (optional)
	Height int `json:"height,omitempty"`
}

// TimeRange a start and end time together
type TimeRange struct {
	// ISO8601 start time (optional)
	StartDateTime time.Time `json:"startDateTime,omitempty"`
	// ISO8601 end time (optional)
	EndDateTime time.Time `json:"endDateTime,omitempty"`
}

/* ----------------------------------------------------------------------------
 * Assignments and Grades
 * ------------------------------------------------------------------------- */

// Scopes for AGS calls
const (
	ScopeLineItem       = "https://purl.imsglobal.org/spec/lti-ags/scope/lineitem"
	ScopeResultReadonly = "https://purl.imsglobal.org/spec/lti-ags/scope/result.readonly"
	ScopeScore          = "https://purl.imsglobal.org/spec/lti-ags/scope/score"
)

// LineItem an object that indicates that an activity is expected to receive scores.
type LineItem struct {
	// ID to uniquely identify the line item within a context
	// Optional when sending/creating a line item, because the ID is created by the platform and sent back
	ID string `json:"id,omitempty"`
	// Label for the line item. If not present, the title of the content item must be used (optional)
	Label string `json:"label,omitempty"`
	// Positive decimal value indicating the maximum score possible for this activity
	ScoreMaximum float32 `json:"scoreMaximum"`
	// Tool provided ID for the resource (optional)
	ResourceID string `json:"resourceId,omitempty"`
	// ID of a resource link in the same context as the line item, to connect to it (optional)
	ResourceLinkID string `json:"resourceLinkId,omitempty"`
	// Additional information about the line item; may be used by the tool to identify line items attached to the same
	// resource or resource link (example: grade, originality, participation) (optional)
	Tag string `json:"tag,omitempty"`
	// ISO8601 start time (optional)
	StartDateTime time.Time `json:"startDateTime,omitempty"`
	// ISO8601 end time (optional)
	EndDateTime time.Time `json:"endDateTime,omitempty"`
}

// Activity progress values for a submitted Grade
const (
	ActivityProgressInitialized = "Initialized"
	ActivityProgressStarted     = "Started"
	ActivityProgressInProgress  = "InProgress"
	ActivityProgressSubmitted   = "Submitted"
	ActivityProgressCompleted   = "Completed"
)

// Grading progress values for a submitted Grade
const (
	GradingProgressFullyGraded   = "FullyGraded"
	GradingProgressPending       = "Pending"
	GradingProgressPendingManual = "PendingManual"
	GradingProgressFailed        = "Failed"
	GradingProgressNotReady      = "NotReady"
)

// Grade represents a score to be sent to the platform
// All fields are required
type Grade struct {
	ScoreGiven       float32   `json:"scoreGiven"`
	ScoreMax         float32   `json:"scoreMaximum"`
	ActivityProgress string    `json:"activityProgress"`
	GradingProgress  string    `json:"gradingProgress"`
	Timestamp        time.Time `json:"timestamp"`
	UserID           string    `json:"userId"`
}

// Result represents a score that is received from the platform
type Result struct {
	UserID        string  `json:"userId"`
	ResultScore   float32 `json:"resultScore"`
	ResultMaximum float32 `json:"resultMaximum"`
	Comment       string  `json:"comment"`
	ID            string  `json:"id"`
	ScoreOf       string  `json:"scoreOf"`
}

/* ----------------------------------------------------------------------------
 * Names and Roles Provisioning
 * ------------------------------------------------------------------------- */

// Scope for NRPS call
const (
	ScopeContextMembershipReadonly = "https://purl.imsglobal.org/spec/lti-nrps/scope/contextmembership.readonly"
)

// MemberResponse holds a list of members along with a context
type MemberResponse struct {
	ID      string   `json:"id"`
	Context Context  `json:"context"`
	Members []Member `json:"members"`
}

// Member contains attributes for a single member
type Member struct {
	Name               string   `json:"name"`
	Picture            string   `json:"picture"`
	GivenName          string   `json:"given_name"`
	FamilyName         string   `json:"family_name"`
	MiddleName         string   `json:"middle_name"`
	Email              string   `json:"email"`
	UserID             string   `json:"user_id"`
	LisPersonSourcedid string   `json:"lis_person_sourcedid"`
	Roles              []string `json:"roles"`
}
