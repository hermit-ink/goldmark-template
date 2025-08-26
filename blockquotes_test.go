package goldmarktemplate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

func TestBasicBlockquotes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple blockquote with template",
			input:    "> {{ .Quote }}",
			expected: "<blockquote>\n<p>{{ .Quote }}</p>\n</blockquote>",
		},
		{
			name:     "blockquote with mixed content",
			input:    "> {{ .Author }} said this",
			expected: "<blockquote>\n<p>{{ .Author }} said this</p>\n</blockquote>",
		},
		{
			name:     "blockquote with emphasis",
			input:    "> *{{ .ImportantQuote }}*",
			expected: "<blockquote>\n<p><em>{{ .ImportantQuote }}</em></p>\n</blockquote>",
		},
		{
			name:     "blockquote with strong",
			input:    "> **{{ .VeryImportant }}**",
			expected: "<blockquote>\n<p><strong>{{ .VeryImportant }}</strong></p>\n</blockquote>",
		},
		{
			name:     "blockquote with code span",
			input:    "> Use `{{ .Command }}` to run",
			expected: "<blockquote>\n<p>Use <code>{{ .Command }}</code> to run</p>\n</blockquote>",
		},
		{
			name:     "blockquote with link",
			input:    "> Read more at [{{ .LinkText }}]({{ .URL }})",
			expected: "<blockquote>\n<p>Read more at <a href=\"{{ .URL }}\">{{ .LinkText }}</a></p>\n</blockquote>",
		},
		{
			name:     "blockquote with complex template",
			input:    "> {{ if .ShowQuote }}{{ .Quote }}{{ else }}No quote{{ end }}",
			expected: "<blockquote>\n<p>{{ if .ShowQuote }}{{ .Quote }}{{ else }}No quote{{ end }}</p>\n</blockquote>",
		},
		{
			name:     "blockquote with multiple templates",
			input:    "> {{ .FirstPart }} and {{ .SecondPart }}",
			expected: "<blockquote>\n<p>{{ .FirstPart }} and {{ .SecondPart }}</p>\n</blockquote>",
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

func TestNestedBlockquotes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "nested blockquote with template",
			input:    "> > {{ .NestedQuote }}",
			expected: "<blockquote>\n<blockquote>\n<p>{{ .NestedQuote }}</p>\n</blockquote>\n</blockquote>",
		},
		{
			name:     "triple nested blockquote",
			input:    "> > > {{ .DeeplyNested }}",
			expected: "<blockquote>\n<blockquote>\n<blockquote>\n<p>{{ .DeeplyNested }}</p>\n</blockquote>\n</blockquote>\n</blockquote>",
		},
		{
			name:     "mixed nesting levels",
			input:    "> {{ .Level1 }}\n> > {{ .Level2 }}\n> {{ .BackToLevel1 }}",
			expected: "<blockquote>\n<p>{{ .Level1 }}</p>\n<blockquote>\n<p>{{ .Level2 }}\n{{ .BackToLevel1 }}</p>\n</blockquote>\n</blockquote>",
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

func TestMultiParagraphBlockquotes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "blockquote with multiple paragraphs",
			input:    "> {{ .FirstParagraph }}\n>\n> {{ .SecondParagraph }}",
			expected: "<blockquote>\n<p>{{ .FirstParagraph }}</p>\n<p>{{ .SecondParagraph }}</p>\n</blockquote>",
		},
		{
			name:     "blockquote with lazy continuation",
			input:    "> {{ .StartOfQuote }}\n{{ .Continuation }}",
			expected: "<blockquote>\n<p>{{ .StartOfQuote }}\n{{ .Continuation }}</p>\n</blockquote>",
		},
		{
			name:     "complex multi-paragraph blockquote",
			input:    "> **{{ .Author }}** said:\n>\n> {{ .Quote }}\n>\n> *{{ .Attribution }}*",
			expected: "<blockquote>\n<p><strong>{{ .Author }}</strong> said:</p>\n<p>{{ .Quote }}</p>\n<p><em>{{ .Attribution }}</em></p>\n</blockquote>",
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

func TestBlockquotesWithOtherElements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "blockquote with list",
			input:    "> {{ .Intro }}:\n>\n> - {{ .Item1 }}\n> - {{ .Item2 }}",
			expected: "<blockquote>\n<p>{{ .Intro }}:</p>\n<ul>\n<li>{{ .Item1 }}</li>\n<li>{{ .Item2 }}</li>\n</ul>\n</blockquote>",
		},
		{
			name:     "blockquote with code block",
			input:    "> Example:\n>\n> ```\n> {{ .CodeExample }}\n> ```",
			expected: "<blockquote>\n<p>Example:</p>\n<pre><code>{{ .CodeExample }}\n</code></pre>\n</blockquote>",
		},
		{
			name:     "blockquote with heading",
			input:    "> ## {{ .HeadingText }}\n>\n> {{ .Content }}",
			expected: "<blockquote>\n<h2>{{ .HeadingText }}</h2>\n<p>{{ .Content }}</p>\n</blockquote>",
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

func TestBlockquotesWithSpecialCharacters(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "blockquote with HTML chars",
			input:    "> {{ .Quote }} > {{ .Author }}",
			expected: "<blockquote>\n<p>{{ .Quote }} &gt; {{ .Author }}</p>\n</blockquote>",
		},
		{
			name:     "blockquote with quotes",
			input:    "> \"{{ .DoubleQuote }}\" and '{{ .SingleQuote }}'",
			expected: "<blockquote>\n<p>&quot;{{ .DoubleQuote }}&quot; and '{{ .SingleQuote }}'</p>\n</blockquote>",
		},
		{
			name:     "blockquote with ampersand",
			input:    "> {{ .Company }} & {{ .Partner }}",
			expected: "<blockquote>\n<p>{{ .Company }} &amp; {{ .Partner }}</p>\n</blockquote>",
		},
		{
			name:     "blockquote attribution format",
			input:    "> {{ .Quote }}\n>\n> — {{ .Author }}",
			expected: "<blockquote>\n<p>{{ .Quote }}</p>\n<p>— {{ .Author }}</p>\n</blockquote>",
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