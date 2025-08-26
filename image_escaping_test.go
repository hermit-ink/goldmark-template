package goldmarktemplate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

func TestImageEscaping(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		description string
	}{
		{
			name: "image with HTML entities in alt text",
			input: `![alt with & < > " chars](image.jpg)`,
			description: "Image alt text should be properly escaped",
		},
		{
			name: "image with HTML entities in URL",
			input: `![alt](image.jpg?param=a&b<c"test)`,
			description: "Image URLs should be properly escaped",
		},
		{
			name: "image with template in alt text",
			input: `![{{ .AltText }}](image.jpg)`,
			description: "Templates in alt text should be preserved",
		},
		{
			name: "image with template in URL",
			input: `![alt]({{ .ImageURL }})`,
			description: "Templates in image URLs should be preserved",
		},
		{
			name: "image with templates and entities",
			input: `![{{ .Alt }} & more]({{ .URL }}&param=test)`,
			description: "Mixed templates and entities in images",
		},
		{
			name: "image with title containing entities",
			input: `![alt](image.jpg "Title with & < >")`,
			description: "Image titles should be escaped",
		},
		{
			name: "image with template in title",
			input: `![alt](image.jpg "{{ .Title }}")`,
			description: "Image titles with templates should be preserved",
		},
		{
			name: "reference image with entities",
			input: "![alt with & < >][ref]\n\n[ref]: image.jpg?param=a&b",
			description: "Reference images should handle entities properly",
		},
		{
			name: "reference image with templates",
			input: "![{{ .Alt }}][ref]\n\n[ref]: {{ .URL }}",
			description: "Reference images with templates",
		},
		{
			name: "complex image with all features",
			input: `![{{ .Alt }} & text]({{ .URL }}&param=value "{{ .Title }} & more")`,
			description: "Images with templates and entities in all attributes",
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
				t.Logf("Image comparison (%s):", tt.description)
				t.Logf("  Input:     %q", tt.input)
				t.Logf("  Goldmark:  %q", goldmarkOutput)
				t.Logf("  Template:  %q", templateOutput)
				
				// Basic validation - templates should be preserved
				if !strings.Contains(templateOutput, "{{") {
					t.Errorf("Template actions not preserved in image\nInput: %q\nOutput: %q", tt.input, templateOutput)
				}
			}
		})
	}
}