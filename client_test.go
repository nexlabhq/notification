//go:build integration
// +build integration

package notification

import (
	"net/http"
	"testing"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/stretchr/testify/assert"
)

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
	httpClient := &http.Client{
		Transport: &hasuraTransport{
			rt:          http.DefaultTransport,
			adminSecret: "hasura",
		},
		Timeout: 30 * time.Second,
	}
	return graphql.NewClient("http://localhost:8080/v1/graphql", httpClient)
}

func TestGetNotificationTemplates(t *testing.T) {
	client := New(newGqlClient())

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
			Topics: []string{"test"},
			Metadata: &NotificationMetadata{
				ImageURL: "https://en.wikipedia.org/static/images/project-logos/enwiki.png",
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(results))
}

func TestCancelNotifications(t *testing.T) {
	client := New(newGqlClient())
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
