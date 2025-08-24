package goldmarktemplate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

func TestHTMLAttributes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "raw HTML with template in attribute",
			input:    "<div id=\"{{ .ID }}\">Content</div>",
			expected: "<div id=\"{{ .ID }}\">Content</div>",
		},
		{
			name:     "raw HTML with template in class attribute",
			input:    "<span class=\"{{ .CSSClass }}\">Text</span>",
			expected: "<p><span class=\"{{ .CSSClass }}\">Text</span></p>",
		},
		{
			name:     "raw HTML with template containing quotes in attribute",
			input:    "<div data-value=\"{{ .Value | quote }}\">Content</div>",
			expected: "<div data-value=\"{{ .Value | quote }}\">Content</div>",
		},
		{
			name:     "raw HTML with multiple template attributes",
			input:    "<a href=\"{{ .URL }}\" title=\"{{ .Title }}\" class=\"{{ .Class }}\">Link</a>",
			expected: "<p><a href=\"{{ .URL }}\" title=\"{{ .Title }}\" class=\"{{ .Class }}\">Link</a></p>",
		},
		{
			name:     "raw HTML attribute with nested templates",
			input:    "<div data-config=\"{{ range .Items }}{{ .Name }},{{ end }}\">Content</div>",
			expected: "<div data-config=\"{{ range .Items }}{{ .Name }},{{ end }}\">Content</div>",
		},
	}

	md := goldmark.New(
		goldmark.WithParser(NewParser()),
		goldmark.WithExtensions(NewExtension()),
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

func TestHTMLContent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "template in raw HTML content",
			input:    "<div>{{ .Content }}</div>",
			expected: "<div>{{ .Content }}</div>",
		},
		{
			name:     "mixed template and HTML with special chars",
			input:    "<p>Value: {{ if .Value > 0 }}{{ .Value }}{{ else }}N/A{{ end }}</p>",
			expected: "<p>Value: {{ if .Value > 0 }}{{ .Value }}{{ else }}N/A{{ end }}</p>",
		},
	}

	md := goldmark.New(
		goldmark.WithParser(NewParser()),
		goldmark.WithExtensions(NewExtension()),
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