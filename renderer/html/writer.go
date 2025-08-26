package html

import (
	"bytes"

	tutil "github.com/hermit-ink/goldmark-template/util"
	ghtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// Template detection utilities
var (
	actionPattern = []byte("{{")
)

// hasAction checks if content contains template actions
func hasAction(content []byte) bool {
	return bytes.Contains(content, actionPattern)
}

// Writer is a custom HTML writer that preserves Go template actions
// without HTML escaping them, while properly handling escaped template cases.
type Writer struct {
	fallback ghtml.Writer
}

// NewWriter creates a new Writer
func NewWriter(opts ...ghtml.WriterOption) ghtml.Writer {
	return &Writer{
		fallback: ghtml.NewWriter(opts...),
	}
}

// Write writes content with normal processing (includes entity resolution and backslash unescaping)
func (w *Writer) Write(writer util.BufWriter, source []byte) {
	if hasAction(source) {
		w.writeWithTemplateSupport(writer, source, true)
	} else {
		w.fallback.Write(writer, source)
	}
}

// SecureWrite writes content with security filtering
func (w *Writer) SecureWrite(writer util.BufWriter, source []byte) {
	if hasAction(source) {
		w.writeWithTemplateSupport(writer, source, false)
	} else {
		w.fallback.SecureWrite(writer, source)
	}
}

// RawWrite writes content while preserving Go template actions (HTML escaping only)
func (w *Writer) RawWrite(writer util.BufWriter, source []byte) {
	if hasAction(source) {
		w.writeWithTemplateSupport(writer, source, false)
	} else {
		w.fallback.RawWrite(writer, source)
	}
}

// writeWithTemplateSupport handles content with template actions
func (w *Writer) writeWithTemplateSupport(writer util.BufWriter, source []byte, processEntities bool) {
	n := 0
	i := 0

	for i < len(source) {
		// Skip non-template characters
		if i >= len(source)-1 || !bytes.HasPrefix(source[i:], actionPattern) {
			i++
			continue
		}

		// Process everything before the action using goldmark's methods
		if n < i {
			beforeAction := source[n:i]
			if processEntities {
				w.fallback.Write(writer, beforeAction)
			} else {
				w.fallback.RawWrite(writer, beforeAction)
			}
		}

		// Find and write the complete template action verbatim
		end := tutil.FindActionEnd(source, i)
		if end <= 0 {
			i++
			continue
		}

		if _, err := writer.Write(source[i:end]); err != nil {
			return
		}
		n = end
		i = end
	}

	// Process remaining content using goldmark's methods
	if n < len(source) {
		remaining := source[n:]
		if processEntities {
			w.fallback.Write(writer, remaining)
		} else {
			w.fallback.RawWrite(writer, remaining)
		}
	}
}

