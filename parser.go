package goldmarktemplate

import (
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

// NewParser creates a parser with our custom link parser replacing the default
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

	return parser.NewParser(
		parser.WithBlockParsers(parser.DefaultBlockParsers()...),
		parser.WithInlineParsers(parsers...),
		parser.WithParagraphTransformers(parser.DefaultParagraphTransformers()...),
	)
}
