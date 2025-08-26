package goldmarktemplate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

func TestHardLineBreaks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "hard line break with template",
			input:    "{{ .FirstLine }}  \n{{ .SecondLine }}",
			expected: "<p>{{ .FirstLine }}<br />\n{{ .SecondLine }}</p>",
		},
		{
			name:     "template before hard line break",
			input:    "Text {{ .Action }}  \nNext line",
			expected: "<p>Text {{ .Action }}<br />\nNext line</p>",
		},
		{
			name:     "template after hard line break",
			input:    "First line  \n{{ .Action }} text",
			expected: "<p>First line<br />\n{{ .Action }} text</p>",
		},
		{
			name:     "multiple templates with hard line breaks",
			input:    "{{ .First }}  \n{{ .Second }}  \n{{ .Third }}",
			expected: "<p>{{ .First }}<br />\n{{ .Second }}<br />\n{{ .Third }}</p>",
		},
		{
			name:     "hard line break in emphasis with template",
			input:    "*{{ .Important }}  \ncontinues here*",
			expected: "<p><em>{{ .Important }}<br />\ncontinues here</em></p>",
		},
		{
			name:     "hard line break with complex template",
			input:    "{{ if .Show }}{{ .Content }}{{ end }}  \nNext line",
			expected: "<p>{{ if .Show }}{{ .Content }}{{ end }}<br />\nNext line</p>",
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

func TestSoftLineBreaks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "soft line break with template",
			input:    "{{ .FirstLine }}\n{{ .SecondLine }}",
			expected: "<p>{{ .FirstLine }}\n{{ .SecondLine }}</p>",
		},
		{
			name:     "template spanning soft line break",
			input:    "Text {{ .Action }}\ncontinues here",
			expected: "<p>Text {{ .Action }}\ncontinues here</p>",
		},
		{
			name:     "multiple soft line breaks with templates",
			input:    "{{ .First }}\n{{ .Second }}\n{{ .Third }}",
			expected: "<p>{{ .First }}\n{{ .Second }}\n{{ .Third }}</p>",
		},
		{
			name:     "soft line break in emphasis with template",
			input:    "*{{ .Important }}\ncontinues here*",
			expected: "<p><em>{{ .Important }}\ncontinues here</em></p>",
		},
		{
			name:     "soft line break with mixed content",
			input:    "Start {{ .Variable }}\nend of paragraph",
			expected: "<p>Start {{ .Variable }}\nend of paragraph</p>",
		},
		{
			name:     "template with internal newlines preserved",
			input:    "{{ range .Items }}\n{{ .Name }}\n{{ end }}",
			expected: "<p>{{ range .Items }}\n{{ .Name }}\n{{ end }}</p>",
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

func TestLineBreaksWithOtherElements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "hard line break in blockquote",
			input:    "> {{ .Quote }}  \n> continues here",
			expected: "<blockquote>\n<p>{{ .Quote }}<br />\ncontinues here</p>\n</blockquote>",
		},
		{
			name:     "hard line break in list item",
			input:    "- {{ .Item }}  \n  continues here",
			expected: "<ul>\n<li>{{ .Item }}<br />\ncontinues here</li>\n</ul>",
		},
		{
			name:     "line break with code span",
			input:    "`{{ .Code }}`  \nNext line",
			expected: "<p><code>{{ .Code }}</code><br />\nNext line</p>",
		},
		{
			name:     "line break with link",
			input:    "[{{ .LinkText }}]({{ .URL }})  \nNext line",
			expected: "<p><a href=\"{{ .URL }}\">{{ .LinkText }}</a><br />\nNext line</p>",
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