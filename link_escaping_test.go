package goldmarktemplate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

func TestLinkEscaping(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		description string
	}{
		{
			name: "link with HTML entities in URL",
			input: `[text](https://example.com?param=a&b<c"test)`,
			description: "URLs should be properly escaped",
		},
		{
			name: "link with HTML entities in text",
			input: `[text with & < > " chars](https://example.com)`,
			description: "Link text should be properly escaped",
		},
		{
			name: "link with template in URL",
			input: `[text]({{ .URL }}&param=value)`,
			description: "Templates in URLs should be preserved",
		},
		{
			name: "link with template in text",
			input: `[{{ .LinkText }}](https://example.com)`,
			description: "Templates in link text should be preserved",
		},
		{
			name: "link with templates and entities",
			input: `[{{ .Text }} & more]({{ .URL }}&param=test)`,
			description: "Mixed templates and entities in links",
		},
		{
			name: "reference link with entities",
			input: "[text with & < >][ref]\n\n[ref]: https://example.com?param=a&b",
			description: "Reference links should handle entities properly",
		},
		{
			name: "reference link with templates",
			input: "[{{ .Text }}][ref]\n\n[ref]: {{ .URL }}",
			description: "Reference links with templates",
		},
		{
			name: "link title with entities",
			input: `[text](https://example.com "Title with & < >")`,
			description: "Link titles should be escaped",
		},
		{
			name: "link title with templates",
			input: `[text](https://example.com "{{ .Title }}")`,
			description: "Link titles with templates should be preserved",
		},
	}

	// Create both goldmark standard and our template-enabled parsers  
	standardMd := goldmark.New()
	templateMd := goldmark.New(
		goldmark.WithExtensions(New()),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get goldmark's standard output
			var goldmarkBuf bytes.Buffer
			err := standardMd.Convert([]byte(tt.input), &goldmarkBuf)
			if err != nil {
				t.Fatalf("Failed to convert with goldmark: %v", err)
			}
			goldmarkOutput := strings.TrimSpace(goldmarkBuf.String())

			// Get our template-enabled output
			var templateBuf bytes.Buffer
			err = templateMd.Convert([]byte(tt.input), &templateBuf)
			if err != nil {
				t.Fatalf("Failed to convert with template extension: %v", err)
			}
			templateOutput := strings.TrimSpace(templateBuf.String())

			// For non-template content, our output should match goldmark exactly
			if !strings.Contains(tt.input, "{{") {
				if templateOutput != goldmarkOutput {
					t.Errorf("%s\nNon-template content should match goldmark exactly\nInput:     %q\nGoldmark:  %q\nTemplate:  %q", 
						tt.description, tt.input, goldmarkOutput, templateOutput)
				}
			} else {
				// For template content, log the difference for analysis
				t.Logf("Link comparison (%s):", tt.description)
				t.Logf("  Input:     %q", tt.input)
				t.Logf("  Goldmark:  %q", goldmarkOutput)
				t.Logf("  Template:  %q", templateOutput)
				
				// Basic validation - templates should be preserved
				if !strings.Contains(templateOutput, "{{") {
					t.Errorf("Template actions not preserved in link\nInput: %q\nOutput: %q", tt.input, templateOutput)
				}
			}
		})
	}
}

func TestAutoLinkEscaping(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		description string
	}{
		{
			name: "autolink with query parameters",
			input: `<https://example.com?a=b&c=d>`,
			description: "Autolinks with query params should be escaped properly",
		},
		{
			name: "email autolink",
			input: `<test@example.com>`,
			description: "Email autolinks should work correctly",
		},
		{
			name: "autolink with template",
			input: `<{{ .URL }}>`,
			description: "Autolinks with templates should preserve templates",
		},
		{
			name: "mixed autolinks",
			input: `Visit <https://example.com?a=b&c=d> or <{{ .URL }}>`,
			description: "Mixed autolinks with and without templates",
		},
	}

	// Create both goldmark standard and our template-enabled parsers
	standardMd := goldmark.New()
	templateMd := goldmark.New(goldmark.WithExtensions(New()))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get goldmark's standard output
			var goldmarkBuf bytes.Buffer
			err := standardMd.Convert([]byte(tt.input), &goldmarkBuf)
			if err != nil {
				t.Fatalf("Failed to convert with goldmark: %v", err)
			}
			goldmarkOutput := strings.TrimSpace(goldmarkBuf.String())

			// Get our template-enabled output
			var templateBuf bytes.Buffer
			err = templateMd.Convert([]byte(tt.input), &templateBuf)
			if err != nil {
				t.Fatalf("Failed to convert with template extension: %v", err)
			}
			templateOutput := strings.TrimSpace(templateBuf.String())

			// For non-template content, our output should match goldmark exactly
			if !strings.Contains(tt.input, "{{") {
				if templateOutput != goldmarkOutput {
					t.Errorf("%s\nNon-template content should match goldmark exactly\nInput:     %q\nGoldmark:  %q\nTemplate:  %q", 
						tt.description, tt.input, goldmarkOutput, templateOutput)
				}
			} else {
				// For template content, log the difference for analysis
				t.Logf("Autolink comparison (%s):", tt.description)
				t.Logf("  Input:     %q", tt.input)
				t.Logf("  Goldmark:  %q", goldmarkOutput)
				t.Logf("  Template:  %q", templateOutput)
				
				// Basic validation - templates should be preserved
				if !strings.Contains(templateOutput, "{{") {
					t.Errorf("Template actions not preserved in autolink\nInput: %q\nOutput: %q", tt.input, templateOutput)
				}
			}
		})
	}
}