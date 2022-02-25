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

func cleanup(t *testing.T, client *graphql.Client) {

	var mutation struct {
		DeleteNotifications struct {
			AffectedRows int `graphql:"affected_rows"`
		} `graphql:"delete_notification(where: $where)"`
	}

	variables := map[string]interface{}{
		"where": notification_bool_exp{},
	}

	err := client.Mutate(context.Background(), &mutation, variables, graphql.OperationName("DeleteNotifications"))
	assert.NoError(t, err)
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

func TestSendNotifications(t *testing.T) {

	client := New(newGqlClient())
	defer cleanup(t, client.client)

	headings := "Test headings"
	contents := "Test contents"
	results, err := client.Send([]*SendNotificationInput{
		{
			ClientName: ToClientName("default", "test2"),
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
			Save: true,
		},
	}, nil)
	assert.NoError(t, err)

	var getQuery struct {
		Notifications []struct {
			ClientName string `graphql:"client_name"`
		} `graphql:"notification(where: $where)"`
	}

	getVariables := map[string]interface{}{
		"where": notification_bool_exp{
			"id": map[string]interface{}{
				"_eq": results.Responses[0].RequestID,
			},
		},
	}
	err = client.client.Query(context.TODO(), &getQuery, getVariables)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(getQuery.Notifications))
	assert.Equal(t, "default,test2", getQuery.Notifications[0].ClientName)
}

func TestCancelNotifications(t *testing.T) {
	client := New(newGqlClient())
	defer cleanup(t, client.client)

	headings := "Test headings"
	contents := "Test contents"
	results, err := client.Send([]*SendNotificationInput{
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
	}, nil)
	assert.NoError(t, err)
	assert.True(t, results.SuccessCount > 0)

	canceledCount, err := client.CancelNotificationsBySubject("test", "test_id")
	assert.NoError(t, err)
	assert.Equal(t, 1, canceledCount)
}
