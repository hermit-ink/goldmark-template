package goldmarktemplate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

func TestBackslashEscapes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "escaped asterisk before template",
			input:    "\\*{{ .NotEmphasis }}",
			expected: "<p>*{{ .NotEmphasis }}</p>",
		},
		{
			name:     "template with escaped asterisk after",
			input:    "{{ .Value }}\\*not emphasis\\*",
			expected: "<p>{{ .Value }}*not emphasis*</p>",
		},
		{
			name:     "escaped brackets around template",
			input:    "\\[{{ .NotLink }}\\]",
			expected: "<p>[{{ .NotLink }}]</p>",
		},
		{
			name:     "escaped underscore before template",
			input:    "\\_{{ .NotEmphasis }}\\_",
			expected: "<p>_{{ .NotEmphasis }}_</p>",
		},
		{
			name:     "escaped backtick before template",
			input:    "\\`{{ .NotCode }}\\`",
			expected: "<p>`{{ .NotCode }}`</p>",
		},
		{
			name:     "escaped hash before template",
			input:    "\\# {{ .NotHeading }}",
			expected: "<p># {{ .NotHeading }}</p>",
		},
		{
			name:     "escaped greater than before template",
			input:    "\\> {{ .NotBlockquote }}",
			expected: "<p>&gt; {{ .NotBlockquote }}</p>",
		},
		{
			name:     "escaped plus before template",
			input:    "\\+ {{ .NotList }}",
			expected: "<p>+ {{ .NotList }}</p>",
		},
		{
			name:     "escaped minus before template",
			input:    "\\- {{ .NotList }}",
			expected: "<p>- {{ .NotList }}</p>",
		},
		{
			name:     "escaped period after number and template",
			input:    "1\\. {{ .NotList }}",
			expected: "<p>1. {{ .NotList }}</p>",
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

func TestBackslashEscapesInTemplates(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "template containing escaped characters",
			input:    "{{ printf \"\\*%s\\*\" .Value }}",
			expected: "<p>{{ printf \"\\*%s\\*\" .Value }}</p>",
		},
		{
			name:     "template with escaped quotes",
			input:    "{{ printf \"\\\"Hello %s\\\"\" .Name }}",
			expected: "<p>{{ printf \"\\\"Hello %s\\\"\" .Name }}</p>",
		},
		{
			name:     "template with escaped backslash",
			input:    "{{ .Path | replace \"\\\\\" \"/\" }}",
			expected: "<p>{{ .Path | replace \"\\\\\" \"/\" }}</p>",
		},
		{
			name:     "mixed escaped and template content",
			input:    "\\*{{ .Important }}\\* and \\`{{ .Code }}\\`",
			expected: "<p>*{{ .Important }}* and `{{ .Code }}`</p>",
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

func TestBackslashEscapesWithMarkdownElements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "escaped emphasis markers in emphasis",
			input:    "*{{ .Text }} \\*not emphasized\\**",
			expected: "<p><em>{{ .Text }} *not emphasized*</em></p>",
		},
		{
			name:     "escaped link brackets in link text",
			input:    "[{{ .Text }} \\[not a link\\]]({{ .URL }})",
			expected: "<p><a href=\"{{ .URL }}\">{{ .Text }} [not a link]</a></p>",
		},
		{
			name:     "escaped code span backticks",
			input:    "`{{ .Code }}` and \\`not code\\`",
			expected: "<p><code>{{ .Code }}</code> and `not code`</p>",
		},
		{
			name:     "escaped in blockquote",
			input:    "> {{ .Quote }} \\> not a nested quote",
			expected: "<blockquote>\n<p>{{ .Quote }} &gt; not a nested quote</p>\n</blockquote>",
		},
		{
			name:     "escaped in list",
			input:    "- {{ .Item }}\n\\- not a list item",
			expected: "<ul>\n<li>{{ .Item }}\n- not a list item</li>\n</ul>",
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

func TestBackslashAtEndOfLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "backslash at end creates hard line break",
			input:    "{{ .FirstLine }}\\\n{{ .SecondLine }}",
			expected: "<p>{{ .FirstLine }}<br />\n{{ .SecondLine }}</p>",
		},
		{
			name:     "template after backslash line break",
			input:    "First line\\\n{{ .Template }} here",
			expected: "<p>First line<br />\n{{ .Template }} here</p>",
		},
		{
			name:     "template before backslash line break",
			input:    "{{ .Template }} here\\\nNext line",
			expected: "<p>{{ .Template }} here<br />\nNext line</p>",
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