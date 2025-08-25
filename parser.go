package goldmarktemplate

import (
	"github.com/hermit-ink/goldmark-template/parser"
	gparser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

// NewParser creates a parser with our custom link parser and reference parser replacing the defaults
func NewParser() gparser.Parser {
	ip := gparser.DefaultInlineParsers()

	parsers := make([]util.PrioritizedValue, 0, len(ip))
	for _, pv := range ip {
		if lp, ok := pv.Value.(interface{ Trigger() []byte }); ok {
			t := lp.Trigger()
			if len(t) == 3 && t[0] == '!' && t[1] == '[' && t[2] == ']' {
				// Looks like a duck, talks like a duck
				parsers = append(parsers, util.Prioritized(parser.NewLinkParser(), 200))
				continue
			}
			if len(t) == 1 && t[0] == '<' && pv.Priority == 300 {
				parsers = append(parsers, util.Prioritized(parser.NewAutoLinkParser(), 300))
				continue
			}
		}
		parsers = append(parsers, pv)
	}

	parsers = append(
		parsers,
		util.Prioritized(parser.NewTemplateActionParser(), 50))

	pt := gparser.DefaultParagraphTransformers()
	transformers := make([]util.PrioritizedValue, 0, len(pt))
	for _, pv := range pt {
		if pv.Value != gparser.LinkReferenceParagraphTransformer {
			transformers = append(transformers, pv)
		}
	}
	// Add our custom reference paragraph transformer
	transformers = append(
		transformers,
		util.Prioritized(parser.NewReferenceParagraphTransformer(), 999))

	return gparser.NewParser(
		gparser.WithBlockParsers(gparser.DefaultBlockParsers()...),
		gparser.WithInlineParsers(parsers...),
		gparser.WithParagraphTransformers(transformers...),
	)
}
