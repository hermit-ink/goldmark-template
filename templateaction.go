package goldmarktemplate

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// TemplateAction represents a Go template action like {{...}}
type TemplateAction struct {
	ast.BaseInline
	Segment text.Segment
	Content []byte
}

// Dump implements Node.Dump.
func (n *TemplateAction) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

// KindTemplateAction is a NodeKind of the TemplateAction node.
var KindTemplateAction = ast.NewNodeKind("TemplateAction")

// Kind implements Node.Kind.
func (n *TemplateAction) Kind() ast.NodeKind {
	return KindTemplateAction
}

// NewTemplateAction returns a new TemplateAction node.
func NewTemplateAction(content []byte, segment text.Segment) *TemplateAction {
	return &TemplateAction{
		Content: content,
		Segment: segment,
	}
}

// templateActionParser is an inline parser for Go template actions
type templateActionParser struct{}

// NewTemplateActionParser returns a new InlineParser that parses go template
// actions
func NewTemplateActionParser() parser.InlineParser {
	return &templateActionParser{}
}

// Trigger returns characters that trigger this parser
func (s *templateActionParser) Trigger() []byte {
	return []byte{'{'}
}

func (s *templateActionParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, segment := block.PeekLine()

	if len(line) < 2 || line[0] != '{' || line[1] != '{' {
		return nil
	}

	i := 2
	inDoubleQuotes := false
	inSingleQuotes := false

	for i < len(line)-1 {
		char := line[i]

		if char == '"' && !inSingleQuotes {
			inDoubleQuotes = !inDoubleQuotes
		} else if char == '\'' && !inDoubleQuotes {
			inSingleQuotes = !inSingleQuotes
		}

		if !inDoubleQuotes && !inSingleQuotes && char == '}' && line[i+1] == '}' {
			content := line[0 : i+2]
			nodeSegment := segment.WithStop(segment.Start + i + 2)
			node := NewTemplateAction(content, nodeSegment)
			block.Advance(i + 2)
			return node
		}
		i++
	}

	return nil
}

func (s *templateActionParser) CloseBlock(parent ast.Node, pc parser.Context) {
	// nothing to do
}

// TemplateActionHTMLRenderer renders TemplateAction nodes directly into the
// output with no HTML/URL escaping
type TemplateActionHTMLRenderer struct {
	html.Config
}

// NewTemplateActionHTMLRenderer returns a new TemplateActionHTMLRenderer
func NewTemplateActionHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &TemplateActionHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs
func (r *TemplateActionHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindTemplateAction, r.render)
}

// render renders template actions as raw content (no HTML encoding)
func (r *TemplateActionHTMLRenderer) render(
	w util.BufWriter, source []byte, n ast.Node, entering bool,
) (ast.WalkStatus, error) {
	if entering {
		if node, ok := n.(*TemplateAction); ok {
			// Write the template action as-is (no HTML encoding)
			_, err := w.Write(node.Content)
			if err != nil {
				return ast.WalkStop, err
			}
		}
	}
	return ast.WalkContinue, nil
}
