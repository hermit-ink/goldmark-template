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

	paragraphTransformers := []util.PrioritizedValue{
		util.Prioritized(LinkReferenceParagraphTransformer, 100),
	}

	return gparser.NewParser(
		gparser.WithBlockParsers(gparser.DefaultBlockParsers()...),
		gparser.WithInlineParsers(inlineParsers...),
		gparser.WithParagraphTransformers(paragraphTransformers...),
	)
}
