package goldmarktemplate

import (
	"bytes"
)

// Template detection utilities
var (
	templatePattern = []byte("{{")
)

// hasTemplate checks if content contains template directives
func hasTemplate(content []byte) bool {
	return bytes.Contains(content, templatePattern)
}