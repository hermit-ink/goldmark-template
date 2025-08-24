package goldmarktemplate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

func TestBasicLinks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "link with template in URL",
			input:    "[Link]({{ .URL }})",
			expected: "<p><a href=\"{{ .URL }}\">Link</a></p>",
		},
		{
			name:     "link with template containing special chars in URL",
			input:    "[Link]({{ .BaseURL }}/path?param={{ .Value }})",
			expected: "<p><a href=\"{{ .BaseURL }}/path?param={{ .Value }}\">Link</a></p>",
		},
		{
			name:     "link with template in title and URL",
			input:    "[{{ .Title }}]({{ .URL }})",
			expected: "<p><a href=\"{{ .URL }}\">{{ .Title }}</a></p>",
		},
		{
			name:     "template with parentheses",
			input:    "[text]({{.Func(arg)}})",
			expected: "<p><a href=\"{{.Func(arg)}}\">text</a></p>",
		},
		{
			name:     "multiple templates in URL",
			input:    "[link](https://{{.Host}}/{{.Path}}?id={{.ID}})",
			expected: "<p><a href=\"https://{{.Host}}/{{.Path}}?id={{.ID}}\">link</a></p>",
		},
	}

	md := goldmark.New(
		goldmark.WithParser(NewTemplatedParser()),
		goldmark.WithExtensions(NewTemplatedHTMLExtension()),
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

func TestLinkTitles(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "link with plain title",
			input:    `[link]({{.URL}} "with a title")`,
			expected: `<p><a href="{{.URL}}" title="with a title">link</a></p>`,
		},
		{
			name:     "link with template in title",
			input:    `[link]({{.URL}} "Title: {{.Title}}")`,
			expected: `<p><a href="{{.URL}}" title="Title: {{.Title}}">link</a></p>`,
		},
		{
			name:     "mailto with template and title",
			input:    `[link](mailto:{{.Email}} "Email {{- .Author }}")`,
			expected: `<p><a href="mailto:{{.Email}}" title="Email {{- .Author }}">link</a></p>`,
		},
		{
			name:     "template with quotes inside",
			input:    `[link]({{.URL}} "{{.Title | quote}}")`,
			expected: `<p><a href="{{.URL}}" title="{{.Title | quote}}">link</a></p>`,
		},
	}

	md := goldmark.New(
		goldmark.WithParser(NewTemplatedParser()),
		goldmark.WithExtensions(NewTemplatedHTMLExtension()),
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

func TestComplexLinkText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "bold template in link text",
			input:    "[**{{.BoldTitle}}** with other text]({{.URL}})",
			expected: `<p><a href="{{.URL}}"><strong>{{.BoldTitle}}</strong> with other text</a></p>`,
		},
		{
			name:     "italic and code with template in link text",
			input:    "[*text* about `{{.CodeReference}}`](https://example.com/{{.Path}})",
			expected: `<p><a href="https://example.com/{{.Path}}"><em>text</em> about <code>{{.CodeReference}}</code></a></p>`,
		},
		{
			name:     "nested inline styles with templates",
			input:    "[**Bold {{.Var1}}** and *italic {{.Var2}}*]({{.URL}})",
			expected: `<p><a href="{{.URL}}"><strong>Bold {{.Var1}}</strong> and <em>italic {{.Var2}}</em></a></p>`,
		},
		{
			name:     "regular link followed by template link",
			input:    "[normal](https://example.com) and [template]({{.URL}})",
			expected: `<p><a href="https://example.com">normal</a> and <a href="{{.URL}}">template</a></p>`,
		},
		{
			name:     "template link followed by regular link",
			input:    "[template]({{.URL}}) and [normal](https://example.com)",
			expected: `<p><a href="{{.URL}}">template</a> and <a href="https://example.com">normal</a></p>`,
		},
	}

	md := goldmark.New(
		goldmark.WithParser(NewTemplatedParser()),
		goldmark.WithExtensions(NewTemplatedHTMLExtension()),
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