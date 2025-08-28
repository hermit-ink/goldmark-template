package goldmarktemplate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func TestBasicAttributeTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "ID attribute with template action",
			input:    `# Heading {id="{{ .HeadingID }}"}`,
			expected: `<h1 id="{{ .HeadingID }}">Heading</h1>`,
		},
		{
			name:     "class attribute with template action",
			input:    `# Heading {class="{{ .CSSClass }}"}`,
			expected: `<h1 class="{{ .CSSClass }}">Heading</h1>`,
		},
		{
			name:     "multiple attributes with templates",
			input:    `# Heading {class="static" data-dynamic="{{ .DynamicClass }}"}`,
			expected: `<h1 class="static" data-dynamic="{{ .DynamicClass }}">Heading</h1>`,
		},
		{
			name:     "key-value attribute with template",
			input:    `# Heading {data-value="{{ .Data }}"}`,
			expected: `<h1 data-value="{{ .Data }}">Heading</h1>`,
		},
		{
			name:     "combined attributes with templates",
			input:    `# Heading {id="{{ .ID }}" class="{{ .Class }}" data-attr="{{ .Data }}"}`,
			expected: `<h1 id="{{ .ID }}" class="{{ .Class }}" data-attr="{{ .Data }}">Heading</h1>`,
		},
	}

	md := goldmark.New(
		goldmark.WithExtensions(WithParserOptions(
			parser.WithAttribute(),
			parser.WithHeadingAttribute(),
		)),
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


func TestComplexAttributeScenarios(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "attribute with complex template expression",
			input:    `# Heading {data-json="{{ .Config | toJson }}"}`,
			expected: `<h1 data-json="{{ .Config | toJson }}">Heading</h1>`,
		},
		{
			name:     "attribute with template containing quotes",
			input:    `# Heading {data-msg="{{ printf \"Hello %s\" .Name }}"}`,
			expected: `<h1 data-msg="{{ printf \"Hello %s\" .Name }}">Heading</h1>`,
		},
		{
			name:     "multiple attributes with different template types",
			input:    `# Heading {id="heading-{{ .ID }}" class="main {{ .Type }}-class" data-count="{{ .Count }}"}`,
			expected: `<h1 id="heading-{{ .ID }}" class="main {{ .Type }}-class" data-count="{{ .Count }}">Heading</h1>`,
		},
		{
			name:     "nested template actions in attribute",
			input:    `# Heading {data-config="{{ if .Debug }}debug{{ else }}prod{{ end }}"}`,
			expected: `<h1 data-config="{{ if .Debug }}debug{{ else }}prod{{ end }}">Heading</h1>`,
		},
	}

	md := goldmark.New(
		goldmark.WithExtensions(WithParserOptions(
			parser.WithAttribute(),
			parser.WithHeadingAttribute(),
		)),
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

