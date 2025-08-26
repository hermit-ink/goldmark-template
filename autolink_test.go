package goldmarktemplate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

func TestAutolinkTemplates(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "action URL in angle brackets should create autolink",
			input:    `<{{ .URL }}>`,
			expected: `<p><a href="{{ .URL }}">{{ .URL }}</a></p>`,
		},
		{
			name:     "action with base URL should create autolink",
			input:    `<{{ .BaseURL }}/page>`,
			expected: `<p><a href="{{ .BaseURL }}/page">{{ .BaseURL }}/page</a></p>`,
		},
		{
			name:     "real autolink should still work",
			input:    `<https://example.com>`,
			expected: `<p><a href="https://example.com">https://example.com</a></p>`,
		},
		{
			name:     "real email autolink should still work",
			input:    `<user@example.com>`,
			expected: `<p><a href="mailto:user@example.com">user@example.com</a></p>`,
		},
		{
			name:     "template and real autolink together",
			input:    `Visit <{{ .URL }}> or <https://example.com>`,
			expected: `<p>Visit <a href="{{ .URL }}">{{ .URL }}</a> or <a href="https://example.com">https://example.com</a></p>`,
		},
		{
			name:     "complex action in angle brackets should create autolink",
			input:    `<{{ printf "%s/%s" .BaseURL .Path }}>`,
			expected: `<p><a href="{{ printf "%s/%s" .BaseURL .Path }}">{{ printf "%s/%s" .BaseURL .Path }}</a></p>`,
		},
		{
			name:     "URL with embedded template action",
			input:    `<https://hermit.ink/{{.Foo}}>`,
			expected: `<p><a href="https://hermit.ink/{{.Foo}}">https://hermit.ink/{{.Foo}}</a></p>`,
		},
		{
			name:     "URL with multiple template actions",
			input:    `<https://{{.Host}}/{{.Path}}>`,
			expected: `<p><a href="https://{{.Host}}/{{.Path}}">https://{{.Host}}/{{.Path}}</a></p>`,
		},
		{
			name:     "URL with template action in query parameter",
			input:    `<https://example.com/page?id={{.ID}}>`,
			expected: `<p><a href="https://example.com/page?id={{.ID}}">https://example.com/page?id={{.ID}}</a></p>`,
		},
		{
			name:     "URL with complex template action embedded",
			input:    `<https://api.{{.Domain}}/v1/{{printf "users/%d" .UserID}}>`,
			expected: `<p><a href="https://api.{{.Domain}}/v1/{{printf "users/%d" .UserID}}">https://api.{{.Domain}}/v1/{{printf "users/%d" .UserID}}</a></p>`,
		},
	}

	md := goldmark.New(
		goldmark.WithExtensions(New()),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
			html.WithXHTML(),
		),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := md.Convert([]byte(tt.input), &buf)
			if err != nil {
				t.Fatalf("Failed to convert markdown: %v", err)
			}

			got := strings.TrimSpace(buf.String())
			expected := strings.TrimSpace(tt.expected)

			if got != expected {
				t.Errorf("Output mismatch\nInput:    %q\nExpected: %q\nGot:      %q", tt.input, expected, got)
			}
		})
	}
}
