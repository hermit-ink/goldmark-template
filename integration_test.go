package goldmarktemplate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func TestWithGFMExtension(t *testing.T) {
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
				goldmark.WithParser(NewParser()),
				goldmark.WithExtensions(
					extension.GFM,
					NewExtension(), // Our extension should work with GFM
				),
				goldmark.WithRendererOptions(
					html.WithUnsafe(),
					html.WithXHTML(),
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
		goldmark.WithParser(NewParser()),
		goldmark.WithExtensions(NewExtension()),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
			html.WithXHTML(),
		),
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
		goldmark.WithParser(NewParser()),
		goldmark.WithExtensions(NewExtension()),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
			html.WithXHTML(),
		),
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