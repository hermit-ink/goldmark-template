package goldmarktemplate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

func TestBasicImages(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "image with action in src",
			input:    "![Alt]({{ .ImageURL }})",
			expected: "<p><img src=\"{{ .ImageURL }}\" alt=\"Alt\" /></p>",
		},
		{
			name:     "image with action in alt and src",
			input:    "![{{ .AltText }}]({{ .ImageURL }})",
			expected: "<p><img src=\"{{ .ImageURL }}\" alt=\"{{ .AltText }}\" /></p>",
		},
		{
			name:     "image with action containing quotes in alt",
			input:    "![{{ .Title | quote }}](image.jpg)",
			expected: "<p><img src=\"image.jpg\" alt=\"{{ .Title | quote }}\" /></p>",
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

func TestComplexImageAlt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "image with inline styles in alt text",
			input:    `![**Bold {{.Alt}}** text]({{.URL}})`,
			expected: `<p><img src="{{.URL}}" alt="Bold {{.Alt}} text" /></p>`,
		},
		{
			name:     "image alt with HTML chars and templates",
			input:    "![{{ .Title }} > {{ .Subtitle }}](image.jpg)",
			expected: "<p><img src=\"image.jpg\" alt=\"{{ .Title }} &gt; {{ .Subtitle }}\" /></p>",
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

func TestImageTitles(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "image with action in alt and title",
			input:    `![A picture of {{.Image}}](https://{{.URL}} "{{.Image}} in colour")`,
			expected: `<p><img src="https://{{.URL}}" alt="A picture of {{.Image}}" title="{{.Image}} in colour" /></p>`,
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
