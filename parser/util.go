package parser

import "bytes"

// containsTemplateAction checks if the given content contains Go template actions
func containsTemplateAction(content []byte) bool {
	return bytes.Contains(content, []byte("{{"))
}