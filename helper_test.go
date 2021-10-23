package notification

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotificationTemplateParser(t *testing.T) {
	templates := []struct {
		Input     NotificationTemplate
		Variables interface{}
		Output    NotificationTemplate
	}{
		{
			NotificationTemplate{
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
			nil,
			NotificationTemplate{
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
		},
		{
			NotificationTemplate{
				ID: "test_template_1",
				Headings: map[string]string{
					"en": "Test headings {{.Foo}}",
					"vi": "Test headings {{.Bar}}",
				},
				Contents: map[string]string{
					"en": "Test contents {{.Foo}}",
					"vi": "Test contents {{.Bar}}",
				},
			},
			struct {
				Foo string
				Bar string
			}{
				"foo",
				"bar",
			},
			NotificationTemplate{
				ID: "test_template_1",
				Headings: map[string]string{
					"en": "Test headings foo",
					"vi": "Test headings bar",
				},
				Contents: map[string]string{
					"en": "Test contents foo",
					"vi": "Test contents bar",
				},
			},
		},
		{
			NotificationTemplate{
				ID: "test_template_1",
				Headings: map[string]string{
					"en": "Test headings {{.Foo}}",
					"vi": "Test headings {{.Bar}}",
				},
				Contents: map[string]string{
					"en": "Test contents {{.Foo}}",
					"vi": "Test contents {{.Bar}}",
				},
			},
			struct {
				Foo string
				Bar string
			}{
				"foo",
				"bar",
			},
			NotificationTemplate{
				ID: "test_template_1",
				Headings: map[string]string{
					"en": "Test headings foo",
					"vi": "Test headings bar",
				},
				Contents: map[string]string{
					"en": "Test contents foo",
					"vi": "Test contents bar",
				},
			},
		},
	}

	for _, template := range templates {
		parsedTemplates, err := ParseTemplate(&template.Input, template.Variables)
		assert.NoError(t, err)

		assert.Equal(t, template.Output.ID, parsedTemplates.ID)
		assert.Equal(t, template.Output.Headings, parsedTemplates.Headings)
		assert.Equal(t, template.Output.Contents, parsedTemplates.Contents)
	}

}
