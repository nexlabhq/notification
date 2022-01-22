package notification

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// ToClientName convert multiple client names to string
func ToClientName(name string, names ...string) string {
	if len(names) == 0 {
		return name
	}

	return strings.Join(append([]string{name}, names...), ",")
}

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
