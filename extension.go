package goldmarktemplate

import (
	"github.com/hermit-ink/goldmark-template/parser"
	"github.com/hermit-ink/goldmark-template/renderer/html"
	"github.com/yuin/goldmark"
	gparser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Extension is a goldmark extension for handling Go template actions
type Extension struct {
	parserOptions []gparser.Option
}

// New creates a new goldmark.Extender for template support
func New() goldmark.Extender {
	return &Extension{}
}

// WithParserOptions creates a new goldmark.Extender for template support with parser options
func WithParserOptions(opts ...gparser.Option) goldmark.Extender {
	return &Extension{parserOptions: opts}
}

// Extend configures the markdown processor to use our custom template action
// handling
func (e *Extension) Extend(m goldmark.Markdown) {
	// Create our new parser
	newParser := parser.ActionAwareParsers()
	
	// Apply user-provided parser options
	if len(e.parserOptions) > 0 {
		newParser.AddOptions(e.parserOptions...)
	}
	
	m.SetParser(newParser)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(html.NewRenderer(), 100),
			util.Prioritized(html.NewTemplateActionHTMLRenderer(), 500),
		),
	)
}
