package goldmarktemplate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

func TestCodeSpans(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple variable in code span",
			input:    "`{{ .Name }}`",
			expected: "<p><code>{{ .Name }}</code></p>",
		},
		{
			name:     "action with comparison in code span",
			input:    "`{{ if .Value > 5 }}high{{ end }}`",
			expected: "<p><code>{{ if .Value > 5 }}high{{ end }}</code></p>",
		},
		{
			name:     "action with less than in code span",
			input:    "`{{ if .Value < 10 }}low{{ end }}`",
			expected: "<p><code>{{ if .Value < 10 }}low{{ end }}</code></p>",
		},
		{
			name:     "action with ampersand in code span",
			input:    "`{{ .Title & .Subtitle }}`",
			expected: "<p><code>{{ .Title & .Subtitle }}</code></p>",
		},
		{
			name:     "multiple actions in one code span",
			input:    "`{{ .First }} and {{ .Second }}`",
			expected: "<p><code>{{ .First }} and {{ .Second }}</code></p>",
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

func TestCodeBlocks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple action in code block",
			input:    "```\n{{ .Name }}\n```",
			expected: "<pre><code>{{ .Name }}\n</code></pre>",
		},
		{
			name:     "action with HTML chars in code block",
			input:    "```\n{{ if .Value > 5 }}<div>High</div>{{ end }}\n```",
			expected: "<pre><code>{{ if .Value > 5 }}&lt;div&gt;High&lt;/div&gt;{{ end }}\n</code></pre>",
		},
		{
			name:     "action in go code block",
			input:    "```go\nfmt.Println({{ .Value }})\n```",
			expected: "<pre><code class=\"language-go\">fmt.Println({{ .Value }})\n</code></pre>",
		},
		{
			name:     "action in html code block",
			input:    "```html\n<div>{{ .Content }}</div>\n```",
			expected: "<pre><code class=\"language-html\">&lt;div&gt;{{ .Content }}&lt;/div&gt;\n</code></pre>",
		},
		{
			name:     "nested actions",
			input:    "`{{ range .Items }}{{ .Name }}{{ end }}`",
			expected: "<p><code>{{ range .Items }}{{ .Name }}{{ end }}</code></p>",
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

func TestEscapedTemplates(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "escaped opening braces",
			input:    "`{{\"{{\"}}`",
			expected: "<p><code>{{\"{{\"}}</code></p>",
		},
		{
			name:     "escaped closing braces",
			input:    "`{{\"}}\"}}`",
			expected: "<p><code>{{\"}}\"}}</code></p>",
		},
		{
			name:     "mixed escaped and regular actions",
			input:    "`{{\"{{\"}}` and `{{ .Name }}`",
			expected: "<p><code>{{\"{{\"}}</code> and <code>{{ .Name }}</code></p>",
		},
		{
			name:     "escaped actions in code block",
			input:    "```\n{{\"{{\"}} and {{\"}}\"}}\n```",
			expected: "<pre><code>{{\"{{\"}} and {{\"}}\"}}\n</code></pre>",
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
