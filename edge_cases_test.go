package goldmarktemplate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func TestNestedTemplates(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "action with string containing braces",
			input:    "`{{ printf \"Value: %d\" .Count }}`",
			expected: "<p><code>{{ printf \"Value: %d\" .Count }}</code></p>",
		},
		{
			name:     "action with string containing }}",
			input:    "`{{ \"a string with }}\" }}`",
			expected: "<p><code>{{ \"a string with }}\" }}</code></p>",
		},
		{
			name:     "deeply nested actions",
			input:    "`{{ if .A }}{{ if .B }}{{ .C }}{{ end }}{{ end }}`",
			expected: "<p><code>{{ if .A }}{{ if .B }}{{ .C }}{{ end }}{{ end }}</code></p>",
		},
	}

	md := goldmark.New(
		goldmark.WithExtensions(New()),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
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

func TestSpecialCharacters(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty template",
			input:    "`{{}}`",
			expected: "<p><code>{{}}</code></p>",
		},
		{
			name:     "action with only spaces",
			input:    "`{{   }}`",
			expected: "<p><code>{{   }}</code></p>",
		},
		{
			name:     "unmatched opening braces",
			input:    "`{{ {{ .Name }}`",
			expected: "<p><code>{{ {{ .Name }}</code></p>",
		},
		{
			name:     "action at start and end",
			input:    "`{{ .Start }}content{{ .End }}`",
			expected: "<p><code>{{ .Start }}content{{ .End }}</code></p>",
		},
		{
			name:     "incomplete template action",
			input:    "`{{ .Name` without closing",
			expected: "<p><code>{{ .Name</code> without closing</p>",
		},
		{
			name:     "action with backticks in string",
			input:    "`{{ `raw string` }}`",
			expected: "<p><code>{{ </code>raw string<code> }}</code></p>",
		},
		{
			name:     "action content with special chars",
			input:    "`<div>&amp;</div>`",
			expected: "<p><code>&lt;div&gt;&amp;amp;&lt;/div&gt;</code></p>",
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

func TestMixedContent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "mixed content",
			input:    "Text with `{{ .Code }}` and `regular code`.",
			expected: "<p>Text with <code>{{ .Code }}</code> and <code>regular code</code>.</p>",
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
