package goldmarktemplate

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type autoLinkParser struct {
}

// NewAutoLinkParser returns a new InlineParser that parses autolinks with template support
func NewAutoLinkParser() parser.InlineParser {
	return &autoLinkParser{}
}

func (s *autoLinkParser) Trigger() []byte {
	return []byte{'<'}
}

func (s *autoLinkParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, segment := block.PeekLine()
	
	// First check if this contains a template - if so, treat as URL autolink
	content := line[1:] // Skip opening '<'
	closePos := bytes.IndexByte(content, '>')
	if closePos < 0 {
		return nil
	}
	
	urlContent := content[:closePos]
	
	// Only treat as autolink if content STARTS with template
	if len(urlContent) >= 2 && urlContent[0] == '{' && urlContent[1] == '{' {
		stop := closePos + 1 // +1 for the '>' 
		value := ast.NewTextSegment(text.NewSegment(segment.Start+1, segment.Start+stop))
		block.Advance(stop + 1)
		return ast.NewAutoLink(ast.AutoLinkURL, value)
	}
	
	// Otherwise, use goldmark's original logic
	stop := util.FindEmailIndex(line[1:])
	typ := ast.AutoLinkType(ast.AutoLinkEmail)
	if stop < 0 {
		stop = util.FindURLIndex(line[1:])
		typ = ast.AutoLinkURL
	}
	if stop < 0 {
		return nil
	}
	stop++
	if stop >= len(line) || line[stop] != '>' {
		return nil
	}
	value := ast.NewTextSegment(text.NewSegment(segment.Start+1, segment.Start+stop))
	block.Advance(stop + 1)
	return ast.NewAutoLink(typ, value)
}