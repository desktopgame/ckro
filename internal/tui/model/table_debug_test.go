package model

import (
	"fmt"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func TestGoldmarkTableParsing(t *testing.T) {
	markdown := `| Header 1 | Header 2 | Header 3 |
|----------|----------|----------|
| Cell 1   | Cell 2   | Cell 3   |
| Cell 4   | Cell 5   | Cell 6   |`

	source := []byte(markdown)

	// Create goldmark parser with table extension
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	// Parse markdown
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	// Debug: print AST structure
	fmt.Println("=== Table AST Structure ===")
	printTableAST(doc, source, 0)

	// Convert to our elements using ConvertMarkdownToElements
	elements := ConvertMarkdownToElements(markdown)

	fmt.Printf("\n=== Converted Table Elements ===\n")
	for i, elem := range elements {
		fmt.Printf("Element %d: %T\n", i, elem)
		printElementStructure(elem, 1)
	}
}

func printTableAST(node ast.Node, source []byte, depth int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	nodeType := fmt.Sprintf("%T", node)
	fmt.Printf("%s%s", indent, nodeType)

	// Print additional info for specific node types
	switch n := node.(type) {
	case *ast.Text:
		fmt.Printf(" (Text: %q)", string(n.Text(source)))
	}
	fmt.Println()

	// Recursively print children
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		printTableAST(child, source, depth+1)
	}
}

func printElementStructure(elem Element, depth int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	for i := 0; i < elem.GetElementCount(); i++ {
		child := elem.GetElement(i)
		fmt.Printf("%sChild %d: %T\n", indent, i, child)
		if child.GetElementCount() > 0 {
			printElementStructure(child, depth+1)
		}
	}
}
