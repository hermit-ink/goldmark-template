package goldmarktemplate

import (
	"bytes"

	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

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