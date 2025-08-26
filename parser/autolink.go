package parser

import (
	"bytes"

	tutil "github.com/hermit-ink/goldmark-template/util"
	"github.com/yuin/goldmark/ast"
	gparser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type autoLinkParser struct{}

// NewAutoLinkParser returns a new InlineParser that parses autolinks with Go
// template action support
func NewAutoLinkParser() gparser.InlineParser {
	return &autoLinkParser{}
}

func (s *autoLinkParser) Trigger() []byte {
	return []byte{'<'}
}

func (s *autoLinkParser) Parse(parent ast.Node, block text.Reader, pc gparser.Context) ast.Node {
	line, segment := block.PeekLine()

	// First check if this contains an action - if so, treat as URL autolink
	content := line[1:] // Skip opening '<'
	closePos := bytes.IndexByte(content, '>')
	if closePos < 0 {
		return nil
	}

	urlContent := content[:closePos]

	// If it starts with an action then it *should* be an autolink.
	//
	// <{{ .URL }}>
	// <{{> will also get treated like an autolink even though its not valid
	// but that's ok
	if len(urlContent) >= 2 && urlContent[0] == '{' && urlContent[1] == '{' {
		stop := closePos + 1 // +1 for the '>'
		value := ast.NewTextSegment(text.NewSegment(segment.Start+1, segment.Start+stop))
		block.Advance(stop + 1)
		return ast.NewAutoLink(ast.AutoLinkURL, value)
	}

	// If it starts with a URL-like string (util.FindURLIndex) and it has a
	// template action in it then construct an autolink ast node and return it
	// <https://......{{.Something}}>
	if util.FindURLIndex(urlContent) > 0 && tutil.ContainsAction(urlContent) {
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
