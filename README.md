# goldmark-template

[![Build Status](https://github.com/hermit-ink/goldmark-template/actions/workflows/test.yml/badge.svg)](https://github.com/hermit-ink/goldmark-template/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/hermit-ink/goldmark-template.svg)](https://pkg.go.dev/github.com/hermit-ink/goldmark-template)
[![Go Report Card](https://goreportcard.com/badge/github.com/hermit-ink/goldmark-template)](https://goreportcard.com/report/github.com/hermit-ink/goldmark-template)
[![Latest Release](https://img.shields.io/github/v/release/hermit-ink/goldmark-template)](https://github.com/hermit-ink/goldmark-template/releases)
[![License](https://img.shields.io/github/license/hermit-ink/goldmark-template)](LICENSE)

A [goldmark](https://github.com/yuin/goldmark) extension that preserves Go template actions (`{{...}}`) in rendered Markdown, preventing HTML escaping and maintaining template syntax wherever it appears so that the HTML output can be used directly by the Go stdlib html/template

## Motivation

You want to use go html/template actions in your markdown.  Markdown is not HTML so you can't execute the markdown as a template and you might not want to do that anyway.  Instead you want to goldmark to simply "ignore" all the go template actions putting them verbatim into the correct places in the output HTML so that you can execute the output from goldmark using go's html/template package.

> [!WARNING]
> This makes your goldmark instance no longer commonmark compliant since go template actions are not valid URLs and you want to be able to use a template action in place of a URL potentially.

## Features

-  **Preserves template actions** in inline code and code blocks
-  **Template-aware parsing** for links, images, and autolinks
-  **Reference link support** with template URLs and titles
-  **Standalone template actions** as inline elements
-  **Full compatibility** with other goldmark extensions (GFM, etc.)
-  **Smart parsing** that handles quotes and nested braces correctly
-  **Comprehensive testing** for 100% compatibility with the existing goldmark parsers and renderers

## Installation

```bash
go get github.com/hermit-ink/goldmark-template
```

## Usage

### Basic Usage

```go
package main

import (
    "bytes"
    "fmt"

    "github.com/yuin/goldmark"
    "github.com/yuin/goldmark/renderer/html"
    goldmarktemplate "github.com/hermit-ink/goldmark-template"
)

func main() {
    md := goldmark.New(
        goldmark.WithExtensions(
            goldmarktemplate.New(),
        ),
        goldmark.WithRendererOptions(
            html.WithUnsafe(), // Required for action preservation in raw HTML
        ),
    )

    input := []byte("# {{ .Title }}\n\n[Link]({{ .URL }})")
    var buf bytes.Buffer
    if err := md.Convert(input, &buf); err != nil {
        panic(err)
    }
    fmt.Println(buf.String())
    // Output: <h1>{{ .Title }}</h1>\n<p><a href="{{ .URL }}">Link</a></p>
}
```

### With Parser Options

Since goldmark-template *replaces* built-in parsers you need to use the goldmark-template
specific way of adding parser options.  The API is identical to goldmark's `WithParserOptions` except its an alternative constructor for the extensions.  Use it just like
you would goldmark's `WithParserOptions`.

```go
package main

import (
    "bytes"
    "fmt"

    "github.com/yuin/goldmark"
    "github.com/yuin/goldmark/parser"
    "github.com/yuin/goldmark/renderer/html"
    goldmarktemplate "github.com/hermit-ink/goldmark-template"
)

func main() {
    md := goldmark.New(
        goldmark.WithExtensions(
            goldmarktemplate.WithParserOptions(
                parser.WithAutoHeadingID(),
                parser.WithAttribute(),
            ),
        ),
        goldmark.WithRendererOptions(
            html.WithUnsafe(),
        ),
    )

    input := []byte("# {{ .Title }}\n\n[Link]({{ .URL }})")
    var buf bytes.Buffer
    if err := md.Convert(input, &buf); err != nil {
        panic(err)
    }
    fmt.Println(buf.String())
    // Output: <h1 id="title">{{ .Title }}</h1>\n<p><a href="{{ .URL }}">Link</a></p>
}
```

### With GFM Extension

```go
md := goldmark.New(
    goldmark.WithExtensions(
        goldmarktemplate.New(), // Must come FIRST
        extension.GFM,
    ),
    goldmark.WithRendererOptions(
        html.WithUnsafe(),
    ),
)
```

## Examples

### Template Actions in Code

`````markdown
Inline: `{{ .Variable }}`
Block:
```go
func main() {
    fmt.Println("{{ .Message }}")
}
```
`````

Output:
```html
<p>Inline: <code>{{ .Variable }}</code></p>
<pre><code class="language-go">func main() {
    fmt.Println("{{ .Message }}")
}
</code></pre>
```

### Template Actions in Links and Images

Input:
```markdown
[{{ .LinkText }}]({{ .URL }})
![{{ .Alt }}]({{ .ImagePath }})
<{{ .BaseURL }}/page>
```

Output:
```html
<p><a href="{{ .URL }}">{{ .LinkText }}</a></p>
<p><img src="{{ .ImagePath }}" alt="{{ .Alt }}"></p>
<p><a href="{{ .BaseURL }}/page">{{ .BaseURL }}/page</a></p>
```

### Reference Links with Templates

```markdown
[Example][ref]

[ref]: {{ .URL }} "{{ .Title }}"
```

Output:
```html
<p><a href="{{ .URL }}" title="{{ .Title }}">Example</a></p>
```

### Standalone Template Actions

```markdown
Welcome {{ .User.Name }}!

Today is {{ .Date }}.
```

Output:
```html
<p>Welcome {{ .User.Name }}!</p>
<p>Today is {{ .Date }}.</p>
```

## Limitations and Caveats

### Actions can only be used as values in attributes

```markdown
# Heading {id="{{ .HeadingID }}"}
# Heading {class="{{ .CSSClass }}"}
# Heading {data-value="{{ .Data }}"}
```

### Extension Order Matters
Always register `goldmark-template` **BEFORE** other extensions that might interfere with template syntax:

```go
goldmark.WithExtensions(
    goldmarktemplate.New(), // First
    extension.GFM,
)
```

### No Template Validation

This extension **does not validate** Go template syntax. Invalid templates pass
through unchanged.

```go
// 1. Process Markdown with goldmark-template
html := processMarkdown(markdown)

// 2. Process the HTML with Go templates
tmpl := template.Must(template.New("").Parse(html))
tmpl.Execute(w, data)
```

## Development

### Running Tests

```bash
make test
```

### Linting

```bash
make lint
```

### Code Coverage

```bash
make test-coverage
```

## Notes

The extension follows goldmark's established patterns:

- **Custom Parsers**: Template-aware parsers for links, autolinks, and reference
definitions.  These are taken directly from the goldmark source with the minimal
possible changes to allow template actions to be preserved untouched.
- **Custom Renderers**:
  - `Renderer` - Overrides standard elements to preserve template actions properly
  within attributes
  - `TemplateActionHTMLRenderer` - Renders standalone template actions
- **Custom AST Node**: `TemplateAction` for actions that do not appear in positions
controlled by other parsers such as images and links

## Contributing

Contributions are welcome! Please ensure:
1. All tests pass
2. Code is properly formatted (`make fmt`)
3. Linting passes (`make lint`)
4. New features include tests

## License

Apache License 2.0 - see [LICENSE](LICENSE) file for details.

## Acknowledgments

Built on top of the excellent [goldmark](https://github.com/yuin/goldmark)
Markdown parser by Yusuke Inuzuka.
