package goldmarktemplate

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// Template detection utilities
var (
	templatePattern = []byte("{{")
)

// hasTemplate checks if content contains template directives
func hasTemplate(content []byte) bool {
	return bytes.Contains(content, templatePattern)
}

// Writer is a custom HTML writer that preserves Go template directives
// without HTML escaping them, while properly handling escaped template cases.
type Writer struct {
	fallback html.Writer
}

// NewWriter creates a new Writer
func NewWriter(opts ...html.WriterOption) html.Writer {
	return &Writer{
		fallback: html.NewWriter(opts...),
	}
}

// Write writes content with normal processing
func (w *Writer) Write(writer util.BufWriter, source []byte) {
	w.fallback.Write(writer, source)
}

// SecureWrite writes content with security filtering
func (w *Writer) SecureWrite(writer util.BufWriter, source []byte) {
	w.fallback.SecureWrite(writer, source)
}

// RawWrite writes content while preserving Go template directives
func (w *Writer) RawWrite(writer util.BufWriter, source []byte) {
	n := 0
	i := 0

	for i < len(source) {
		// Check for template directive start
		if i < len(source)-1 && bytes.HasPrefix(source[i:], templatePattern) {
			// Write everything before the directive (with escaping)
			if err := w.writeEscaped(writer, source[n:i]); err != nil {
				// If write fails, attempt to write remaining content and return
				_ = w.writeEscaped(writer, source[i:])
				return
			}

			// Find the end of template directive
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

// findDirectiveEnd finds the end of a template directive, handling nested templates
// and string literals
func (w *Writer) findDirectiveEnd(source []byte, start int) int {
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
func (w *Writer) writeEscaped(writer util.BufWriter, source []byte) error {
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


// Extension is a goldmark extension for handling Go templates
type Extension struct{}

// Extend configures the markdown processor to use our custom template handling
func (e *Extension) Extend(m goldmark.Markdown) {
	// Don't add parser here - it's handled by NewTemplatedParser()
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(NewRenderer(), 100),
		),
	)
}

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

func (r *Renderer) writeLines(w util.BufWriter, source []byte, n ast.Node) error {
	l := n.Lines().Len()
	for i := range l {
		line := n.Lines().At(i)
		r.Writer.RawWrite(w, line.Value(source))
	}
	return nil
}

// NewExtension creates a new goldmark.Extender for template support
func NewExtension() goldmark.Extender {
	return &Extension{}
}

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
