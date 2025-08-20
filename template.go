package goldmarktemplate

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// TemplateWriter is a custom HTML writer that preserves Go template directives
// without HTML escaping them, while properly handling escaped template cases.
type TemplateWriter struct {
	fallback html.Writer
}

// NewTemplateWriter creates a new TemplateWriter
func NewTemplateWriter(opts ...html.WriterOption) html.Writer {
	return &TemplateWriter{
		fallback: html.NewWriter(opts...),
	}
}

// Write writes content with normal processing
func (w *TemplateWriter) Write(writer util.BufWriter, source []byte) {
	w.fallback.Write(writer, source)
}

// SecureWrite writes content with security filtering
func (w *TemplateWriter) SecureWrite(writer util.BufWriter, source []byte) {
	w.fallback.SecureWrite(writer, source)
}

// RawWrite writes content while preserving Go template directives
func (w *TemplateWriter) RawWrite(writer util.BufWriter, source []byte) {
	n := 0
	i := 0

	for i < len(source) {
		// Check for template directive start
		if i < len(source)-1 && source[i] == '{' && source[i+1] == '{' {
			// Write everything before the directive (with escaping)
			if err := w.writeEscaped(writer, source[n:i]); err != nil {
				// If write fails, attempt to write remaining content and return
				_ = w.writeEscaped(writer, source[i:])
				return
			}

			// Check for escaped template syntax
			if w.isEscapedTemplate(source, i) {
				// Handle escaped template: {{"{{"}} or {{"}}"}}
				if w.matchesAt(source, i, []byte(`{{"{{"}}`)) {
					// Output literal {{
					if _, err := writer.WriteString("{{"); err != nil {
						return
					}
					i += 8
					n = i
					continue
				} else if w.matchesAt(source, i, []byte(`{{"}}"}}`)) {
					// Output literal }}
					if _, err := writer.WriteString("}}"); err != nil {
						return
					}
					i += 8
					n = i
					continue
				}
			}

			// Regular template directive - find the end
			end := w.findDirectiveEnd(source, i+2)
			if end > 0 {
				// Write the directive without escaping
				if _, err := writer.Write(source[i:end]); err != nil {
					return
				}
				i = end
				n = end
				continue
			}
		}
		i++
	}

	// Write remaining content with escaping
	_ = w.writeEscaped(writer, source[n:])
}

// isEscapedTemplate checks if the template at position is an escaped template
func (w *TemplateWriter) isEscapedTemplate(source []byte, pos int) bool {
	return w.matchesAt(source, pos, []byte(`{{"{{"}}`)) ||
		w.matchesAt(source, pos, []byte(`{{"}}"}}`))
}

// matchesAt checks if source matches pattern at position
func (w *TemplateWriter) matchesAt(source []byte, pos int, pattern []byte) bool {
	if pos+len(pattern) > len(source) {
		return false
	}
	return bytes.Equal(source[pos:pos+len(pattern)], pattern)
}

// findDirectiveEnd finds the end of a template directive, handling nested templates
// and string literals
func (w *TemplateWriter) findDirectiveEnd(source []byte, start int) int {
	depth := 1
	i := start
	inString := false
	inRawString := false

	for i < len(source)-1 {
		// Handle raw strings (backticks)
		if source[i] == '`' && !inString {
			inRawString = !inRawString
			i++
			continue
		}

		// Handle regular strings
		if source[i] == '"' && !inRawString && (i == start || source[i-1] != '\\') {
			inString = !inString
			i++
			continue
		}

		// Only process template markers outside strings
		if !inString && !inRawString {
			if source[i] == '{' && i+1 < len(source) && source[i+1] == '{' {
				depth++
				i += 2
			} else if source[i] == '}' && i+1 < len(source) && source[i+1] == '}' {
				depth--
				if depth == 0 {
					return i + 2
				}
				i += 2
			} else {
				i++
			}
		} else {
			i++
		}
	}

	return -1 // No matching closing found
}

// writeEscaped writes content with HTML escaping
func (w *TemplateWriter) writeEscaped(writer util.BufWriter, source []byte) error {
	for _, b := range source {
		if escaped := util.EscapeHTMLByte(b); escaped != nil {
			if _, err := writer.Write(escaped); err != nil {
				return err
			}
		} else {
			if err := writer.WriteByte(b); err != nil {
				return err
			}
		}
	}
	return nil
}

// TemplateExtension is a goldmark extension for handling Go templates
type TemplateExtension struct{}

// Extend configures the markdown processor to use our custom template handling
func (e *TemplateExtension) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(NewTemplateRenderer(), 100),
		),
	)
}

// TemplateRenderer is a custom renderer that uses TemplateWriter
type TemplateRenderer struct {
	html.Config
}

// NewTemplateRenderer creates a new TemplateRenderer
func NewTemplateRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &TemplateRenderer{
		Config: html.NewConfig(),
	}
	r.Writer = NewTemplateWriter()
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs registers rendering functions for code blocks and spans
func (r *TemplateRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindCodeBlock, r.renderCodeBlock)
	reg.Register(ast.KindFencedCodeBlock, r.renderFencedCodeBlock)
	reg.Register(ast.KindCodeSpan, r.renderCodeSpan)
}

func (r *TemplateRenderer) renderCodeBlock(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
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

func (r *TemplateRenderer) renderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
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

func (r *TemplateRenderer) renderCodeSpan(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
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

func (r *TemplateRenderer) writeLines(w util.BufWriter, source []byte, n ast.Node) error {
	l := n.Lines().Len()
	for i := 0; i < l; i++ {
		line := n.Lines().At(i)
		r.Writer.RawWrite(w, line.Value(source))
	}
	return nil
}

// New creates a new goldmark.Extender for template support
func New() goldmark.Extender {
	return &TemplateExtension{}
}
