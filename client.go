package notification

import (
	"context"
	"strings"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/hgiasac/graphql-utils/client"
)

// Client represents a generic notification client
type Client struct {
	client client.Client
}

// New creates a Client instance
func New(client client.Client) *Client {
	return &Client{
		client: client,
	}
}

// Send sends create notifications request
func (c *Client) Send(inputs []*SendNotificationInput, variables map[string]string) (*SendNotificationOutput, error) {
	if len(inputs) == 0 {
		return &SendNotificationOutput{}, nil
	}

	for _, input := range inputs {
		if input.SendAfter.IsZero() {
			input.SendAfter = time.Now()
		}
	}
	var mutation struct {
		SendNotifications SendNotificationOutput `graphql:"sendNotifications(data: $data, variables: $variables)"`
	}

	inputVariables := map[string]interface{}{
		"data":      inputs,
		"variables": json(variables),
	}

	err := c.client.Mutate(context.Background(), &mutation, inputVariables, graphql.OperationName("SendNotifications"))
	if err != nil {
		return nil, err
	}

	return &mutation.SendNotifications, nil
}

// CancelNotificationsBySubject cancel and update notifications by subject
func (c *Client) CancelNotificationsBySubject(subjectType string, subjectId string) (int, error) {
	variables := map[string]interface{}{
		"subject_id": map[string]string{
			"_eq": subjectId,
		},
		"send_after": map[string]time.Time{
			"_gt": time.Now(),
		},
	}

	if subjectType != "" {
		variables["subject_type"] = map[string]string{
			"_eq": subjectType,
		}
	}
	return c.CancelNotifications(variables)
}

// CancelNotifications cancel and update notifications
func (c *Client) CancelNotifications(where map[string]interface{}) (int, error) {

	var mutation struct {
		UpdateNotifications struct {
			AffectedRows int `graphql:"affected_rows"`
		} `graphql:"update_notification(where: $where, _set: $setValues)"`
	}

	variables := map[string]interface{}{
		"where": notification_bool_exp(where),
		"setValues": notification_set_input{
			"closed":  true,
			"visible": false,
		},
	}

	err := c.client.Mutate(context.Background(), &mutation, variables, graphql.OperationName("CancelNotifications"))
	if err != nil {
		return 0, err
	}

	return mutation.UpdateNotifications.AffectedRows, nil
}

// ToClientName convert multiple client names to string
func ToClientName(name string, names ...string) string {
	if len(names) == 0 {
		return name
	}

	return strings.Join(append([]string{name}, names...), ",")
}
