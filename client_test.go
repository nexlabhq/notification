//go:build integration
// +build integration

package notification

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/stretchr/testify/assert"
)

func cleanup(t *testing.T, client *Client) {
	_, err := client.DeleteNotifications(map[string]interface{}{})
	assert.Nil(t, err)

	_, err = client.DeleteNotifications(map[string]interface{}{})
	assert.Nil(t, err)
}

// hasuraTransport transport for Hasura GraphQL Client
type hasuraTransport struct {
	adminSecret string
	headers     map[string]string
	// keep a reference to the client's original transport
	rt http.RoundTripper
}

// RoundTrip set header data before executing http request
func (t *hasuraTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if t.adminSecret != "" {
		r.Header.Set("X-Hasura-Admin-Secret", t.adminSecret)
	}
	for k, v := range t.headers {
		r.Header.Set(k, v)
	}
	return t.rt.RoundTrip(r)
}

func newGqlClient() *graphql.Client {
	adminSecret := os.Getenv("HASURA_GRAPHQL_ADMIN_SECRET")
	httpClient := &http.Client{
		Transport: &hasuraTransport{
			rt:          http.DefaultTransport,
			adminSecret: adminSecret,
		},
		Timeout: 30 * time.Second,
	}
	return graphql.NewClient(os.Getenv("DATA_URL"), httpClient)
}

func TestGetNotificationTemplates(t *testing.T) {

	client := New(newGqlClient())
	defer cleanup(t, client)

	templates, err := client.UpsertTemplates([]*NotificationTemplate{
		{
			ID: "test_template",
			Headings: map[string]string{
				"en": "Test headings en",
				"vi": "Test headings vi",
			},
			Contents: map[string]string{
				"en": "Test contents en",
				"vi": "Test contents vi",
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(templates))

	returnedTemplates, err := client.GetTemplateByIDs("test_template")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(templates))

	assert.Equal(t, templates[0].ID, returnedTemplates["test_template"].ID)
	assert.Equal(t, templates[0].Headings, returnedTemplates["test_template"].Headings)
	assert.Equal(t, templates[0].Contents, returnedTemplates["test_template"].Contents)
}

func TestSendNotifications(t *testing.T) {

	client := New(newGqlClient())
	defer cleanup(t, client)

	headings := "Test headings"
	contents := "Test contents"
	results, err := client.Send([]*NotificationInput{
		{
			ClientName: ToClientName("app1", "app2"),
			Headings: map[string]string{
				"en": headings,
				"vi": headings,
			},
			Contents: map[string]string{
				"en": contents,
				"vi": contents,
			},
			Topics: []string{"test"},
			Metadata: &NotificationMetadata{
				ImageURL: "https://en.wikipedia.org/static/images/project-logos/enwiki.png",
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(results))

	var getQuery struct {
		Notifications []struct {
			ClientName string `graphql:"client_name"`
		} `graphql:"notification(where: $where)"`
	}

	getVariables := map[string]interface{}{
		"where": notification_bool_exp{
			"id": map[string]interface{}{
				"_eq": results[0],
			},
		},
	}
	err = client.client.Query(context.TODO(), &getQuery, getVariables)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(getQuery.Notifications))
	assert.Equal(t, "app1,app2", getQuery.Notifications[0].ClientName)
}

func TestCancelNotifications(t *testing.T) {
	client := New(newGqlClient())
	defer cleanup(t, client)

	headings := "Test headings"
	contents := "Test contents"
	results, err := client.Send([]*NotificationInput{
		{
			Headings: map[string]string{
				"en": headings,
				"vi": headings,
			},
			Contents: map[string]string{
				"en": contents,
				"vi": contents,
			},
			SubjectType: "test",
			SubjectID:   "test_id",
			SendAfter:   time.Now().Add(time.Hour),
			Topics:      []string{"test"},
			Metadata: &NotificationMetadata{
				ImageURL: "https://en.wikipedia.org/static/images/project-logos/enwiki.png",
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(results))

	canceledCount, err := client.CancelNotificationsBySubject("test", "test_id")
	assert.NoError(t, err)
	assert.Equal(t, 1, canceledCount)
}
