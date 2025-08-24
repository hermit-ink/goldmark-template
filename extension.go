package goldmarktemplate

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Extension is a goldmark extension for handling Go templates
type Extension struct{}

// NewExtension creates a new goldmark.Extender for template support
func NewExtension() goldmark.Extender {
	return &Extension{}
}

// Extend configures the markdown processor to use our custom template handling
func (e *Extension) Extend(m goldmark.Markdown) {
	// Don't add parser here - it's handled by NewParser()
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(NewRenderer(), 100),
		),
	)
}