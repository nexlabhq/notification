package notification

import (
	"encoding/json"
	"sort"
	"strings"
	"time"
)

type notification_bool_exp map[string]interface{}
type notification_set_input map[string]interface{}
type notification_template_bool_exp map[string]interface{}

type NotificationMetadata struct {
	CaseId    string                 `json:"case_id"`
	SessionId string                 `json:"session_id"`
	Color     string                 `json:"color"`
	URL       string                 `json:"url"`
	ImageURL  string                 `json:"image_url"`
	Subtitles map[string]interface{} `json:"subtitles"`
}

type notificationUserInput struct {
	UserID string `json:"user_id"`
}
type notificationUsersInput struct {
	Data []notificationUserInput `json:"data"`
}

type NotificationInput struct {
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
}

type notification_insert_input struct {
	*NotificationInput
	Users *notificationUsersInput `json:"users,omitempty"`
}

type NotificationTemplate struct {
	ID       string               `graphql:"id" json:"id"`
	Headings map[string]string    `graphql:"headings" json:"headings"`
	Contents map[string]string    `graphql:"contents"  json:"contents"`
	Metadata NotificationMetadata `graphql:"metadata"  json:"metadata"`
}

type notification_template_insert_input NotificationTemplate

type NotificationTemplateRaw struct {
	ID       string          `graphql:"id" json:"id"`
	Headings json.RawMessage `graphql:"headings" json:"headings"`
	Contents json.RawMessage `graphql:"contents"  json:"contents"`
	Metadata json.RawMessage `graphql:"metadata"  json:"metadata"`
}

func (ntr *NotificationTemplateRaw) Parse() (*NotificationTemplate, error) {

	var headings map[string]string
	var contents map[string]string
	var metadata NotificationMetadata
	if len(ntr.Contents) > 0 {
		err := json.Unmarshal(ntr.Contents, &contents)
		if err != nil {
			return nil, err
		}
	}
	if len(ntr.Headings) > 0 {
		err := json.Unmarshal(ntr.Headings, &headings)
		if err != nil {
			return nil, err
		}
	}

	if len(ntr.Metadata) > 0 {
		err := json.Unmarshal(ntr.Metadata, &metadata)
		if err != nil {
			return nil, err
		}
	}
	return &NotificationTemplate{
		ID:       ntr.ID,
		Headings: headings,
		Contents: contents,
		Metadata: metadata,
	}, nil
}

// uniqueStrings is the special array string that only store unique values
type uniqueStrings map[string]bool

// Add append new value or skip if it's existing
func (us uniqueStrings) Add(values ...string) {
	for _, s := range values {
		if _, ok := us[s]; !ok {
			us[s] = true
		}
	}
}

// IsEmpty check if the array is empty
func (us uniqueStrings) IsEmpty() bool {
	return len(us) == 0
}

// Value return
func (us uniqueStrings) Value() []string {
	results := make([]string, 0, len(us))
	for k := range us {
		results = append(results, k)
	}
	return results
}

// String implement string interface
func (us uniqueStrings) String() string {
	results := us.Value()
	sort.Strings(results)
	return strings.Join(results, ",")
}
