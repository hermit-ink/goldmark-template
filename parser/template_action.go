package parser

import (
	"github.com/hermit-ink/goldmark-template/ast"
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

	i := 2
	// TODO: backticks
	inDoubleQuotes := false
	inSingleQuotes := false

	for i < len(line)-1 {
		char := line[i]

		if char == '"' && !inSingleQuotes {
			inDoubleQuotes = !inDoubleQuotes
		} else if char == '\'' && !inDoubleQuotes {
			inSingleQuotes = !inSingleQuotes
		}

		if !inDoubleQuotes && !inSingleQuotes && char == '}' && line[i+1] == '}' {
			content := line[0 : i+2]
			nodeSegment := segment.WithStop(segment.Start + i + 2)
			node := ast.NewTemplateAction(content, nodeSegment)
			block.Advance(i + 2)
			return node
		}
		i++
	}

	return nil
}

func (s *templateActionParser) CloseBlock(parent gast.Node, pc gparser.Context) {
	// nothing to do
}
