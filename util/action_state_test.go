package util

import (
	"testing"
)

func TestActionStateBasics(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []bool // expected InAction state at each position
	}{
		{
			name:     "simple template action",
			input:    "{{ .Var }}",
			expected: []bool{true, true, true, true, true, true, true, true, false, false},
		},
		{
			name:     "template action with text before and after",
			input:    "before {{ .Var }} after",
			expected: []bool{false, false, false, false, false, false, false, true, true, true, true, true, true, true, true, false, false, false, false, false, false, false, false},
		},
		{
			name:     "nested template actions",
			input:    "{{ if {{ .Condition }} }}",
			expected: []bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, false, false},
		},
		{
			name:     "no template actions",
			input:    "plain text",
			expected: []bool{false, false, false, false, false, false, false, false, false, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracker := NewActionState()
			line := []byte(tt.input)
			
			if len(tt.expected) != len(line) {
				t.Fatalf("Test setup error: expected slice length %d doesn't match input length %d", len(tt.expected), len(line))
			}

			for i := 0; i < len(line); i++ {
				tracker.ProcessChar(line, i)
				if tracker.InAction() != tt.expected[i] {
					t.Errorf("Position %d in %q: expected inAction=%v, got %v", i, tt.input, tt.expected[i], tracker.InAction())
				}
			}
		})
	}
}

func TestActionStateWithQuotes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []bool
	}{
		{
			name:     "double quotes in template action",
			input:    `{{ "hello" }}`,
			expected: []bool{true, true, true, true, true, true, true, true, true, true, true, false, false},
		},
		{
			name:     "single quotes in template action", 
			input:    `{{ 'hello' }}`,
			expected: []bool{true, true, true, true, true, true, true, true, true, true, true, false, false},
		},
		{
			name:     "backticks in template action",
			input:    "{{ `hello` }}",
			expected: []bool{true, true, true, true, true, true, true, true, true, true, true, false, false},
		},
		{
			name:     "closing braces inside double quotes",
			input:    `{{ "}} not closed" }}`,
			expected: []bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, false, false},
		},
		{
			name:     "closing braces inside single quotes",
			input:    `{{ '}} not closed' }}`,
			expected: []bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, false, false},
		},
		{
			name:     "closing braces inside backticks",
			input:    "{{ `}} not closed` }}",
			expected: []bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, false, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracker := NewActionState()
			line := []byte(tt.input)
			
			if len(tt.expected) != len(line) {
				t.Fatalf("Test setup error: expected slice length %d doesn't match input length %d", len(tt.expected), len(line))
			}

			for i := 0; i < len(line); i++ {
				tracker.ProcessChar(line, i)
				if tracker.InAction() != tt.expected[i] {
					t.Errorf("Position %d in %q: expected inAction=%v, got %v", i, tt.input, tt.expected[i], tracker.InAction())
				}
			}
		})
	}
}

func TestActionStateWithEscapedCharacters(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []bool
	}{
		{
			name:     "escaped double quote",
			input:    `{{ "hello \"world\"" }}`,
			expected: []bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, false, false},
		},
		{
			name:     "escaped single quote",
			input:    `{{ 'hello \'world\'' }}`,
			expected: []bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, false, false},
		},
		{
			name:     "escaped backslash",
			input:    `{{ "hello \\" }}`,
			expected: []bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, false, false},
		},
		{
			name:     "double escaped backslash",
			input:    `{{ "hello \\\\" }}`,
			expected: []bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, false, false},
		},
		{
			name:     "escaped closing braces",
			input:    `{{ "\\}}" }}`,
			expected: []bool{true, true, true, true, true, true, true, true, true, true, false, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracker := NewActionState()
			line := []byte(tt.input)
			
			if len(tt.expected) != len(line) {
				t.Fatalf("Test setup error: expected slice length %d doesn't match input length %d", len(tt.expected), len(line))
			}

			for i := 0; i < len(line); i++ {
				tracker.ProcessChar(line, i)
				if tracker.InAction() != tt.expected[i] {
					t.Errorf("Position %d in %q: expected inAction=%v, got %v", i, tt.input, tt.expected[i], tracker.InAction())
				}
			}
		})
	}
}

func TestActionStateComplexCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []bool
	}{
		{
			name:     "multiple template actions",
			input:    "{{ .A }} and {{ .B }}",
			expected: []bool{true, true, true, true, true, true, false, false, false, false, false, false, false, true, true, true, true, true, true, false, false},
		},
		{
			name:     "nested braces with quotes",
			input:    `{{ if (gt {{ .Count }} 0) "yes" "no" }}`,
			expected: []bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, false, false},
		},
		{
			name:     "template action with function call",
			input:    `{{ printf "Hello %s" .Name }}`,
			expected: []bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, false, false},
		},
		{
			name:     "template action with pipe",
			input:    `{{ .Value | upper }}`,
			expected: []bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, false, false},
		},
		{
			name:     "template action with range",
			input:    `{{ range .Items }}{{ . }}{{ end }}`,
			expected: []bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, false, false, true, true, true, true, true, false, false, true, true, true, true, true, true, true, false, false},
		},
		{
			name:     "template action with complex backslash escaping",
			input:    `{{ .Path | replace "\\\\" "/" }}`,
			expected: []bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, false, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracker := NewActionState()
			line := []byte(tt.input)
			
			if len(tt.expected) != len(line) {
				t.Fatalf("Test setup error: expected slice length %d doesn't match input length %d", len(tt.expected), len(line))
			}

			for i := 0; i < len(line); i++ {
				tracker.ProcessChar(line, i)
				if tracker.InAction() != tt.expected[i] {
					t.Errorf("Position %d in %q: expected inAction=%v, got %v", i, tt.input, tt.expected[i], tracker.InAction())
				}
			}
		})
	}
}

func TestFindActionEnd(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		startPos int
		expected int
	}{
		{
			name:     "simple action",
			input:    "{{ .Var }}",
			startPos: 0,
			expected: 10,
		},
		{
			name:     "action with quotes",
			input:    `{{ "hello" }}`,
			startPos: 0,
			expected: 13,
		},
		{
			name:     "action with nested braces in quotes",
			input:    `{{ "}} not end" }}`,
			startPos: 0,
			expected: 18,
		},
		{
			name:     "action in middle of string",
			input:    "prefix {{ .Var }} suffix",
			startPos: 7,
			expected: 17,
		},
		{
			name:     "nested actions",
			input:    "{{ if {{ .Cond }} }}",
			startPos: 0,
			expected: 20,
		},
		{
			name:     "invalid start position",
			input:    "no action here",
			startPos: 0,
			expected: -1,
		},
		{
			name:     "incomplete action",
			input:    "{{ incomplete",
			startPos: 0,
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindActionEnd([]byte(tt.input), tt.startPos)
			if result != tt.expected {
				t.Errorf("FindActionEnd(%q, %d): expected %d, got %d", tt.input, tt.startPos, tt.expected, result)
			}
		})
	}
}

func TestActionStateEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: []bool{},
		},
		{
			name:     "single brace",
			input:    "{",
			expected: []bool{false},
		},
		{
			name:     "double brace at end",
			input:    "text{{",
			expected: []bool{false, false, false, false, true, true},
		},
		{
			name:     "incomplete closing",
			input:    "{{ .Var }",
			expected: []bool{true, true, true, true, true, true, true, true, true},
		},
		{
			name:     "only opening braces",
			input:    "{{",
			expected: []bool{true, true},
		},
		{
			name:     "only closing braces",
			input:    "}}",
			expected: []bool{false, false},
		},
		{
			name:     "mixed braces",
			input:    "{{{ }}}",
			expected: []bool{true, true, true, true, true, false, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracker := NewActionState()
			line := []byte(tt.input)
			
			if len(tt.expected) != len(line) {
				t.Fatalf("Test setup error: expected slice length %d doesn't match input length %d", len(tt.expected), len(line))
			}

			for i := 0; i < len(line); i++ {
				tracker.ProcessChar(line, i)
				if tracker.InAction() != tt.expected[i] {
					t.Errorf("Position %d in %q: expected inAction=%v, got %v", i, tt.input, tt.expected[i], tracker.InAction())
				}
			}
		})
	}
}