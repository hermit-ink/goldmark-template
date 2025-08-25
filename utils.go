package goldmarktemplate

import (
	"bytes"
)

// Template detection utilities
var (
	actionPattern = []byte("{{")
)

// hasAction checks if content contains template actions
func hasAction(content []byte) bool {
	return bytes.Contains(content, actionPattern)
}
