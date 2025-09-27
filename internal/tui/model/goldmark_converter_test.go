package model

import (
	"fmt"
	"testing"
)

func TestConvertMarkdownToElements(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expected int // Expected number of top-level elements
	}{
		{
			name:     "Simple paragraph",
			markdown: "Hello world",
			expected: 1,
		},
		{
			name:     "Heading",
			markdown: "# Hello World",
			expected: 1,
		},
		{
			name:     "Multiple paragraphs",
			markdown: "First paragraph\n\nSecond paragraph",
			expected: 2,
		},
		{
			name:     "Code block",
			markdown: "```go\nfunc main() {\n    fmt.Println(\"Hello\")\n}\n```",
			expected: 1,
		},
		{
			name:     "List",
			markdown: "- Item 1\n- Item 2\n- Item 3",
			expected: 1,
		},
		{
			name:     "Mixed content",
			markdown: "# Title\n\nThis is a paragraph with **bold** and *italic* text.\n\n- List item 1\n- List item 2",
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elements := ConvertMarkdownToElements(tt.markdown)
			if len(elements) != tt.expected {
				t.Errorf("Expected %d elements, got %d", tt.expected, len(elements))
			}
		})
	}
}

func TestHeadingElement(t *testing.T) {
	markdown := "# Level 1\n## Level 2\n### Level 3"
	elements := ConvertMarkdownToElements(markdown)

	if len(elements) != 3 {
		t.Fatalf("Expected 3 elements, got %d", len(elements))
	}

	// Check first heading
	if heading, ok := elements[0].(*HeadingElement); ok {
		if heading.Level != 1 {
			t.Errorf("Expected level 1, got %d", heading.Level)
		}
	} else {
		t.Error("Expected HeadingElement")
	}

	// Check second heading
	if heading, ok := elements[1].(*HeadingElement); ok {
		if heading.Level != 2 {
			t.Errorf("Expected level 2, got %d", heading.Level)
		}
	} else {
		t.Error("Expected HeadingElement")
	}
}

func TestParagraphWithInlineElements(t *testing.T) {
	markdown := "This is **bold** and *italic* text with `code`."
	elements := ConvertMarkdownToElements(markdown)

	if len(elements) != 1 {
		t.Fatalf("Expected 1 element, got %d", len(elements))
	}

	paragraph, ok := elements[0].(*ParagraphElement)
	if !ok {
		t.Fatal("Expected ParagraphElement")
	}

	if paragraph.GetElementCount() == 0 {
		t.Error("Expected paragraph to have child elements")
	}
}

func TestCodeBlockElement(t *testing.T) {
	markdown := "```go\nfunc main() {\n    fmt.Println(\"Hello\")\n}\n```"
	elements := ConvertMarkdownToElements(markdown)

	if len(elements) != 1 {
		t.Fatalf("Expected 1 element, got %d", len(elements))
	}

	codeBlock, ok := elements[0].(*CodeBlockElement)
	if !ok {
		t.Fatal("Expected CodeBlockElement")
	}

	if codeBlock.Language != "go" {
		t.Errorf("Expected language 'go', got '%s'", codeBlock.Language)
	}

	if codeBlock.Text == "" {
		t.Error("Expected code block to have text content")
	}
}

func TestListElement(t *testing.T) {
	markdown := "1. First item\n2. Second item\n3. Third item"
	elements := ConvertMarkdownToElements(markdown)

	if len(elements) != 1 {
		t.Fatalf("Expected 1 element, got %d", len(elements))
	}

	list, ok := elements[0].(*ListElement)
	if !ok {
		t.Fatal("Expected ListElement")
	}

	if !list.Ordered {
		t.Error("Expected ordered list")
	}

	if list.GetElementCount() != 3 {
		t.Errorf("Expected 3 list items, got %d", list.GetElementCount())
	}
}

func TestLinkElement(t *testing.T) {
	markdown := "[Google](https://google.com)"
	elements := ConvertMarkdownToElements(markdown)

	if len(elements) != 1 {
		t.Fatalf("Expected 1 element, got %d", len(elements))
	}

	paragraph, ok := elements[0].(*ParagraphElement)
	if !ok {
		t.Fatal("Expected ParagraphElement")
	}

	if paragraph.GetElementCount() == 0 {
		t.Fatal("Expected paragraph to have child elements")
	}

	link, ok := paragraph.GetElement(0).(*LinkElement)
	if !ok {
		t.Fatal("Expected LinkElement")
	}

	if link.URL != "https://google.com" {
		t.Errorf("Expected URL 'https://google.com', got '%s'", link.URL)
	}
}

func TestPlainDocumentRender(t *testing.T) {
	doc := &PlainDocument{
		// Styled: true,
	}
	doc.Init()

	// Insert some markdown content
	doc.InsertString("# Hello World\n\nThis is a **test** document.")

	elements := doc.Render()

	if len(elements) == 0 {
		t.Error("Expected elements to be rendered")
	}

	// Should have heading and paragraph
	if len(elements) < 2 {
		t.Errorf("Expected at least 2 elements, got %d", len(elements))
	}
}

/*
func TestPlainDocumentRenderPlainText(t *testing.T) {
	doc := &PlainDocument{
		Styled: false, // Plain text mode
	}
	doc.Init()

	doc.InsertString("Hello World")

	elements := doc.Render()

	if len(elements) != 1 {
		t.Errorf("Expected 1 element, got %d", len(elements))
	}

	// Should be LineContainerElement in plain text mode
	if _, ok := elements[0].(*LineContainerElement); !ok {
		t.Error("Expected LineContainerElement in plain text mode")
	}
}
*/

func TestCompplex(t *testing.T) {
	code := `
# Markdown syntax guide

## Headers

# This is a Heading h1
## This is a Heading h2
###### This is a Heading h6

## Emphasis

*This text will be italic*  
_This will also be italic_

**This text will be bold**  
__This will also be bold__

_You **can** combine them_

## Lists

### Unordered

* Item
  * Item
    * Item
      * [ ] Item
      * [x] Item
    * Item

### Ordered

1. Item 1
2. Item 2
3. Item 3
  1. [ ] Item 3a
  2. Item 3b

## Images

![This is an alt text.](/image/sample.webp "This is a sample image.")

## Links

You may be using [Markdown Live Preview](https://markdownlivepreview.com/).

## Blockquotes

> Markdown is a lightweight markup language with plain-text-formatting syntax, created in 2004 by John Gruber with Aaron Swartz.
>
>> Markdown is often used to format readme files, for writing messages in online discussion forums, and to create rich text using a plain text editor.

## Tables

| Left columns  | Right columns |
| :--- | ---:|
| left foo      | right foo     |
| left bar      | right bar     |
| left baz      | right baz     |

## Blocks of code
`

	ConvertMarkdownToElements(code)
}

func TestTable(t *testing.T) {
	code := `| Functional option | Type | Description |
| ----------------- | ---- | ----------- |
| goldmark.WithParser | parser.Parser  | This option must be passed before goldmark.WithParserOptions and goldmark.WithExtensions |
| goldmark.WithRenderer | renderer.Renderer  | This option must be passed before goldmark.WithRendererOptions and goldmark.WithExtensions  |
| goldmark.WithParserOptions | ...parser.Option  |  |
| goldmark.WithRendererOptions | ...renderer.Option |  |
| goldmark.WithExtensions | ...goldmark.Extender  |  |
`

	elements := ConvertMarkdownToElements(code)
	if table, ok := elements[0].(*TableElement); ok {
		for _, child := range table.Children {
			fmt.Printf("%#v\n", child)
		}
	}
}
