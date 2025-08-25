package ast

import (
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// TemplateAction represents a Go template action like {{...}}
type TemplateAction struct {
	gast.BaseInline
	Segment text.Segment
	Content []byte
}

// Dump implements Node.Dump.
func (n *TemplateAction) Dump(source []byte, level int) {
	gast.DumpHelper(n, source, level, nil, nil)
}

// KindTemplateAction is a NodeKind of the TemplateAction node.
var KindTemplateAction = gast.NewNodeKind("TemplateAction")

// Kind implements Node.Kind.
func (n *TemplateAction) Kind() gast.NodeKind {
	return KindTemplateAction
}

// NewTemplateAction returns a new TemplateAction node.
func NewTemplateAction(content []byte, segment text.Segment) *TemplateAction {
	return &TemplateAction{
		Content: content,
		Segment: segment,
	}
}
