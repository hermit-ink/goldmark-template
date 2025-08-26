package util

import "bytes"

// ContainsAction checks if the given content contains Go template actions
func ContainsAction(content []byte) bool {
	return bytes.Contains(content, []byte("{{"))
}

// actionState tracks state when parsing through template actions
type actionState struct {
	inAction        bool
	templateDepth   int
	inDoubleQuotes  bool
	inSingleQuotes  bool
	inBackticks     bool
}

// NewActionState creates a new template action state tracker
func NewActionState() *actionState {
	return &actionState{}
}

// ProcessChar processes a character and updates template action state
// Returns true if the character should be ignored for other parsing logic
func (t *actionState) ProcessChar(line []byte, i int) bool {
	if i >= len(line) {
		return false
	}

	char := line[i]

	// Check if current character is escaped (but not inside backticks where escapes don't exist)
	isEscaped := !t.inBackticks && i > 0 && line[i-1] == '\\'

	// Track template action boundaries
	if !t.inAction && i < len(line)-1 && char == '{' && line[i+1] == '{' {
		t.inAction = true
		t.templateDepth = 1
		return false
	} else if t.inAction {
		// Track quote/backtick state within template actions
		if !isEscaped {
			if char == '"' && !t.inSingleQuotes && !t.inBackticks {
				t.inDoubleQuotes = !t.inDoubleQuotes
			} else if char == '\'' && !t.inDoubleQuotes && !t.inBackticks {
				t.inSingleQuotes = !t.inSingleQuotes
			} else if char == '`' && !t.inDoubleQuotes && !t.inSingleQuotes {
				t.inBackticks = !t.inBackticks
			}
		}

		// Check for template action end only when not inside quotes/backticks
		if !t.inDoubleQuotes && !t.inSingleQuotes && !t.inBackticks &&
		   i < len(line)-1 && char == '}' && line[i+1] == '}' {
			t.templateDepth--
			if t.templateDepth == 0 {
				t.inAction = false
			}
		} else if !t.inDoubleQuotes && !t.inSingleQuotes && !t.inBackticks &&
		          i < len(line)-1 && char == '{' && line[i+1] == '{' {
			t.templateDepth++
		}
	}

	return false
}

// InTemplateAction returns true if currently inside a template action
func (t *actionState) InAction() bool {
	return t.inAction
}

// FindActionEnd finds the end of a template action starting from position startPos
// Returns the position after the closing }} or -1 if not found
func FindActionEnd(line []byte, startPos int) int {
	if startPos+2 >= len(line) || line[startPos] != '{' || line[startPos+1] != '{' {
		return -1
	}

	tracker := NewActionState()

	tracker.ProcessChar(line, startPos)

	for i := startPos + 1; i < len(line)-1; i++ {
		tracker.ProcessChar(line, i)

		if !tracker.InAction() {
			if i < len(line)-1 && line[i] == '}' && line[i+1] == '}' {
				return i + 2
			}
		}
	}

	return -1
}
