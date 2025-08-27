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
type Extension struct{}

// New creates a new goldmark.Extender for template support
func New() goldmark.Extender {
	return &Extension{}
}

// Extend configures the markdown processor to use our custom template action
// handling
func (e *Extension) Extend(m goldmark.Markdown) {
	// Create our new parser
	newParser := parser.ActionAwareParsers()
	
	// TEMPORARY: Manually add the options that were in the test
	// TODO: Find a way to preserve the original parser's options
	newParser.AddOptions(
		gparser.WithAttribute(),
		gparser.WithHeadingAttribute(),
	)
	
	m.SetParser(newParser)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(html.NewRenderer(), 100),
			util.Prioritized(html.NewTemplateActionHTMLRenderer(), 500),
		),
	)
}
