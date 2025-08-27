package parser

import (
	gparser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

func ActionAwareParsers() gparser.Parser {
	inlineParsers := []util.PrioritizedValue{
		util.Prioritized(NewCodeSpanParser(), 100),
		util.Prioritized(NewLinkParser(), 200),
		util.Prioritized(NewAutoLinkParser(), 300),
		util.Prioritized(gparser.NewRawHTMLParser(), 400),
		util.Prioritized(gparser.NewEmphasisParser(), 500),
		util.Prioritized(NewTemplateActionParser(), 600),
	}

	blockParsers := []util.PrioritizedValue{
		util.Prioritized(gparser.NewSetextHeadingParser(), 100),
		util.Prioritized(gparser.NewThematicBreakParser(), 200),
		util.Prioritized(gparser.NewListParser(), 300),
		util.Prioritized(gparser.NewListItemParser(), 400),
		util.Prioritized(gparser.NewCodeBlockParser(), 500),
		util.Prioritized(NewATXHeadingParser(), 600),
		util.Prioritized(gparser.NewFencedCodeBlockParser(), 700),
		util.Prioritized(gparser.NewBlockquoteParser(), 800),
		util.Prioritized(gparser.NewHTMLBlockParser(), 900),
		util.Prioritized(gparser.NewParagraphParser(), 1000),
	}

	paragraphTransformers := []util.PrioritizedValue{
		util.Prioritized(LinkReferenceParagraphTransformer, 100),
	}

	return gparser.NewParser(
		gparser.WithBlockParsers(blockParsers...),
		gparser.WithInlineParsers(inlineParsers...),
		gparser.WithParagraphTransformers(paragraphTransformers...),
	)
}
