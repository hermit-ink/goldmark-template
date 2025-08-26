package goldmarktemplate

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/mermaid"
)

func TestMermaidIntegration(t *testing.T) {
	// Skip if mmdc (mermaid CLI) is not available (e.g., in CI environments)
	if _, err := exec.LookPath("mmdc"); err != nil {
		t.Skip("mmdc (mermaid CLI) not available, skipping server-side mermaid integration test")
	}
	tests := []struct {
		name                 string
		input                string
		expectedTemplateText string // Template action that should appear in SVG text content
	}{
		{
			name: "template action in mermaid node label",
			input: `~~~mermaid
graph TD
    A["{{ .StartNode }}"] --> B[End]
~~~`,
			expectedTemplateText: "<p>{{ .StartNode }}</p>",
		},
		{
			name: "template action in sequence diagram",
			input: `~~~mermaid
sequenceDiagram
    participant A as {{ .ActorA }}
    A->>B: {{ .Message }}
~~~`,
			expectedTemplateText: "{{ .ActorA }}",
		},
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			New(),
			&mermaid.Extender{
				RenderMode: mermaid.RenderModeServer,
			},
		),
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

			output := buf.String()

			// Verify SVG is generated (server-side rendering working)
			if !strings.Contains(output, "<svg") {
				t.Errorf("Expected SVG output from server-side mermaid rendering, but none found")
			}

			// Verify template action appears exactly as expected in the SVG
			if !strings.Contains(output, tt.expectedTemplateText) {
				t.Errorf("Expected template text %q not found in SVG output\nOutput: %s", tt.expectedTemplateText, output)
			}

			// Verify the output is wrapped in the expected mermaid div structure
			if !strings.Contains(output, `<div class="mermaid">`) {
				t.Errorf("Expected mermaid div wrapper not found in output")
			}
		})
	}
}