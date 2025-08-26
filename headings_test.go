package goldmarktemplate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

func TestATXHeadings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "h1 with template",
			input:    "# {{ .Title }}",
			expected: "<h1>{{ .Title }}</h1>",
		},
		{
			name:     "h2 with template",
			input:    "## {{ .Subtitle }}",
			expected: "<h2>{{ .Subtitle }}</h2>",
		},
		{
			name:     "h3 with template",
			input:    "### {{ .Section }}",
			expected: "<h3>{{ .Section }}</h3>",
		},
		{
			name:     "h6 with template",
			input:    "###### {{ .SmallHeading }}",
			expected: "<h6>{{ .SmallHeading }}</h6>",
		},
		{
			name:     "heading with mixed content",
			input:    "# Welcome to {{ .SiteName }}",
			expected: "<h1>Welcome to {{ .SiteName }}</h1>",
		},
		{
			name:     "heading with multiple templates",
			input:    "## {{ .Category }}: {{ .Title }}",
			expected: "<h2>{{ .Category }}: {{ .Title }}</h2>",
		},
		{
			name:     "heading with template and emphasis",
			input:    "# {{ .Title }} - *{{ .Tagline }}*",
			expected: "<h1>{{ .Title }} - <em>{{ .Tagline }}</em></h1>",
		},
		{
			name:     "heading with template and strong",
			input:    "## **{{ .Important }}** Notice",
			expected: "<h2><strong>{{ .Important }}</strong> Notice</h2>",
		},
		{
			name:     "heading with complex template",
			input:    "# {{ if .Title }}{{ .Title }}{{ else }}Untitled{{ end }}",
			expected: "<h1>{{ if .Title }}{{ .Title }}{{ else }}Untitled{{ end }}</h1>",
		},
		{
			name:     "heading with template containing spaces",
			input:    "## {{ .User.Name }}",
			expected: "<h2>{{ .User.Name }}</h2>",
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

func TestSetextHeadings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "setext h1 with template",
			input:    "{{ .Title }}\n=======",
			expected: "<h1>{{ .Title }}</h1>",
		},
		{
			name:     "setext h2 with template",
			input:    "{{ .Subtitle }}\n-------",
			expected: "<h2>{{ .Subtitle }}</h2>",
		},
		{
			name:     "setext h1 with mixed content",
			input:    "Welcome to {{ .SiteName }}\n==================",
			expected: "<h1>Welcome to {{ .SiteName }}</h1>",
		},
		{
			name:     "setext h2 with emphasis",
			input:    "*{{ .Important }}* Update\n--------------------",
			expected: "<h2><em>{{ .Important }}</em> Update</h2>",
		},
		{
			name:     "setext h1 with multiple templates",
			input:    "{{ .FirstName }} {{ .LastName }}\n======================",
			expected: "<h1>{{ .FirstName }} {{ .LastName }}</h1>",
		},
		{
			name:     "setext h2 with complex template",
			input:    "{{ range .Items }}{{ .Name }} {{ end }}\n------------------------",
			expected: "<h2>{{ range .Items }}{{ .Name }} {{ end }}</h2>",
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

func TestHeadingsWithSpecialCharacters(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "heading with HTML chars in template",
			input:    "# {{ .Title }} > {{ .Subtitle }}",
			expected: "<h1>{{ .Title }} &gt; {{ .Subtitle }}</h1>",
		},
		{
			name:     "heading with quotes in template",
			input:    `## "{{ .Quote }}"`,
			expected: "<h2>&quot;{{ .Quote }}&quot;</h2>",
		},
		{
			name:     "heading with ampersand",
			input:    "# {{ .Company }} & Associates",
			expected: "<h1>{{ .Company }} &amp; Associates</h1>",
		},
		{
			name:     "heading with code span",
			input:    "# Using `{{ .Function }}` in Go",
			expected: "<h1>Using <code>{{ .Function }}</code> in Go</h1>",
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