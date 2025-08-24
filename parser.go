package goldmarktemplate

import (
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

// NewParser creates a parser with our custom link parser replacing the default
func NewParser() parser.Parser {
	// Get default parsers
	blockParsers := parser.DefaultBlockParsers()
	inlineParsers := parser.DefaultInlineParsers()
	paragraphTransformers := parser.DefaultParagraphTransformers()
	
	// Remove the default link parser and add ours
	filteredInlineParsers := make([]util.PrioritizedValue, 0, len(inlineParsers))
	for _, pv := range inlineParsers {
		// Check if this is goldmark's default link parser by type assertion
		// The default link parser is an unexported *linkParser type
		if _, isLinkParser := pv.Value.(interface{ Trigger() []byte }); isLinkParser {
			// Check if triggers match ['!', '[', ']'] 
			triggers := pv.Value.(interface{ Trigger() []byte }).Trigger()
			if len(triggers) == 3 && triggers[0] == '!' && triggers[1] == '[' && triggers[2] == ']' {
				// This is the default link parser, skip it
				continue
			}
		}
		filteredInlineParsers = append(filteredInlineParsers, pv)
	}
	
	// Add our templated link parser at the same priority (200)
	filteredInlineParsers = append(filteredInlineParsers, 
		util.Prioritized(NewLinkParser(), 200))
	
	return parser.NewParser(
		parser.WithBlockParsers(blockParsers...),
		parser.WithInlineParsers(filteredInlineParsers...),
		parser.WithParagraphTransformers(paragraphTransformers...),
	)
}