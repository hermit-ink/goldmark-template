package html

import (
	"bytes"

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

// Write writes content with normal processing
func (w *Writer) Write(writer util.BufWriter, source []byte) {
	w.fallback.Write(writer, source)
}

// SecureWrite writes content with security filtering
func (w *Writer) SecureWrite(writer util.BufWriter, source []byte) {
	w.fallback.SecureWrite(writer, source)
}

// RawWrite writes content while preserving Go template actions
func (w *Writer) RawWrite(writer util.BufWriter, source []byte) {
	n := 0
	i := 0

	for i < len(source) {
		if i < len(source)-1 && bytes.HasPrefix(source[i:], actionPattern) {
			// Write everything before the action (with escaping)
			if err := w.writeEscaped(writer, source[n:i]); err != nil {
				_ = w.writeEscaped(writer, source[i:])
				return
			}

			end := w.findActionEnd(source, i+2)
			if end > 0 {
				if _, err := writer.Write(source[i:end]); err != nil {
					return
				}
				n = end
				i = end
			} else {
				i++
			}
		} else {
			i++
		}
	}

	// Write remaining content (with escaping)
	if n < len(source) {
		_ = w.writeEscaped(writer, source[n:])
	}
}

func (w *Writer) writeEscaped(writer util.BufWriter, source []byte) error {
	for _, b := range source {
		switch b {
		case '<':
			if _, err := writer.WriteString("&lt;"); err != nil {
				return err
			}
		case '>':
			if _, err := writer.WriteString("&gt;"); err != nil {
				return err
			}
		case '&':
			if _, err := writer.WriteString("&amp;"); err != nil {
				return err
			}
		case '"':
			if _, err := writer.WriteString("&quot;"); err != nil {
				return err
			}
		default:
			if err := writer.WriteByte(b); err != nil {
				return err
			}
		}
	}
	return nil
}

func (w *Writer) findActionEnd(source []byte, start int) int {
	inDoubleQuotes := false
	inSingleQuotes := false

	for i := start; i < len(source)-1; i++ {
		char := source[i]

		if char == '"' && !inSingleQuotes {
			inDoubleQuotes = !inDoubleQuotes
		} else if char == '\'' && !inDoubleQuotes {
			inSingleQuotes = !inSingleQuotes
		}

		if !inDoubleQuotes && !inSingleQuotes && char == '}' && source[i+1] == '}' {
			return i + 2
		}
	}

	return -1
}
