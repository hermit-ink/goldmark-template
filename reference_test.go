package goldmarktemplate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

func TestReferenceLinks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "basic reference link",
			input: `[Example][1]

[1]: https://example.com "Title"`,
			expected: `<p><a href="https://example.com" title="Title">Example</a></p>`,
		},
		{
			name: "reference link with template URL",
			input: `[Example][1]

[1]: {{ .URL }} "Title"`,
			expected: `<p><a href="{{ .URL }}" title="Title">Example</a></p>`,
		},
		{
			name: "reference link with template URL containing spaces",
			input: `[Example][1]

[1]: {{ .Site.BaseURL }}/{{ .RelPermalink }} "{{ .Title }}"`,
			expected: `<p><a href="{{ .Site.BaseURL }}/{{ .RelPermalink }}" title="{{ .Title }}">Example</a></p>`,
		},
		{
			name: "shortcut reference link",
			input: `[Example]

[Example]: {{ .URL }} "{{ .Title }}"`,
			expected: `<p><a href="{{ .URL }}" title="{{ .Title }}">Example</a></p>`,
		},
		{
			name: "collapsed reference link",
			input: `[Example][]

[Example]: {{ .URL }} "{{ .Title }}"`,
			expected: `<p><a href="{{ .URL }}" title="{{ .Title }}">Example</a></p>`,
		},
		{
			name: "multiple reference links",
			input: `[First][1] and [Second][2]

[1]: {{ .FirstURL }} "First Title"
[2]: {{ .SecondURL }} "Second Title"`,
			expected: `<p><a href="{{ .FirstURL }}" title="First Title">First</a> and <a href="{{ .SecondURL }}" title="Second Title">Second</a></p>`,
		},
	}

	md := goldmark.New(
		goldmark.WithParser(NewParser()),
		goldmark.WithExtensions(NewExtension()),
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