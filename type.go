package notification

import (
	"time"
)

const (
	AllClients = "all"
)

type json map[string]interface{}
type notification_bool_exp map[string]interface{}
type notification_set_input map[string]interface{}

type NotificationMetadata struct {
	CaseId    string                 `json:"case_id"`
	SessionId string                 `json:"session_id"`
	Color     string                 `json:"color"`
	URL       string                 `json:"url"`
	ImageURL  string                 `json:"image_url"`
	Subtitles map[string]interface{} `json:"subtitles"`
}

type SendNotificationInput struct {
	AppID       string                `json:"api_id,omitempty"`
	ClientName  string                `json:"client_name,omitempty"`
	TemplateID  string                `json:"template_id,omitempty"`
	Broadcast   bool                  `json:"broadcast"`
	Headings    map[string]string     `json:"headings,omitempty"`
	Contents    map[string]string     `json:"contents,omitempty"`
	SubjectType string                `json:"subject_type,omitempty"`
	SubjectID   string                `json:"subject_id,omitempty"`
	Topics      []string              `json:"topics,omitempty"`
	UserIDs     []string              `json:"user_ids,omitempty"`
	SendAfter   time.Time             `json:"send_after,omitempty"`
	Data        map[string]string     `json:"data,omitempty"`
	Metadata    *NotificationMetadata `json:"metadata,omitempty"`
	Visible     bool                  `json:"visible,omitempty"`
	Save        bool                  `json:"save,omitempty"`
}

type SendResponse struct {
	Success           bool        `json:"success" graphql:"success"`
	RateLimitExceeded bool        `json:"rate_limit_exceeded" graphql:"rate_limit_exceeded"`
	ClientName        string      `json:"client_name,omitempty" graphql:"client_name"`
	RequestID         string      `json:"request_id,omitempty" graphql:"request_id"`
	MessageID         string      `json:"message_id,omitempty" graphql:"message_id"`
	Error             interface{} `json:"error,omitempty" graphql:"error"`
}

type SendNotificationOutput struct {
	Responses    []*SendResponse `json:"responses" graphql:"responses"`
	SuccessCount int             `json:"success_count" graphql:"success_count"`
	FailureCount int             `json:"failure_count" graphql:"failure_count"`
}
