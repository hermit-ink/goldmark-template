package goldmarktemplate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

func TestHTMLEscapingConsistency(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedGoldmark string // What goldmark would output
		description    string
	}{
		{
			name:  "basic HTML entities in text",
			input: `Text with & < > " characters`,
			expectedGoldmark: `<p>Text with &amp; &lt; &gt; &quot; characters</p>`,
			description: "Basic HTML escaping should match goldmark exactly",
		},
		{
			name:  "HTML entities with template action",
			input: `Text with & < > {{ .Value }} " characters`,
			expectedGoldmark: `<p>Text with &amp; &lt; &gt; {{ .Value }} &quot; characters</p>`,
			description: "HTML escaping around template actions",
		},
		{
			name:  "numeric character references",
			input: `Text with &#123; &#x41; entities`,
			expectedGoldmark: `<p>Text with { A entities</p>`,
			description: "Goldmark resolves numeric entities, we might not",
		},
		{
			name:  "named HTML entities",
			input: `Text with &amp; &lt; &gt; &quot; entities`,
			expectedGoldmark: `<p>Text with &amp; &lt; &gt; " entities</p>`,
			description: "Goldmark resolves named entities, we might not",
		},
		{
			name:  "backslash escaping",
			input: `Text with \* \& \< escaped chars`,
			expectedGoldmark: `<p>Text with * &amp; &lt; escaped chars</p>`,
			description: "Goldmark unescapes backslash-punctuation, we might not",
		},
		{
			name:  "mixed entities and templates",
			input: `&amp; {{ .Value }} &lt; more text &gt;`,
			expectedGoldmark: `<p>&amp; {{ .Value }} &lt; more text &gt;</p>`,
			description: "Entities around templates should be handled consistently",
		},
		{
			name:  "template in HTML attribute with entities",
			input: `<div title="Value: {{ .Val }} &amp; more">Content</div>`,
			expectedGoldmark: `<div title="Value: {{ .Val }} &amp; more">Content</div>`,
			description: "Entities in attributes with templates",
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

			// Compare expected vs actual goldmark (sanity check)
			if goldmarkOutput != tt.expectedGoldmark {
				t.Logf("WARNING: Expected goldmark output doesn't match actual:")
				t.Logf("  Input:    %q", tt.input)
				t.Logf("  Expected: %q", tt.expectedGoldmark)
				t.Logf("  Actual:   %q", goldmarkOutput)
			}

			// For non-template content, our output should match goldmark exactly
			if !strings.Contains(tt.input, "{{") {
				if templateOutput != goldmarkOutput {
					t.Errorf("%s\nNon-template content should match goldmark exactly\nInput:     %q\nGoldmark:  %q\nTemplate:  %q", 
						tt.description, tt.input, goldmarkOutput, templateOutput)
				}
			} else {
				// For template content, log the difference for analysis
				t.Logf("Template content comparison (%s):", tt.description)
				t.Logf("  Input:     %q", tt.input)
				t.Logf("  Goldmark:  %q", goldmarkOutput)
				t.Logf("  Template:  %q", templateOutput)
			}
		})
	}
}