package goldmarktemplate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

func TestBasicEmphasis(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "italic with template using asterisks",
			input:    "*{{ .Text }}*",
			expected: "<p><em>{{ .Text }}</em></p>",
		},
		{
			name:     "italic with template using underscores",
			input:    "_{{ .Text }}_",
			expected: "<p><em>{{ .Text }}</em></p>",
		},
		{
			name:     "bold with template using asterisks",
			input:    "**{{ .Text }}**",
			expected: "<p><strong>{{ .Text }}</strong></p>",
		},
		{
			name:     "bold with template using underscores",
			input:    "__{{ .Text }}__",
			expected: "<p><strong>{{ .Text }}</strong></p>",
		},
		{
			name:     "italic with mixed content",
			input:    "*Welcome {{ .Name }}*",
			expected: "<p><em>Welcome {{ .Name }}</em></p>",
		},
		{
			name:     "bold with mixed content",
			input:    "**Hello {{ .User }}**",
			expected: "<p><strong>Hello {{ .User }}</strong></p>",
		},
		{
			name:     "emphasis with multiple templates",
			input:    "*{{ .First }} and {{ .Second }}*",
			expected: "<p><em>{{ .First }} and {{ .Second }}</em></p>",
		},
		{
			name:     "emphasis with complex template",
			input:    "*{{ if .Show }}{{ .Content }}{{ end }}*",
			expected: "<p><em>{{ if .Show }}{{ .Content }}{{ end }}</em></p>",
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

func TestNestedEmphasis(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "nested emphasis with template",
			input:    "***{{ .Text }}***",
			expected: "<p><em><strong>{{ .Text }}</strong></em></p>",
		},
		{
			name:     "mixed nested emphasis with template",
			input:    "**_{{ .Text }}_**",
			expected: "<p><strong><em>{{ .Text }}</em></strong></p>",
		},
		{
			name:     "reverse nested emphasis with template",
			input:    "_**{{ .Text }}**_",
			expected: "<p><em><strong>{{ .Text }}</strong></em></p>",
		},
		{
			name:     "complex nested emphasis",
			input:    "***{{ .Important }} text***",
			expected: "<p><em><strong>{{ .Important }} text</strong></em></p>",
		},
		{
			name:     "nested emphasis with multiple templates",
			input:    "***{{ .First }} and {{ .Second }}***",
			expected: "<p><em><strong>{{ .First }} and {{ .Second }}</strong></em></p>",
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

func TestEmphasisWithOtherElements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "emphasis with code span",
			input:    "*Use `{{ .Function }}` here*",
			expected: "<p><em>Use <code>{{ .Function }}</code> here</em></p>",
		},
		{
			name:     "emphasis with link",
			input:    "*Visit [{{ .SiteName }}]({{ .URL }})*",
			expected: "<p><em>Visit <a href=\"{{ .URL }}\">{{ .SiteName }}</a></em></p>",
		},
		{
			name:     "code span with emphasis inside",
			input:    "`*{{ .Variable }}*`",
			expected: "<p><code>*{{ .Variable }}*</code></p>",
		},
		{
			name:     "link with emphasis",
			input:    "[*{{ .LinkText }}*]({{ .URL }})",
			expected: "<p><a href=\"{{ .URL }}\"><em>{{ .LinkText }}</em></a></p>",
		},
		{
			name:     "mixed emphasis and regular text",
			input:    "This is *{{ .Important }}* and this is **{{ .VeryImportant }}**.",
			expected: "<p>This is <em>{{ .Important }}</em> and this is <strong>{{ .VeryImportant }}</strong>.</p>",
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

func TestEmphasisWithSpecialCharacters(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "emphasis with HTML chars",
			input:    "*{{ .Text }} > {{ .Other }}*",
			expected: "<p><em>{{ .Text }} &gt; {{ .Other }}</em></p>",
		},
		{
			name:     "emphasis with quotes",
			input:    "*\"{{ .Quote }}\"*",
			expected: "<p><em>&quot;{{ .Quote }}&quot;</em></p>",
		},
		{
			name:     "emphasis with ampersand",
			input:    "**{{ .Company }} & {{ .Partner }}**",
			expected: "<p><strong>{{ .Company }} &amp; {{ .Partner }}</strong></p>",
		},
		{
			name:     "emphasis with asterisks in template",
			input:    "*{{ .Value }}* and *more text*",
			expected: "<p><em>{{ .Value }}</em> and <em>more text</em></p>",
		},
		{
			name:     "emphasis with underscores in template",
			input:    "_{{ .Variable }}_name_",
			expected: "<p>_{{ .Variable }}<em>name</em></p>",
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

func TestEmphasisDelimiterHandling(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "templates with internal asterisks",
			input:    "{{ .Value }} * {{ .Other }}",
			expected: "<p>{{ .Value }} * {{ .Other }}</p>",
		},
		{
			name:     "templates with internal underscores",
			input:    "{{ .User_Name }} and {{ .Other_Value }}",
			expected: "<p>{{ .User_Name }} and {{ .Other_Value }}</p>",
		},
		{
			name:     "unmatched emphasis markers",
			input:    "*{{ .Start }} but no end",
			expected: "<p>*{{ .Start }} but no end</p>",
		},
		{
			name:     "emphasis around template boundaries",
			input:    "before *{{ .Text }}* after",
			expected: "<p>before <em>{{ .Text }}</em> after</p>",
		},
		{
			name:     "asterisk with spaces creates list",
			input:    "* {{ .SpacedTemplate }} *",
			expected: "<ul>\n<li>{{ .SpacedTemplate }} *</li>\n</ul>",
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