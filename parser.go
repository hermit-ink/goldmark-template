package goldmarktemplate

import (
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

// NewParser creates a parser with our custom link parser and reference parser replacing the defaults
func NewParser() parser.Parser {
	ip := parser.DefaultInlineParsers()

	parsers := make([]util.PrioritizedValue, 0, len(ip))
	for _, pv := range ip {
		if lp, ok := pv.Value.(interface{ Trigger() []byte }); ok {
			t := lp.Trigger()
			if len(t) == 3 && t[0] == '!' && t[1] == '[' && t[2] == ']' {
				// Looks like a duck, talks like a duck
				parsers = append(parsers, util.Prioritized(NewLinkParser(), 200))
				continue
			}
		}
		parsers = append(parsers, pv)
	}

	pt := parser.DefaultParagraphTransformers()
	transformers := make([]util.PrioritizedValue, 0, len(pt))
	for _, pv := range pt {
		if pv.Value != parser.LinkReferenceParagraphTransformer {
			transformers = append(transformers, pv)
		}
	}
	// Add our custom reference paragraph transformer
	transformers = append(
		transformers,
		util.Prioritized(NewReferenceParagraphTransformer(), 999))

	return parser.NewParser(
		parser.WithBlockParsers(parser.DefaultBlockParsers()...),
		parser.WithInlineParsers(parsers...),
		parser.WithParagraphTransformers(transformers...),
	)
}
