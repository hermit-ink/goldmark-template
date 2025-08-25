package goldmarktemplate

import (
	"github.com/hermit-ink/goldmark-template/renderer/html"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Extension is a goldmark extension for handling Go template actions
type Extension struct{}

// NewExtension creates a new goldmark.Extender for template support
func NewExtension() goldmark.Extender {
	return &Extension{}
}

// Extend configures the markdown processor to use our custom template action
// handling
func (e *Extension) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(html.NewRenderer(), 100),
			util.Prioritized(html.NewTemplateActionHTMLRenderer(), 500),
		),
	)
}
