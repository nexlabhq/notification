package notification

import (
	"bytes"
	"fmt"
	"text/template"
)

func ParseTemplate(nt *NotificationTemplate, variables interface{}) (*NotificationTemplate, error) {
	headings := make(map[string]string)
	contents := make(map[string]string)

	for k, v := range nt.Headings {
		tName := fmt.Sprintf("%s:%s:%s", nt.ID, "heading", k)
		t, err := template.New(tName).Parse(v)
		if err != nil {
			return nil, err
		}
		var b bytes.Buffer
		if err = t.Execute(&b, variables); err != nil {
			return nil, err
		}
		headings[k] = b.String()
	}

	for k, v := range nt.Contents {
		tName := fmt.Sprintf("%s:%s:%s", nt.ID, "content", k)
		t, err := template.New(tName).Parse(v)
		if err != nil {
			return nil, err
		}
		var b bytes.Buffer
		if err = t.Execute(&b, variables); err != nil {
			return nil, err
		}
		contents[k] = b.String()
	}

	return &NotificationTemplate{
		ID:       nt.ID,
		Headings: headings,
		Contents: contents,
	}, nil
}
