package goldmarktemplate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func TestTableExtensionIntegration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "table with template actions in cells",
			input: `| Name | Value |
|------|-------|
| User | {{ .UserName }} |
| Age  | {{ .UserAge }} |`,
			expected: `<table>
<thead>
<tr>
<th>Name</th>
<th>Value</th>
</tr>
</thead>
<tbody>
<tr>
<td>User</td>
<td>{{ .UserName }}</td>
</tr>
<tr>
<td>Age</td>
<td>{{ .UserAge }}</td>
</tr>
</tbody>
</table>`,
		},
		{
			name: "table with template actions in headers",
			input: `| {{ .ColumnA }} | {{ .ColumnB }} |
|----------------|----------------|
| Value 1        | Value 2        |`,
			expected: `<table>
<thead>
<tr>
<th>{{ .ColumnA }}</th>
<th>{{ .ColumnB }}</th>
</tr>
</thead>
<tbody>
<tr>
<td>Value 1</td>
<td>Value 2</td>
</tr>
</tbody>
</table>`,
		},
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			New(),
			extension.Table,
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
			html.WithXHTML(),
		),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := md.Convert([]byte(tt.input), &buf)
			if err != nil {
				t.Fatalf("Failed to convert markdown: %v", err)
			}

			got := strings.TrimSpace(buf.String())
			expected := strings.TrimSpace(tt.expected)

			if got != expected {
				t.Errorf("Output mismatch\nInput:    %q\nExpected: %q\nGot:      %q", tt.input, expected, got)
			}
		})
	}
}

func TestStrikethroughExtensionIntegration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "strikethrough with template action",
			input:    "~~{{ .DeletedText }}~~",
			expected: "<p><del>{{ .DeletedText }}</del></p>",
		},
		{
			name:     "strikethrough with mixed content",
			input:    "~~Delete {{ .Item }} now~~",
			expected: "<p><del>Delete {{ .Item }} now</del></p>",
		},
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			New(),
			extension.Strikethrough,
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
			html.WithXHTML(),
		),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := md.Convert([]byte(tt.input), &buf)
			if err != nil {
				t.Fatalf("Failed to convert markdown: %v", err)
			}

			got := strings.TrimSpace(buf.String())
			expected := strings.TrimSpace(tt.expected)

			if got != expected {
				t.Errorf("Output mismatch\nInput:    %q\nExpected: %q\nGot:      %q", tt.input, expected, got)
			}
		})
	}
}

func TestTaskListExtensionIntegration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "task list with template actions",
			input: `- [x] {{ .CompletedTask }}
- [ ] {{ .PendingTask }}`,
			expected: `<ul>
<li><input checked="" disabled="" type="checkbox" /> {{ .CompletedTask }}</li>
<li><input disabled="" type="checkbox" /> {{ .PendingTask }}</li>
</ul>`,
		},
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			New(),
			extension.TaskList,
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
			html.WithXHTML(),
		),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := md.Convert([]byte(tt.input), &buf)
			if err != nil {
				t.Fatalf("Failed to convert markdown: %v", err)
			}

			got := strings.TrimSpace(buf.String())
			expected := strings.TrimSpace(tt.expected)

			if got != expected {
				t.Errorf("Output mismatch\nInput:    %q\nExpected: %q\nGot:      %q", tt.input, expected, got)
			}
		})
	}
}

func TestFootnoteExtensionIntegration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "footnote content with template action",
			input: `Here is a footnote[^1].

[^1]: {{ .FootnoteContent }}`,
			expected: `<p>Here is a footnote<sup id="fnref:1"><a href="#fn:1" class="footnote-ref" role="doc-noteref">1</a></sup>.</p>
<div class="footnotes" role="doc-endnotes">
<hr />
<ol>
<li id="fn:1">
<p>{{ .FootnoteContent }}&#160;<a href="#fnref:1" class="footnote-backref" role="doc-backlink">&#x21a9;&#xfe0e;</a></p>
</li>
</ol>
</div>`,
		},
		{
			name: "footnote with template in content",
			input: `Text with footnote[^1].

[^1]: Note {{ .Var }}`,
			expected: `<p>Text with footnote<sup id="fnref:1"><a href="#fn:1" class="footnote-ref" role="doc-noteref">1</a></sup>.</p>
<div class="footnotes" role="doc-endnotes">
<hr />
<ol>
<li id="fn:1">
<p>Note {{ .Var }}&#160;<a href="#fnref:1" class="footnote-backref" role="doc-backlink">&#x21a9;&#xfe0e;</a></p>
</li>
</ol>
</div>`,
		},
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			New(),
			extension.Footnote,
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
			html.WithXHTML(),
		),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := md.Convert([]byte(tt.input), &buf)
			if err != nil {
				t.Fatalf("Failed to convert markdown: %v", err)
			}

			got := strings.TrimSpace(buf.String())
			expected := strings.TrimSpace(tt.expected)

			if got != expected {
				t.Errorf("Output mismatch\nInput:    %q\nExpected: %q\nGot:      %q", tt.input, expected, got)
			}
		})
	}
}

func TestDefinitionListExtensionIntegration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "definition list with template actions",
			input: `{{ .Term }}
: {{ .Definition }}

Another term
: Definition with {{ .Variable }}`,
			expected: `<dl>
<dt>{{ .Term }}</dt>
<dd>{{ .Definition }}</dd>
<dt>Another term</dt>
<dd>Definition with {{ .Variable }}</dd>
</dl>`,
		},
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			New(),
			extension.DefinitionList,
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
			html.WithXHTML(),
		),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := md.Convert([]byte(tt.input), &buf)
			if err != nil {
				t.Fatalf("Failed to convert markdown: %v", err)
			}

			got := strings.TrimSpace(buf.String())
			expected := strings.TrimSpace(tt.expected)

			if got != expected {
				t.Errorf("Output mismatch\nInput:    %q\nExpected: %q\nGot:      %q", tt.input, expected, got)
			}
		})
	}
}

func TestLinkifyExtensionIntegration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "URL followed by template action",
			input: `Visit https://example.com for {{ .Info }}`,
			expected: `<p>Visit <a href="https://example.com">https://example.com</a> for {{ .Info }}</p>`,
		},
		{
			name: "template action with URL-like content",
			input: `Check out {{ .BaseURL }}/path/to/resource`,
			expected: `<p>Check out {{ .BaseURL }}/path/to/resource</p>`,
		},
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			New(),
			extension.Linkify,
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
			html.WithXHTML(),
		),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := md.Convert([]byte(tt.input), &buf)
			if err != nil {
				t.Fatalf("Failed to convert markdown: %v", err)
			}

			got := strings.TrimSpace(buf.String())
			expected := strings.TrimSpace(tt.expected)

			if got != expected {
				t.Errorf("Output mismatch\nInput:    %q\nExpected: %q\nGot:      %q", tt.input, expected, got)
			}
		})
	}
}

func TestTypographerExtensionIntegration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "typographer with template actions",
			input: `"{{ .Quote }}" said {{ .Author }}.`,
			expected: `<p>&ldquo;{{ .Quote }}&rdquo; said {{ .Author }}.</p>`,
		},
		{
			name: "en dashes with template actions", 
			input: `{{ .Title }} -- {{ .Subtitle }}`,
			expected: `<p>{{ .Title }} &ndash; {{ .Subtitle }}</p>`,
		},
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			New(),
			extension.Typographer,
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
			html.WithXHTML(),
		),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := md.Convert([]byte(tt.input), &buf)
			if err != nil {
				t.Fatalf("Failed to convert markdown: %v", err)
			}

			got := strings.TrimSpace(buf.String())
			expected := strings.TrimSpace(tt.expected)

			if got != expected {
				t.Errorf("Output mismatch\nInput:    %q\nExpected: %q\nGot:      %q", tt.input, expected, got)
			}
		})
	}
}

func TestCJKExtensionIntegration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "CJK text with template actions",
			input: `こんにちは{{ .Name }}さん`,
			expected: `<p>こんにちは{{ .Name }}さん</p>`,
		},
		{
			name: "Chinese text with template",
			input: `你好{{ .World }}世界`,
			expected: `<p>你好{{ .World }}世界</p>`,
		},
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			New(),
			extension.CJK,
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
			html.WithXHTML(),
		),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := md.Convert([]byte(tt.input), &buf)
			if err != nil {
				t.Fatalf("Failed to convert markdown: %v", err)
			}

			got := strings.TrimSpace(buf.String())
			expected := strings.TrimSpace(tt.expected)

			if got != expected {
				t.Errorf("Output mismatch\nInput:    %q\nExpected: %q\nGot:      %q", tt.input, expected, got)
			}
		})
	}
}