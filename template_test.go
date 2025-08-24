package goldmarktemplate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func TestTemplateExtension(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Basic template directives in code spans
		{
			name:     "simple variable in code span",
			input:    "`{{ .Name }}`",
			expected: "<p><code>{{ .Name }}</code></p>",
		},
		{
			name:     "template with comparison in code span",
			input:    "`{{ if .Value > 5 }}high{{ end }}`",
			expected: "<p><code>{{ if .Value > 5 }}high{{ end }}</code></p>",
		},
		{
			name:     "template with less than in code span",
			input:    "`{{ if .Value < 10 }}low{{ end }}`",
			expected: "<p><code>{{ if .Value < 10 }}low{{ end }}</code></p>",
		},
		{
			name:     "template with ampersand in code span",
			input:    "`{{ .Title & .Subtitle }}`",
			expected: "<p><code>{{ .Title & .Subtitle }}</code></p>",
		},

		// Escaped template cases
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
			name:     "mixed escaped and regular templates",
			input:    "`{{\"{{\"}}` and `{{ .Name }}`",
			expected: "<p><code>{{\"{{\"}}</code> and <code>{{ .Name }}</code></p>",
		},

		// Templates in code blocks
		{
			name: "simple template in code block",
			input: "```\n{{ .Name }}\n```",
			expected: "<pre><code>{{ .Name }}\n</code></pre>\n",
		},
		{
			name: "template with HTML chars in code block",
			input: "```\n{{ if .Value > 5 }}<div>High</div>{{ end }}\n```",
			expected: "<pre><code>{{ if .Value > 5 }}&lt;div&gt;High&lt;/div&gt;{{ end }}\n</code></pre>\n",
		},
		{
			name: "escaped templates in code block",
			input: "```\n{{\"{{\"}} and {{\"}}\"}}\n```",
			expected: "<pre><code>{{\"{{\"}} and {{\"}}\"}}\n</code></pre>\n",
		},

		// Nested templates
		{
			name:     "nested template directives",
			input:    "`{{ range .Items }}{{ .Name }}{{ end }}`",
			expected: "<p><code>{{ range .Items }}{{ .Name }}{{ end }}</code></p>",
		},
		{
			name:     "template with string containing braces",
			input:    "`{{ printf \"Value: %d\" .Count }}`",
			expected: "<p><code>{{ printf \"Value: %d\" .Count }}</code></p>",
		},
		{
			name:     "template with string containing }}",
			input:    "`{{ \"a string with }}\" }}`",
			expected: "<p><code>{{ \"a string with }}\" }}</code></p>",
		},

		// Edge cases
		{
			name:     "incomplete template directive (later step will find the error)",
			input:    "`{{ .Name` without closing",
			expected: "<p><code>{{ .Name</code> without closing</p>",
		},
		{
			name:     "template with backticks in string",
			input:    "`{{ `raw string` }}`",
			expected: "<p><code>{{ </code>raw string<code> }}</code></p>",
		},
		{
			name:     "multiple templates in one code span",
			input:    "`{{ .First }} and {{ .Second }}`",
			expected: "<p><code>{{ .First }} and {{ .Second }}</code></p>",
		},

		// Regular markdown should still work
		{
			name:     "non-template content with special chars",
			input:    "`<div>&amp;</div>`",
			expected: "<p><code>&lt;div&gt;&amp;amp;&lt;/div&gt;</code></p>",
		},
		{
			name:     "mixed content",
			input:    "Text with `{{ .Code }}` and `regular code`.",
			expected: "<p>Text with <code>{{ .Code }}</code> and <code>regular code</code>.</p>\n",
		},

		// Fenced code blocks with language
		{
			name: "template in go code block",
			input: "```go\nfmt.Println({{ .Value }})\n```",
			expected: "<pre><code class=\"language-go\">fmt.Println({{ .Value }})\n</code></pre>\n",
		},
		{
			name: "template in html code block",
			input: "```html\n<div>{{ .Content }}</div>\n```",
			expected: "<pre><code class=\"language-html\">&lt;div&gt;{{ .Content }}&lt;/div&gt;\n</code></pre>\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := goldmark.New(
				goldmark.WithExtensions(
					NewTemplatedHTMLExtension(),
				),
			)

			var buf bytes.Buffer
			err := md.Convert([]byte(tt.input), &buf)
			if err != nil {
				t.Fatalf("Failed to convert markdown: %v", err)
			}

			got := buf.String()
			// Normalize whitespace for comparison
			got = strings.TrimSpace(got)
			expected := strings.TrimSpace(tt.expected)

			if got != expected {
				t.Errorf("Output mismatch\nInput:    %q\nExpected: %q\nGot:      %q", tt.input, expected, got)
			}
		})
	}
}

func TestTemplateExtensionWithOtherExtensions(t *testing.T) {
	// Test that our extension works well with other goldmark extensions
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "template with GFM table",
			input:    "| Col |\n|-----|\n| `{{ .Value }}` |",
			expected: "<table>\n<thead>\n<tr>\n<th>Col</th>\n</tr>\n</thead>\n<tbody>\n<tr>\n<td><code>{{ .Value }}</code></td>\n</tr>\n</tbody>\n</table>",
		},
		{
			name:     "template with strikethrough",
			input:    "~~old~~ `{{ .New }}`",
			expected: "<p><del>old</del> <code>{{ .New }}</code></p>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := goldmark.New(
				goldmark.WithExtensions(
					extension.GFM,
					NewTemplatedHTMLExtension(), // Our extension should work with GFM
				),
			)

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

func BenchmarkTemplateExtension(b *testing.B) {
	md := goldmark.New(
		goldmark.WithExtensions(NewTemplatedHTMLExtension()),
	)

	input := []byte("Text with `{{ .Code }}` and `{{ if .Value > 5 }}high{{ end }}` templates.")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		_ = md.Convert(input, &buf)
	}
}

func BenchmarkTemplateExtensionComplex(b *testing.B) {
	md := goldmark.New(
		goldmark.WithExtensions(NewTemplatedHTMLExtension()),
	)

	input := []byte(`
# Template Examples

Here's a simple variable: ` + "`{{ .Name }}`" + `

## Code Block

` + "```go" + `
func main() {
    fmt.Println("{{ .Message }}")
    if {{ .Condition }} {
        {{ range .Items }}
        processItem({{ . }})
        {{ end }}
    }
}
` + "```" + `

## Escaped Templates

Show literal braces: ` + "`{{\"{{\"}}`" + ` and ` + "`{{\"}}\"}}`" + `
`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		_ = md.Convert(input, &buf)
	}
}

func TestAttributeAndURLContexts(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Link URLs with templates
		{
			name:     "link with template in URL",
			input:    "[Link]({{ .URL }})",
			expected: "<p><a href=\"{{ .URL }}\">Link</a></p>",
		},
		{
			name:     "link with template containing special chars in URL",
			input:    "[Link]({{ .BaseURL }}/path?param={{ .Value }})",
			expected: "<p><a href=\"{{ .BaseURL }}/path?param={{ .Value }}\">Link</a></p>",
		},
		{
			name:     "link with template in title and URL",
			input:    "[{{ .Title }}]({{ .URL }})",
			expected: "<p><a href=\"{{ .URL }}\">{{ .Title }}</a></p>",
		},
		{
			name:     "reference link with template in URL",
			input:    "[Link][ref]\n\n[ref]: {{ .URL }}",
			expected: "<p><a href=\"{{ .URL }}\">Link</a></p>",
		},

		// Image sources with templates
		{
			name:     "image with template in src",
			input:    "![Alt]({{ .ImageURL }})",
			expected: "<p><img src=\"{{ .ImageURL }}\" alt=\"Alt\" /></p>",
		},
		{
			name:     "image with template in alt and src",
			input:    "![{{ .AltText }}]({{ .ImageURL }})",
			expected: "<p><img src=\"{{ .ImageURL }}\" alt=\"{{ .AltText }}\" /></p>",
		},
		{
			name:     "image with template containing quotes in alt",
			input:    "![{{ .Title | quote }}](image.jpg)",
			expected: "<p><img src=\"image.jpg\" alt=\"{{ .Title | quote }}\" /></p>",
		},

		// Raw HTML attributes with templates
		{
			name:     "raw HTML with template in attribute",
			input:    "<div id=\"{{ .ID }}\">Content</div>",
			expected: "<div id=\"{{ .ID }}\">Content</div>",
		},
		{
			name:     "raw HTML with template in class attribute",
			input:    "<span class=\"{{ .CSSClass }}\">Text</span>",
			expected: "<p><span class=\"{{ .CSSClass }}\">Text</span></p>",
		},
		{
			name:     "raw HTML with template containing quotes in attribute",
			input:    "<div data-value=\"{{ .Value | quote }}\">Content</div>",
			expected: "<div data-value=\"{{ .Value | quote }}\">Content</div>",
		},
		{
			name:     "raw HTML with multiple template attributes",
			input:    "<a href=\"{{ .URL }}\" title=\"{{ .Title }}\" class=\"{{ .Class }}\">Link</a>",
			expected: "<p><a href=\"{{ .URL }}\" title=\"{{ .Title }}\" class=\"{{ .Class }}\">Link</a></p>",
		},

		// Templates in HTML content (not attributes)
		{
			name:     "template in raw HTML content",
			input:    "<div>{{ .Content }}</div>",
			expected: "<div>{{ .Content }}</div>",
		},
		{
			name:     "mixed template and HTML with special chars",
			input:    "<p>Value: {{ if .Value > 0 }}{{ .Value }}{{ else }}N/A{{ end }}</p>",
			expected: "<p>Value: {{ if .Value > 0 }}{{ .Value }}{{ else }}N/A{{ end }}</p>",
		},

		// Edge cases with escaping
		{
			name:     "link URL with escaped template",
			input:    "[Link]({{\"{{\"}} .URL {{\"}}\"}})",
			expected: "<p><a href=\"{{\"{{\"}} .URL {{\"}}\"}}\">Link</a></p>",
		},
		{
			name:     "image alt with HTML chars and templates",
			input:    "![{{ .Title }} > {{ .Subtitle }}](image.jpg)",
			expected: "<p><img src=\"image.jpg\" alt=\"{{ .Title }} > {{ .Subtitle }}\" /></p>",
		},
		{
			name:     "raw HTML attribute with nested templates",
			input:    "<div data-config=\"{{ range .Items }}{{ .Name }},{{ end }}\">Content</div>",
			expected: "<div data-config=\"{{ range .Items }}{{ .Name }},{{ end }}\">Content</div>",
		},
	}

	md := goldmark.New(
		goldmark.WithParser(NewTemplatedParser()),
		goldmark.WithExtensions(NewTemplatedHTMLExtension()),
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

func TestEdgeCases(t *testing.T) {
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
			name:     "template with only spaces",
			input:    "`{{   }}`",
			expected: "<p><code>{{   }}</code></p>",
		},
		{
			name:     "unmatched opening braces",
			input:    "`{{ {{ .Name }}`",
			expected: "<p><code>{{ {{ .Name }}</code></p>",
		},
		{
			name:     "template at start and end",
			input:    "`{{ .Start }}content{{ .End }}`",
			expected: "<p><code>{{ .Start }}content{{ .End }}</code></p>",
		},
		{
			name:     "deeply nested templates",
			input:    "`{{ if .A }}{{ if .B }}{{ .C }}{{ end }}{{ end }}`",
			expected: "<p><code>{{ if .A }}{{ if .B }}{{ .C }}{{ end }}{{ end }}</code></p>",
		},
	}

	md := goldmark.New(
		goldmark.WithExtensions(NewTemplatedHTMLExtension()),
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
