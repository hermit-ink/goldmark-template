package goldmarktemplate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

func TestUnorderedLists(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple unordered list with templates",
			input:    "- {{ .Item1 }}\n- {{ .Item2 }}",
			expected: "<ul>\n<li>{{ .Item1 }}</li>\n<li>{{ .Item2 }}</li>\n</ul>",
		},
		{
			name:     "unordered list with mixed content",
			input:    "- First: {{ .First }}\n- Second: {{ .Second }}",
			expected: "<ul>\n<li>First: {{ .First }}</li>\n<li>Second: {{ .Second }}</li>\n</ul>",
		},
		{
			name:     "unordered list with emphasis",
			input:    "- *{{ .Important }}*\n- **{{ .VeryImportant }}**",
			expected: "<ul>\n<li><em>{{ .Important }}</em></li>\n<li><strong>{{ .VeryImportant }}</strong></li>\n</ul>",
		},
		{
			name:     "unordered list with code spans",
			input:    "- Use `{{ .Function }}`\n- Call `{{ .Method }}`",
			expected: "<ul>\n<li>Use <code>{{ .Function }}</code></li>\n<li>Call <code>{{ .Method }}</code></li>\n</ul>",
		},
		{
			name:     "unordered list with links",
			input:    "- [{{ .LinkText }}]({{ .URL }})\n- [Documentation]({{ .DocsURL }})",
			expected: "<ul>\n<li><a href=\"{{ .URL }}\">{{ .LinkText }}</a></li>\n<li><a href=\"{{ .DocsURL }}\">Documentation</a></li>\n</ul>",
		},
		{
			name:     "unordered list with complex templates",
			input:    "- {{ if .ShowFirst }}{{ .First }}{{ end }}\n- {{ range .Items }}{{ .Name }} {{ end }}",
			expected: "<ul>\n<li>{{ if .ShowFirst }}{{ .First }}{{ end }}</li>\n<li>{{ range .Items }}{{ .Name }} {{ end }}</li>\n</ul>",
		},
		{
			name:     "unordered list with plus marker",
			input:    "+ {{ .ItemA }}\n+ {{ .ItemB }}",
			expected: "<ul>\n<li>{{ .ItemA }}</li>\n<li>{{ .ItemB }}</li>\n</ul>",
		},
		{
			name:     "unordered list with asterisk marker",
			input:    "* {{ .ItemX }}\n* {{ .ItemY }}",
			expected: "<ul>\n<li>{{ .ItemX }}</li>\n<li>{{ .ItemY }}</li>\n</ul>",
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

func TestOrderedLists(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple ordered list with templates",
			input:    "1. {{ .First }}\n2. {{ .Second }}",
			expected: "<ol>\n<li>{{ .First }}</li>\n<li>{{ .Second }}</li>\n</ol>",
		},
		{
			name:     "ordered list with mixed content",
			input:    "1. Step: {{ .Step1 }}\n2. Next: {{ .Step2 }}",
			expected: "<ol>\n<li>Step: {{ .Step1 }}</li>\n<li>Next: {{ .Step2 }}</li>\n</ol>",
		},
		{
			name:     "ordered list with emphasis",
			input:    "1. *{{ .Action1 }}*\n2. **{{ .Action2 }}**",
			expected: "<ol>\n<li><em>{{ .Action1 }}</em></li>\n<li><strong>{{ .Action2 }}</strong></li>\n</ol>",
		},
		{
			name:     "ordered list with custom start",
			input:    "5. {{ .FifthItem }}\n6. {{ .SixthItem }}",
			expected: "<ol start=\"5\">\n<li>{{ .FifthItem }}</li>\n<li>{{ .SixthItem }}</li>\n</ol>",
		},
		{
			name:     "ordered list with parentheses marker",
			input:    "1) {{ .FirstItem }}\n2) {{ .SecondItem }}",
			expected: "<ol>\n<li>{{ .FirstItem }}</li>\n<li>{{ .SecondItem }}</li>\n</ol>",
		},
		{
			name:     "ordered list with code spans",
			input:    "1. Run `{{ .Command1 }}`\n2. Execute `{{ .Command2 }}`",
			expected: "<ol>\n<li>Run <code>{{ .Command1 }}</code></li>\n<li>Execute <code>{{ .Command2 }}</code></li>\n</ol>",
		},
		{
			name:     "ordered list with complex templates",
			input:    "1. {{ printf \"Step %d: %s\" .Number .Description }}\n2. {{ .Instructions | upper }}",
			expected: "<ol>\n<li>{{ printf \"Step %d: %s\" .Number .Description }}</li>\n<li>{{ .Instructions | upper }}</li>\n</ol>",
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

func TestNestedLists(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "nested unordered list",
			input:    "- {{ .Parent }}\n  - {{ .Child }}",
			expected: "<ul>\n<li>{{ .Parent }}\n<ul>\n<li>{{ .Child }}</li>\n</ul>\n</li>\n</ul>",
		},
		{
			name:     "nested ordered list",
			input:    "1. {{ .MainStep }}\n   1. {{ .SubStep }}",
			expected: "<ol>\n<li>{{ .MainStep }}\n<ol>\n<li>{{ .SubStep }}</li>\n</ol>\n</li>\n</ol>",
		},
		{
			name:     "mixed nested lists",
			input:    "1. {{ .OrderedItem }}\n   - {{ .UnorderedSub }}",
			expected: "<ol>\n<li>{{ .OrderedItem }}\n<ul>\n<li>{{ .UnorderedSub }}</li>\n</ul>\n</li>\n</ol>",
		},
		{
			name:     "deeply nested list",
			input:    "- {{ .Level1 }}\n  - {{ .Level2 }}\n    - {{ .Level3 }}",
			expected: "<ul>\n<li>{{ .Level1 }}\n<ul>\n<li>{{ .Level2 }}\n<ul>\n<li>{{ .Level3 }}</li>\n</ul>\n</li>\n</ul>\n</li>\n</ul>",
		},
		{
			name:     "multiple nested items",
			input:    "1. {{ .First }}\n   - {{ .FirstSub1 }}\n   - {{ .FirstSub2 }}\n2. {{ .Second }}",
			expected: "<ol>\n<li>{{ .First }}\n<ul>\n<li>{{ .FirstSub1 }}</li>\n<li>{{ .FirstSub2 }}</li>\n</ul>\n</li>\n<li>{{ .Second }}</li>\n</ol>",
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

func TestListsWithSpecialCharacters(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "list with HTML chars",
			input:    "- {{ .Item }} > {{ .Comparison }}\n- {{ .Other }} & more",
			expected: "<ul>\n<li>{{ .Item }} &gt; {{ .Comparison }}</li>\n<li>{{ .Other }} &amp; more</li>\n</ul>",
		},
		{
			name:     "list with quotes",
			input:    "1. \"{{ .Quote }}\"\n2. '{{ .SingleQuote }}'",
			expected: "<ol>\n<li>&quot;{{ .Quote }}&quot;</li>\n<li>'{{ .SingleQuote }}'</li>\n</ol>",
		},
		{
			name:     "list with multi-line items",
			input:    "- {{ .FirstLine }}\n  {{ .SecondLine }}\n- {{ .NextItem }}",
			expected: "<ul>\n<li>{{ .FirstLine }}\n{{ .SecondLine }}</li>\n<li>{{ .NextItem }}</li>\n</ul>",
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