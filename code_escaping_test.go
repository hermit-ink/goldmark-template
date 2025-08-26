package goldmarktemplate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
)

func TestCodeBlockEscaping(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		description string
	}{
		{
			name: "code block with HTML entities",
			input: "```\nCode with & < > \" characters\n```",
			description: "Code blocks should escape HTML characters but preserve templates",
		},
		{
			name: "code block with template actions",
			input: "```\nCode with {{ .Value }} & more\n```",
			description: "Code blocks with templates should preserve templates but escape other chars",
		},
		{
			name: "code block with entities and templates mixed",
			input: "```\n&amp; {{ .Value }} &lt; more\n```",
			description: "Mixed content should handle entities and templates correctly",
		},
		{
			name: "indented code block with HTML chars",
			input: "    Code with & < > \" characters",
			description: "Indented code blocks should also escape HTML properly",
		},
		{
			name: "indented code block with templates",
			input: "    {{ .Value }} & more < chars",
			description: "Indented code blocks with templates",
		},
		{
			name: "fenced code with language and HTML",
			input: "```go\nif value < 10 && flag {\n    fmt.Println(\"test\")\n}\n```",
			description: "Fenced code blocks with language should escape HTML",
		},
		{
			name: "fenced code with language and templates",
			input: "```go\nif {{ .Value }} < 10 && flag {\n    fmt.Println(\"{{ .Message }}\")\n}\n```",
			description: "Fenced code blocks with templates should preserve them",
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
				t.Logf("Code block comparison (%s):", tt.description)
				t.Logf("  Input:     %q", tt.input)
				t.Logf("  Goldmark:  %q", goldmarkOutput)
				t.Logf("  Template:  %q", templateOutput)
				
				// Basic validation - templates should be preserved
				if !strings.Contains(templateOutput, "{{") {
					t.Errorf("Template actions not preserved in code block\nInput: %q\nOutput: %q", tt.input, templateOutput)
				}
			}
		})
	}
}

func TestCodeSpanEscaping(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		description string
	}{
		{
			name: "code span with HTML entities",
			input: "Inline `code with & < > \" chars` here",
			description: "Code spans should escape HTML characters",
		},
		{
			name: "code span with template actions",
			input: "Inline `{{ .Value }} & more` code",
			description: "Code spans with templates should preserve templates",
		},
		{
			name: "code span with entities and templates",
			input: "Code `&amp; {{ .Value }} &lt;` span",
			description: "Mixed entities and templates in code spans",
		},
		{
			name: "multiple code spans",
			input: "Code `{{ .A }}` and `{{ .B }}` spans",
			description: "Multiple code spans with templates",
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
				t.Logf("Code span comparison (%s):", tt.description)
				t.Logf("  Input:     %q", tt.input)
				t.Logf("  Goldmark:  %q", goldmarkOutput)
				t.Logf("  Template:  %q", templateOutput)
				
				// Basic validation - templates should be preserved
				if !strings.Contains(templateOutput, "{{") {
					t.Errorf("Template actions not preserved in code span\nInput: %q\nOutput: %q", tt.input, templateOutput)
				}
			}
		})
	}
}