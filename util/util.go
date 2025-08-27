package util

import (
	"bytes"
)

// ContainsAction checks if the given content contains Go template actions
func ContainsAction(content []byte) bool {
	return bytes.Contains(content, []byte("{{"))
}
