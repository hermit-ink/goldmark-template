package goldmarktemplate

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// Simple test parser to understand trigger behavior
type debugParser struct{}

func (p *debugParser) Trigger() []byte {
	return []byte{'[', ']'}
}

func (p *debugParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	fmt.Printf("DEBUG-PARSER: Called with line[0]=%c, line=%q\n", line[0], string(line))
	return nil // Don't actually parse anything
}

func TestParserTriggers(t *testing.T) {
	// Test with default setup plus our debug parser
	md := goldmark.New(
		goldmark.WithExtensions(NewTemplatedHTMLExtension()),
		goldmark.WithParserOptions(
			parser.WithInlineParsers(
				util.Prioritized(&debugParser{}, 50),   // Low priority debug parser
			),
		),
		goldmark.WithRendererOptions(html.WithXHTML()),
	)

	input := "[Link]({{ .URL }})"
	fmt.Printf("\n=== Testing input: %q ===\n", input)
	
	var buf bytes.Buffer
	err := md.Convert([]byte(input), &buf)
	if err != nil {
		t.Fatalf("Failed to convert: %v", err)
	}
	
	fmt.Printf("Output: %q\n", buf.String())
}