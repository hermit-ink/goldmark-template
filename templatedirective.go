package goldmarktemplate

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// TemplateDirective represents a Go template directive like {{...}}
type TemplateDirective struct {
	ast.BaseInline
	Segment text.Segment
	Content []byte
}

// Dump implements Node.Dump.
func (n *TemplateDirective) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

// KindTemplateDirective is a NodeKind of the TemplateDirective node.
var KindTemplateDirective = ast.NewNodeKind("TemplateDirective")

// Kind implements Node.Kind.
func (n *TemplateDirective) Kind() ast.NodeKind {
	return KindTemplateDirective
}

// NewTemplateDirective returns a new TemplateDirective node.
func NewTemplateDirective(content []byte, segment text.Segment) *TemplateDirective {
	return &TemplateDirective{
		Content: content,
		Segment: segment,
	}
}

// templateDirectiveParser is an inline parser for Go template directives
type templateDirectiveParser struct{}

// NewTemplateDirectiveParser returns a new InlineParser that parses template directives
func NewTemplateDirectiveParser() parser.InlineParser {
	return &templateDirectiveParser{}
}

// Trigger returns characters that trigger this parser
func (s *templateDirectiveParser) Trigger() []byte {
	return []byte{'{'}
}

func (s *templateDirectiveParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, segment := block.PeekLine()

	if len(line) < 2 || line[0] != '{' || line[1] != '{' {
		return nil
	}

	i := 2
	insideDoubleQuote := false
	insideSingleQuote := false

	for i < len(line)-1 {
		char := line[i]

		if char == '"' && !insideSingleQuote {
			insideDoubleQuote = !insideDoubleQuote
		} else if char == '\'' && !insideDoubleQuote {
			insideSingleQuote = !insideSingleQuote
		}

		if !insideDoubleQuote && !insideSingleQuote && char == '}' && line[i+1] == '}' {
			content := line[0 : i+2]
			nodeSegment := segment.WithStop(segment.Start + i + 2)
			node := NewTemplateDirective(content, nodeSegment)
			block.Advance(i + 2)
			return node
		}
		i++
	}

	return nil
}

func (s *templateDirectiveParser) CloseBlock(parent ast.Node, pc parser.Context) {
	// nothing to do
}

// TemplateDirectiveHTMLRenderer renders TemplateDirective nodes
type TemplateDirectiveHTMLRenderer struct {
	html.Config
}

// NewTemplateDirectiveHTMLRenderer returns a new TemplateDirectiveHTMLRenderer
func NewTemplateDirectiveHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &TemplateDirectiveHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs
func (r *TemplateDirectiveHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindTemplateDirective, r.renderTemplateDirective)
}

// renderTemplateDirective renders template directive as raw content (no HTML encoding)
func (r *TemplateDirectiveHTMLRenderer) renderTemplateDirective(
	w util.BufWriter, source []byte, n ast.Node, entering bool,
) (ast.WalkStatus, error) {
	if entering {
		if node, ok := n.(*TemplateDirective); ok {
			// Write the template directive as-is (no HTML encoding)
			_, err := w.Write(node.Content)
			if err != nil {
				return ast.WalkStop, err
			}
		}
	}
	return ast.WalkContinue, nil
}
