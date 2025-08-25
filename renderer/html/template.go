package html

import (
	"github.com/hermit-ink/goldmark-template/ast"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	ghtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// TemplateActionHTMLRenderer renders TemplateAction nodes directly into the
// output with no HTML/URL escaping
type TemplateActionHTMLRenderer struct {
	ghtml.Config
}

// NewTemplateActionHTMLRenderer returns a new TemplateActionHTMLRenderer
func NewTemplateActionHTMLRenderer(opts ...ghtml.Option) renderer.NodeRenderer {
	r := &TemplateActionHTMLRenderer{
		Config: ghtml.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs
func (r *TemplateActionHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindTemplateAction, r.render)
}

// render renders template actions as raw content (no HTML encoding)
func (r *TemplateActionHTMLRenderer) render(
	w util.BufWriter, source []byte, n gast.Node, entering bool,
) (gast.WalkStatus, error) {
	if entering {
		if node, ok := n.(*ast.TemplateAction); ok {
			// Write the template action as-is (no HTML encoding)
			_, err := w.Write(node.Content)
			if err != nil {
				return gast.WalkStop, err
			}
		}
	}
	return gast.WalkContinue, nil
}
