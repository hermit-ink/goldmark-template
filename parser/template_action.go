package parser

import (
	"github.com/hermit-ink/goldmark-template/ast"
	tutil "github.com/hermit-ink/goldmark-template/util"
	gast "github.com/yuin/goldmark/ast"
	gparser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// templateActionParser is an inline parser for Go template actions
type templateActionParser struct{}

// NewTemplateActionParser returns a new InlineParser that parses go template
// actions
func NewTemplateActionParser() gparser.InlineParser {
	return &templateActionParser{}
}

// Trigger returns characters that trigger this parser
func (s *templateActionParser) Trigger() []byte {
	return []byte{'{'}
}

func (s *templateActionParser) Parse(parent gast.Node, block text.Reader, pc gparser.Context) gast.Node {
	line, segment := block.PeekLine()

	if len(line) < 2 || line[0] != '{' || line[1] != '{' {
		return nil
	}


	endPos := tutil.FindActionEnd(line, 0)
	if endPos == -1 {
		return nil
	}

	content := line[0:endPos]
	nodeSegment := segment.WithStop(segment.Start + endPos)
	node := ast.NewTemplateAction(content, nodeSegment)
	block.Advance(endPos)
	return node
}

func (s *templateActionParser) CloseBlock(parent gast.Node, pc gparser.Context) {
	// nothing to do
}
