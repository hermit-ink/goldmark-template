package html

import (
	"bytes"

	"github.com/hermit-ink/goldmark-template/ast"
	tutil "github.com/hermit-ink/goldmark-template/util"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	ghtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// Renderer is a custom renderer that uses Writer
type Renderer struct {
	ghtml.Config
}

// NewRenderer creates a new Renderer
func NewRenderer(opts ...ghtml.Option) renderer.NodeRenderer {
	r := &Renderer{
		Config: ghtml.NewConfig(),
	}
	r.Writer = NewWriter()
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs registers rendering functions for code blocks and spans
func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(gast.KindCodeBlock, r.renderCodeBlock)
	reg.Register(gast.KindFencedCodeBlock, r.renderFencedCodeBlock)
	reg.Register(gast.KindCodeSpan, r.renderCodeSpan)
	reg.Register(gast.KindLink, r.renderLink)
	reg.Register(gast.KindImage, r.renderImage)
	reg.Register(gast.KindAutoLink, r.renderAutoLink)
}

func (r *Renderer) renderCodeBlock(w util.BufWriter, source []byte, n gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		if _, err := w.WriteString("<pre><code>"); err != nil {
			return gast.WalkStop, err
		}
		if err := r.writeLines(w, source, n); err != nil {
			return gast.WalkStop, err
		}
	} else {
		if _, err := w.WriteString("</code></pre>\n"); err != nil {
			return gast.WalkStop, err
		}
	}
	return gast.WalkContinue, nil
}

func (r *Renderer) renderFencedCodeBlock(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	n := node.(*gast.FencedCodeBlock)
	if entering {
		if _, err := w.WriteString("<pre><code"); err != nil {
			return gast.WalkStop, err
		}
		language := n.Language(source)
		if language != nil {
			if _, err := w.WriteString(" class=\"language-"); err != nil {
				return gast.WalkStop, err
			}
			r.Writer.Write(w, language)
			if _, err := w.WriteString("\""); err != nil {
				return gast.WalkStop, err
			}
		}
		if err := w.WriteByte('>'); err != nil {
			return gast.WalkStop, err
		}
		if err := r.writeLines(w, source, n); err != nil {
			return gast.WalkStop, err
		}
	} else {
		if _, err := w.WriteString("</code></pre>\n"); err != nil {
			return gast.WalkStop, err
		}
	}
	return gast.WalkContinue, nil
}

func (r *Renderer) renderCodeSpan(w util.BufWriter, source []byte, n gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		if n.Attributes() != nil {
			if _, err := w.WriteString("<code"); err != nil {
				return gast.WalkStop, err
			}
			ghtml.RenderAttributes(w, n, ghtml.CodeAttributeFilter)
			if err := w.WriteByte('>'); err != nil {
				return gast.WalkStop, err
			}
		} else {
			if _, err := w.WriteString("<code>"); err != nil {
				return gast.WalkStop, err
			}
		}
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			segment := c.(*gast.Text).Segment
			value := segment.Value(source)
			if bytes.HasSuffix(value, []byte("\n")) {
				r.Writer.RawWrite(w, value[:len(value)-1])
				r.Writer.RawWrite(w, []byte(" "))
			} else {
				r.Writer.RawWrite(w, value)
			}
		}
		return gast.WalkSkipChildren, nil
	}
	if _, err := w.WriteString("</code>"); err != nil {
		return gast.WalkStop, err
	}
	return gast.WalkContinue, nil
}

func (r *Renderer) renderImage(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	if !entering {
		return gast.WalkContinue, nil
	}

	n := node.(*gast.Image)
	if _, err := w.WriteString("<img"); err != nil {
		return gast.WalkStop, err
	}
	if err := r.writeAttribute(w, "src", n.Destination); err != nil {
		return gast.WalkStop, err
	}
	if err := r.writeAttribute(w, "alt", r.extractTextContent(n, source)); err != nil {
		return gast.WalkStop, err
	}
	if err := r.writeAttribute(w, "title", n.Title); err != nil {
		return gast.WalkStop, err
	}
	if r.XHTML {
		if _, err := w.WriteString(" />"); err != nil {
			return gast.WalkStop, err
		}
	} else {
		if _, err := w.WriteString(">"); err != nil {
			return gast.WalkStop, err
		}
	}
	return gast.WalkSkipChildren, nil
}

func (r *Renderer) renderLink(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	n := node.(*gast.Link)
	if entering {
		if _, err := w.WriteString("<a"); err != nil {
			return gast.WalkStop, err
		}
		if err := r.writeAttribute(w, "href", n.Destination); err != nil {
			return gast.WalkStop, err
		}
		if err := r.writeAttribute(w, "title", n.Title); err != nil {
			return gast.WalkStop, err
		}
		if err := w.WriteByte('>'); err != nil {
			return gast.WalkStop, err
		}
	} else {
		if _, err := w.WriteString("</a>"); err != nil {
			return gast.WalkStop, err
		}
	}
	return gast.WalkContinue, nil
}

func (r *Renderer) extractTextContent(n gast.Node, source []byte) []byte {
	var buf bytes.Buffer
	_ = gast.Walk(n, func(node gast.Node, entering bool) (gast.WalkStatus, error) {
		if entering {
			if text, ok := node.(*gast.Text); ok {
				buf.Write(text.Segment.Value(source))
			} else if td, ok := node.(*ast.TemplateAction); ok {
				buf.Write(td.Content)
			}
		}
		return gast.WalkContinue, nil
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

	// Determine if this is a URL attribute that needs URL escaping
	isURLAttribute := name == "href" || name == "src"
	
	if hasAction(value) {
		// For values with templates, we need to handle URL vs HTML escaping properly
		r.writeAttributeWithTemplates(w, value, isURLAttribute)
	} else {
		// For values without templates, use goldmark's standard processing
		if isURLAttribute {
			// URL attributes need URL escaping then HTML escaping (like goldmark does)
			urlEscaped := util.URLEscape(value, true)
			htmlEscaped := util.EscapeHTML(urlEscaped)
			if _, err := w.Write(htmlEscaped); err != nil {
				return err
			}
		} else {
			// Non-URL attributes just need HTML escaping
			r.Writer.RawWrite(w, value)
		}
	}

	if _, err := w.WriteString("\""); err != nil {
		return err
	}
	return nil
}

// writeAttributeWithTemplates handles attribute values containing template actions
func (r *Renderer) writeAttributeWithTemplates(w util.BufWriter, value []byte, isURLAttribute bool) error {
	actionPattern := []byte("{{")
	n := 0
	i := 0

	for i < len(value) {
		// Skip non-template characters
		if i >= len(value)-1 || !bytes.HasPrefix(value[i:], actionPattern) {
			i++
			continue
		}

		// Process everything before the action
		if n < i {
			beforeAction := value[n:i]
			if isURLAttribute {
				// URL attributes need URL escaping then HTML escaping (like goldmark does)
				urlEscaped := util.URLEscape(beforeAction, true)
				htmlEscaped := util.EscapeHTML(urlEscaped)
				if _, err := w.Write(htmlEscaped); err != nil {
					return err
				}
			} else {
				// Non-URL attributes just need HTML escaping
				r.Writer.RawWrite(w, beforeAction)
			}
		}

		// Find and write the complete template action verbatim
		end := tutil.FindActionEnd(value, i)
		if end <= 0 {
			i++
			continue
		}

		if _, err := w.Write(value[i:end]); err != nil {
			return err
		}
		n = end
		i = end
	}

	// Process remaining content
	if n < len(value) {
		remaining := value[n:]
		if isURLAttribute {
			// URL attributes need URL escaping then HTML escaping (like goldmark does)
			urlEscaped := util.URLEscape(remaining, true)
			htmlEscaped := util.EscapeHTML(urlEscaped)
			if _, err := w.Write(htmlEscaped); err != nil {
				return err
			}
		} else {
			// Non-URL attributes just need HTML escaping
			r.Writer.RawWrite(w, remaining)
		}
	}
	return nil
}

func (r *Renderer) renderAutoLink(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	if !entering {
		return gast.WalkContinue, nil
	}

	n := node.(*gast.AutoLink)
	url := n.URL(source)
	label := n.Label(source)

	if n.AutoLinkType == gast.AutoLinkEmail {
		if _, err := w.WriteString(`<a href="mailto:`); err != nil {
			return gast.WalkStop, err
		}
	} else {
		if _, err := w.WriteString(`<a href="`); err != nil {
			return gast.WalkStop, err
		}
	}

	// Use raw write to preserve templates in URLs
	if hasAction(url) {
		if _, err := w.Write(url); err != nil {
			return gast.WalkStop, err
		}
	} else {
		r.Writer.Write(w, url)
	}

	if _, err := w.WriteString(`">`); err != nil {
		return gast.WalkStop, err
	}
	r.Writer.RawWrite(w, label)
	if _, err := w.WriteString(`</a>`); err != nil {
		return gast.WalkStop, err
	}
	return gast.WalkSkipChildren, nil
}

func (r *Renderer) writeLines(w util.BufWriter, source []byte, n gast.Node) error {
	l := n.Lines().Len()
	for i := range l {
		line := n.Lines().At(i)
		r.Writer.RawWrite(w, line.Value(source))
	}
	return nil
}
