package notification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hasura/go-graphql-client"
)

type Client struct {
	client *graphql.Client
}

func New(client *graphql.Client) *Client {
	return &Client{
		client: client,
	}
}

func (c *Client) Send(inputs []*NotificationInput) ([]string, error) {
	if len(inputs) == 0 {
		return []string{}, nil
	}

	var mutation struct {
		CreateNotifications struct {
			Returning []struct {
				ID string `graphql:"id"`
			} `graphql:"returning"`
		} `graphql:"insert_notification(objects: $objects)"`
	}

	objects := make([]notification_insert_input, 0)
	for _, n := range inputs {
		if len(n.Headings) == 0 && len(n.Contents) == 0 {
			continue
		}
		item := notification_insert_input{
			NotificationInput: n,
		}
		if len(n.UserIDs) > 0 {
			item.Users = &notificationUsersInput{}
			for _, uid := range n.UserIDs {
				item.Users.Data = append(item.Users.Data, notificationUserInput{
					UserID: uid,
				})
			}
		}
		if item.SendAfter.IsZero() {
			item.SendAfter = time.Now()
		}
		objects = append(objects, item)
	}
	variables := map[string]interface{}{
		"objects": objects,
	}

	err := c.client.Mutate(context.Background(), &mutation, variables, graphql.OperationName("CreateNotifications"))
	if err != nil {
		return nil, err
	}

	if len(mutation.CreateNotifications.Returning) == 0 {
		return nil, errors.New("insert zero notification")
	}

	var results []string
	for _, r := range mutation.CreateNotifications.Returning {
		results = append(results, r.ID)
	}
	return results, nil
}

func (c *Client) SendWithTemplate(inputs []*NotificationInput, variables interface{}) ([]string, error) {
	if len(inputs) == 0 {
		return []string{}, nil
	}

	templateIDs := uniqueStrings{}
	newInputs := make([]*NotificationInput, 0, len(inputs))

	for _, item := range inputs {
		if item.TemplateID != "" {
			templateIDs.Add(item.TemplateID)
		}
	}

	if !templateIDs.IsEmpty() {
		templates, err := c.GetTemplateByIDs(templateIDs.Value()...)
		if err != nil {
			return nil, err
		}

		for _, item := range inputs {
			if item.TemplateID != "" {
				template, ok := templates[item.TemplateID]
				if !ok {
					return nil, fmt.Errorf("notification template not found: %s", item.TemplateID)
				}
				newItem, err := ParseTemplate(template, variables)
				if err != nil {
					return nil, err
				}
				item.Headings = newItem.Headings
				item.Contents = newItem.Contents
				newInputs = append(newInputs, item)
			} else {
				newInputs = append(newInputs, item)
			}
		}
	}
	return c.Send(newInputs)
}

func (c *Client) GetTemplateByIDs(ids ...string) (map[string]*NotificationTemplate, error) {
	results := make(map[string]*NotificationTemplate)
	if len(ids) == 0 {
		return results, nil
	}

	var query struct {
		NotificationTemplates []NotificationTemplateRaw `graphql:"notification_template(where: $where)" json:"notification_template"`
	}

	variables := map[string]interface{}{
		"where": notification_template_bool_exp{
			"id": map[string]interface{}{
				"_in": ids,
			},
		},
	}

	err := c.client.Query(context.Background(), &query, variables, graphql.OperationName("GetNotificationTemplatesByIds"))
	if err != nil {
		return nil, err
	}

	for _, ntr := range query.NotificationTemplates {
		nt, err := ntr.Parse()
		if err != nil {
			return nil, err
		}
		results[nt.ID] = nt
	}
	return results, nil
}

func (c *Client) UpsertTemplates(inputs []*NotificationTemplate) ([]*NotificationTemplate, error) {
	if len(inputs) == 0 {
		return []*NotificationTemplate{}, nil
	}

	var mutation struct {
		UpsertNotificationTemplates struct {
			Returning []NotificationTemplateRaw `graphql:"returning" json:"returning"`
		} `graphql:"insert_notification_template(objects: $objects, on_conflict: { constraint: notification_template_pkey, update_columns: [contents, headings] })" json:"insert_notification_template"`
	}

	objects := make([]notification_template_insert_input, len(inputs))
	for i, nt := range inputs {
		objects[i] = notification_template_insert_input(*nt)
	}

	variables := map[string]interface{}{
		"objects": objects,
	}

	bytes, err := c.client.NamedMutateRaw(context.Background(), "UpsertNotificationTemplates", &mutation, variables)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(*bytes, &mutation)
	if err != nil {
		return nil, err
	}

	results := make([]*NotificationTemplate, 0, len(mutation.UpsertNotificationTemplates.Returning))

	for _, ntr := range mutation.UpsertNotificationTemplates.Returning {
		nt, err := ntr.Parse()
		if err != nil {
			return nil, err
		}
		results = append(results, nt)
	}
	return results, nil
}

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

	err := c.client.Mutate(context.Background(), &mutation, variables, graphql.OperationName("UpsertNotificationTemplates"))
	if err != nil {
		return 0, err
	}

	return mutation.UpdateNotifications.AffectedRows, nil
}
