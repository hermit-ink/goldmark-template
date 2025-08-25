package goldmarktemplate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

func TestEmphasis(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "template with asterisks in markdown",
			input:    `*{{ .BoldText }}* and **{{ .VeryBold }}**`,
			expected: `<p><em>{{ .BoldText }}</em> and <strong>{{ .VeryBold }}</strong></p>`,
		},
		{
			name:     "template with underscores in markdown",
			input:    `_{{ .ItalicText }}_ and __{{ .BoldText }}__`,
			expected: `<p><em>{{ .ItalicText }}</em> and <strong>{{ .BoldText }}</strong></p>`,
		},
		{
			name:     "template containing asterisks",
			input:    `{{ printf "*%s*" .Text }}`,
			expected: `<p>{{ printf "*%s*" .Text }}</p>`,
		},
		{
			name:     "mixed emphasis and templates",
			input:    `*Bold* and {{ .Template }} and **{{ .BoldTemplate }}**`,
			expected: `<p><em>Bold</em> and {{ .Template }} and <strong>{{ .BoldTemplate }}</strong></p>`,
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

func TestComplexTemplateExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "template with markdown-like characters",
			input:    `{{ printf "[%s](%s)" .Title .URL }}`,
			expected: `<p>{{ printf "[%s](%s)" .Title .URL }}</p>`,
		},
		{
			name:     "template with HTML-like content",
			input:    `{{ printf "<div>%s</div>" .Content }}`,
			expected: `<p>{{ printf "<div>%s</div>" .Content }}</p>`,
		},
		{
			name:     "nested templates with special chars",
			input:    `{{ if .Show }}{{ printf "*%s*" .Text }}{{ end }}`,
			expected: `<p>{{ if .Show }}{{ printf "*%s*" .Text }}{{ end }}</p>`,
		},
		{
			name:     "template in code with autolink-like content",
			input:    "`{{ printf \"<https://example.com>\" }}`",
			expected: `<p><code>{{ printf "<https://example.com>" }}</code></p>`,
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
