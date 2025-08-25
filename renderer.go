package goldmarktemplate

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// Renderer is a custom renderer that uses Writer
type Renderer struct {
	html.Config
}

// NewRenderer creates a new Renderer
func NewRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &Renderer{
		Config: html.NewConfig(),
	}
	r.Writer = NewWriter()
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs registers rendering functions for code blocks and spans
func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindCodeBlock, r.renderCodeBlock)
	reg.Register(ast.KindFencedCodeBlock, r.renderFencedCodeBlock)
	reg.Register(ast.KindCodeSpan, r.renderCodeSpan)
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindImage, r.renderImage)
	reg.Register(ast.KindAutoLink, r.renderAutoLink)
}

func (r *Renderer) renderCodeBlock(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		if _, err := w.WriteString("<pre><code>"); err != nil {
			return ast.WalkStop, err
		}
		if err := r.writeLines(w, source, n); err != nil {
			return ast.WalkStop, err
		}
	} else {
		if _, err := w.WriteString("</code></pre>\n"); err != nil {
			return ast.WalkStop, err
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.FencedCodeBlock)
	if entering {
		if _, err := w.WriteString("<pre><code"); err != nil {
			return ast.WalkStop, err
		}
		language := n.Language(source)
		if language != nil {
			if _, err := w.WriteString(" class=\"language-"); err != nil {
				return ast.WalkStop, err
			}
			r.Writer.Write(w, language)
			if _, err := w.WriteString("\""); err != nil {
				return ast.WalkStop, err
			}
		}
		if err := w.WriteByte('>'); err != nil {
			return ast.WalkStop, err
		}
		if err := r.writeLines(w, source, n); err != nil {
			return ast.WalkStop, err
		}
	} else {
		if _, err := w.WriteString("</code></pre>\n"); err != nil {
			return ast.WalkStop, err
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderCodeSpan(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		if n.Attributes() != nil {
			if _, err := w.WriteString("<code"); err != nil {
				return ast.WalkStop, err
			}
			html.RenderAttributes(w, n, html.CodeAttributeFilter)
			if err := w.WriteByte('>'); err != nil {
				return ast.WalkStop, err
			}
		} else {
			if _, err := w.WriteString("<code>"); err != nil {
				return ast.WalkStop, err
			}
		}
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			segment := c.(*ast.Text).Segment
			value := segment.Value(source)
			if bytes.HasSuffix(value, []byte("\n")) {
				r.Writer.RawWrite(w, value[:len(value)-1])
				r.Writer.RawWrite(w, []byte(" "))
			} else {
				r.Writer.RawWrite(w, value)
			}
		}
		return ast.WalkSkipChildren, nil
	}
	if _, err := w.WriteString("</code>"); err != nil {
		return ast.WalkStop, err
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*ast.Image)
	if _, err := w.WriteString("<img"); err != nil {
		return ast.WalkStop, err
	}
	if err := r.writeAttribute(w, "src", n.Destination); err != nil {
		return ast.WalkStop, err
	}
	if err := r.writeAttribute(w, "alt", r.extractTextContent(n, source)); err != nil {
		return ast.WalkStop, err
	}
	if err := r.writeAttribute(w, "title", n.Title); err != nil {
		return ast.WalkStop, err
	}
	if r.XHTML {
		if _, err := w.WriteString(" />"); err != nil {
			return ast.WalkStop, err
		}
	} else {
		if _, err := w.WriteString(">"); err != nil {
			return ast.WalkStop, err
		}
	}
	return ast.WalkSkipChildren, nil
}

func (r *Renderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Link)
	if entering {
		if _, err := w.WriteString("<a"); err != nil {
			return ast.WalkStop, err
		}
		if err := r.writeAttribute(w, "href", n.Destination); err != nil {
			return ast.WalkStop, err
		}
		if err := r.writeAttribute(w, "title", n.Title); err != nil {
			return ast.WalkStop, err
		}
		if err := w.WriteByte('>'); err != nil {
			return ast.WalkStop, err
		}
	} else {
		if _, err := w.WriteString("</a>"); err != nil {
			return ast.WalkStop, err
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) extractTextContent(n ast.Node, source []byte) []byte {
	var buf bytes.Buffer
	ast.Walk(n, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if text, ok := node.(*ast.Text); ok {
				buf.Write(text.Segment.Value(source))
			}
		}
		return ast.WalkContinue, nil
	})
	return buf.Bytes()
}

// writeAttribute writes an HTML attribute with template preservation
func (r *Renderer) writeAttribute(w util.BufWriter, name string, value []byte) error {
	if value == nil {
		return nil
	}
	if _, err := w.WriteString(" " + name + "=\""); err != nil {
		return err
	}

	// For attribute values, if they contain templates, preserve the entire content
	// to avoid escaping characters between template directives
	if hasTemplate(value) {
		if _, err := w.Write(value); err != nil {
			return err
		}
	} else {
		r.Writer.RawWrite(w, value)
	}

	if _, err := w.WriteString("\""); err != nil {
		return err
	}
	return nil
}

func (r *Renderer) renderAutoLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*ast.AutoLink)
	url := n.URL(source)
	label := n.Label(source)

	if n.AutoLinkType == ast.AutoLinkEmail {
		if _, err := w.WriteString(`<a href="mailto:`); err != nil {
			return ast.WalkStop, err
		}
	} else {
		if _, err := w.WriteString(`<a href="`); err != nil {
			return ast.WalkStop, err
		}
	}

	// Use raw write to preserve templates in URLs
	if hasTemplate(url) {
		if _, err := w.Write(url); err != nil {
			return ast.WalkStop, err
		}
	} else {
		r.Writer.Write(w, url)
	}

	if _, err := w.WriteString(`">`); err != nil {
		return ast.WalkStop, err
	}
	r.Writer.RawWrite(w, label)
	if _, err := w.WriteString(`</a>`); err != nil {
		return ast.WalkStop, err
	}
	return ast.WalkSkipChildren, nil
}

func (r *Renderer) writeLines(w util.BufWriter, source []byte, n ast.Node) error {
	l := n.Lines().Len()
	for i := range l {
		line := n.Lines().At(i)
		r.Writer.RawWrite(w, line.Value(source))
	}
	return nil
}