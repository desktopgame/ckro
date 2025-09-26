package model

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// GoldmarkConverter converts goldmark AST to Element structure
type GoldmarkConverter struct {
	source []byte
}

// NewGoldmarkConverter creates a new converter
func NewGoldmarkConverter(source []byte) *GoldmarkConverter {
	return &GoldmarkConverter{source: source}
}

// Convert converts markdown text to Element array using goldmark
func ConvertMarkdownToElements(markdownText string) []Element {
	source := []byte(markdownText)
	converter := NewGoldmarkConverter(source)

	// Create goldmark parser with extensions
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

	// Convert AST to Elements
	elements := []Element{}
	for child := doc.FirstChild(); child != nil; child = child.NextSibling() {
		if element := converter.convertNode(child); element != nil {
			elements = append(elements, element)
		}
	}

	return elements
}

// convertNode converts a goldmark AST node to Element
func (c *GoldmarkConverter) convertNode(node ast.Node) Element {
	switch n := node.(type) {
	case *ast.Document:
		return c.convertDocument(n)
	case *ast.Paragraph:
		return c.convertParagraph(n)
	case *ast.Heading:
		return c.convertHeading(n)
	case *ast.CodeBlock:
		return c.convertCodeBlock(n)
	case *ast.FencedCodeBlock:
		return c.convertFencedCodeBlock(n)
	case *ast.Blockquote:
		return c.convertBlockquote(n)
	case *ast.List:
		return c.convertList(n)
	case *ast.ListItem:
		return c.convertListItem(n)
	case *ast.ThematicBreak:
		return c.convertThematicBreak(n)
	case *ast.Text:
		return c.convertText(n)
	case *ast.Emphasis:
		if n.Level == 2 {
			return c.convertStrong(n)
		}
		return c.convertEmphasis(n)
	case *ast.CodeSpan:
		return c.convertCodeSpan(n)
	case *ast.Link:
		return c.convertLink(n)
	case *ast.Image:
		return c.convertImage(n)
	case *ast.AutoLink:
		return c.convertAutoLink(n)
	case *ast.RawHTML:
		return c.convertRawHTML(n)
	default:
		// Handle table elements from extension
		if c.isTableNode(node) {
			return c.convertTableNode(node)
		}
		// For unknown nodes, try to convert children
		return c.convertGenericNode(node)
	}
}

// Helper function to get position from AST node
func (c *GoldmarkConverter) getPosition(node ast.Node) (Position, Position) {
	// Default positions
	startPos := Position{Row: 0, Column: 0}
	endPos := Position{Row: 0, Column: 0}

	// Check if this is a block node by checking its kind
	kind := node.Kind()
	isBlockNode := kind == ast.KindDocument ||
		kind == ast.KindParagraph ||
		kind == ast.KindHeading ||
		kind == ast.KindCodeBlock ||
		kind == ast.KindFencedCodeBlock ||
		kind == ast.KindBlockquote ||
		kind == ast.KindList ||
		kind == ast.KindListItem ||
		kind == ast.KindThematicBreak

	if isBlockNode && node.Lines().Len() > 0 {
		segment := node.Lines().At(0)
		startPos = c.byteOffsetToPosition(segment.Start)
		endPos = c.byteOffsetToPosition(segment.Stop)
		return startPos, endPos
	}

	// For inline nodes or nodes without lines, try to get position from parent
	if parent := node.Parent(); parent != nil {
		parentKind := parent.Kind()
		isParentBlock := parentKind == ast.KindDocument ||
			parentKind == ast.KindParagraph ||
			parentKind == ast.KindHeading ||
			parentKind == ast.KindCodeBlock ||
			parentKind == ast.KindFencedCodeBlock ||
			parentKind == ast.KindBlockquote ||
			parentKind == ast.KindList ||
			parentKind == ast.KindListItem ||
			parentKind == ast.KindThematicBreak

		if isParentBlock && parent.Lines().Len() > 0 {
			segment := parent.Lines().At(0)
			startPos = c.byteOffsetToPosition(segment.Start)
			endPos = c.byteOffsetToPosition(segment.Stop)
		}
	}

	return startPos, endPos
}

// Convert byte offset to row/column position
func (c *GoldmarkConverter) byteOffsetToPosition(offset int) Position {
	if offset > len(c.source) {
		offset = len(c.source)
	}

	row := 0
	col := 0
	for i := 0; i < offset; i++ {
		if c.source[i] == '\n' {
			row++
			col = 0
		} else {
			col++
		}
	}
	return Position{Row: row, Column: col}
}

// Convert children nodes
func (c *GoldmarkConverter) convertChildren(node ast.Node) []Element {
	children := []Element{}
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if element := c.convertNode(child); element != nil {
			children = append(children, element)
		}
	}
	return children
}

// Document conversion
func (c *GoldmarkConverter) convertDocument(node *ast.Document) Element {
	startPos, endPos := c.getPosition(node)
	return &DocumentElement{
		StartPosition: startPos,
		EndPosition:   endPos,
		Children:      c.convertChildren(node),
	}
}

// Paragraph conversion
func (c *GoldmarkConverter) convertParagraph(node *ast.Paragraph) Element {
	startPos, endPos := c.getPosition(node)
	return &ParagraphElement{
		StartPosition: startPos,
		EndPosition:   endPos,
		Children:      c.convertChildren(node),
	}
}

// Heading conversion
func (c *GoldmarkConverter) convertHeading(node *ast.Heading) Element {
	startPos, endPos := c.getPosition(node)
	return &HeadingElement{
		StartPosition: startPos,
		EndPosition:   endPos,
		Level:         node.Level,
		Children:      c.convertChildren(node),
	}
}

// CodeBlock conversion
func (c *GoldmarkConverter) convertCodeBlock(node *ast.CodeBlock) Element {
	startPos, endPos := c.getPosition(node)
	var buf bytes.Buffer
	for i := 0; i < node.Lines().Len(); i++ {
		line := node.Lines().At(i)
		buf.Write(line.Value(c.source))
	}
	return &CodeBlockElement{
		StartPosition: startPos,
		EndPosition:   endPos,
		Language:      "",
		Text:          buf.String(),
	}
}

// FencedCodeBlock conversion
func (c *GoldmarkConverter) convertFencedCodeBlock(node *ast.FencedCodeBlock) Element {
	startPos, endPos := c.getPosition(node)
	var buf bytes.Buffer
	for i := 0; i < node.Lines().Len(); i++ {
		line := node.Lines().At(i)
		buf.Write(line.Value(c.source))
	}

	language := ""
	if node.Info != nil {
		language = string(node.Info.Text(c.source))
	}

	return &CodeBlockElement{
		StartPosition: startPos,
		EndPosition:   endPos,
		Language:      language,
		Text:          buf.String(),
	}
}

// Blockquote conversion
func (c *GoldmarkConverter) convertBlockquote(node *ast.Blockquote) Element {
	startPos, endPos := c.getPosition(node)
	return &BlockquoteElement{
		StartPosition: startPos,
		EndPosition:   endPos,
		Children:      c.convertChildren(node),
	}
}

// List conversion
func (c *GoldmarkConverter) convertList(node *ast.List) Element {
	startPos, endPos := c.getPosition(node)
	return &ListElement{
		StartPosition: startPos,
		EndPosition:   endPos,
		Ordered:       node.IsOrdered(),
		Tight:         node.IsTight,
		Start:         node.Start,
		Children:      c.convertChildren(node),
	}
}

// ListItem conversion
func (c *GoldmarkConverter) convertListItem(node *ast.ListItem) Element {
	startPos, endPos := c.getPosition(node)
	return &ListItemElement{
		StartPosition: startPos,
		EndPosition:   endPos,
		Children:      c.convertChildren(node),
	}
}

// ThematicBreak conversion
func (c *GoldmarkConverter) convertThematicBreak(node *ast.ThematicBreak) Element {
	startPos, endPos := c.getPosition(node)
	return &ThematicBreakElement{
		StartPosition: startPos,
		EndPosition:   endPos,
	}
}

// Text conversion
func (c *GoldmarkConverter) convertText(node *ast.Text) Element {
	startPos, endPos := c.getPosition(node)
	text := string(node.Text(c.source))
	return &TextElement{
		StartPosition: startPos,
		EndPosition:   endPos,
		Text:          text,
	}
}

// Emphasis conversion
func (c *GoldmarkConverter) convertEmphasis(node *ast.Emphasis) Element {
	startPos, endPos := c.getPosition(node)
	return &EmphasisElement{
		StartPosition: startPos,
		EndPosition:   endPos,
		Children:      c.convertChildren(node),
	}
}

// Strong conversion
func (c *GoldmarkConverter) convertStrong(node *ast.Emphasis) Element {
	startPos, endPos := c.getPosition(node)
	return &StrongElement{
		StartPosition: startPos,
		EndPosition:   endPos,
		Children:      c.convertChildren(node),
	}
}

// CodeSpan conversion
func (c *GoldmarkConverter) convertCodeSpan(node *ast.CodeSpan) Element {
	startPos, endPos := c.getPosition(node)
	text := string(node.Text(c.source))
	return &CodeElement{
		StartPosition: startPos,
		EndPosition:   endPos,
		Text:          text,
	}
}

// Link conversion
func (c *GoldmarkConverter) convertLink(node *ast.Link) Element {
	startPos, endPos := c.getPosition(node)
	url := string(node.Destination)
	title := ""
	if node.Title != nil {
		title = string(node.Title)
	}
	return &LinkElement{
		StartPosition: startPos,
		EndPosition:   endPos,
		URL:           url,
		Title:         title,
		Children:      c.convertChildren(node),
	}
}

// Image conversion
func (c *GoldmarkConverter) convertImage(node *ast.Image) Element {
	startPos, endPos := c.getPosition(node)
	url := string(node.Destination)
	title := ""
	if node.Title != nil {
		title = string(node.Title)
	}

	// Get alt text from children
	alt := ""
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if textNode, ok := child.(*ast.Text); ok {
			alt += string(textNode.Text(c.source))
		}
	}

	return &ImageElement{
		StartPosition: startPos,
		EndPosition:   endPos,
		URL:           url,
		Title:         title,
		Alt:           alt,
	}
}

// AutoLink conversion
func (c *GoldmarkConverter) convertAutoLink(node *ast.AutoLink) Element {
	startPos, endPos := c.getPosition(node)
	url := string(node.URL(c.source))
	return &LinkElement{
		StartPosition: startPos,
		EndPosition:   endPos,
		URL:           url,
		Title:         "",
		Children: []Element{
			&TextElement{
				StartPosition: startPos,
				EndPosition:   endPos,
				Text:          url,
			},
		},
	}
}

// RawHTML conversion (treat as text for now)
func (c *GoldmarkConverter) convertRawHTML(node *ast.RawHTML) Element {
	startPos, endPos := c.getPosition(node)
	var buf bytes.Buffer
	for i := 0; i < node.Segments.Len(); i++ {
		segment := node.Segments.At(i)
		buf.Write(segment.Value(c.source))
	}
	return &TextElement{
		StartPosition: startPos,
		EndPosition:   endPos,
		Text:          buf.String(),
	}
}

// Check if node is a table-related node
func (c *GoldmarkConverter) isTableNode(node ast.Node) bool {
	kind := node.Kind()
	return strings.Contains(kind.String(), "Table")
}

// Convert table nodes (improved implementation)
func (c *GoldmarkConverter) convertTableNode(node ast.Node) Element {
	startPos, endPos := c.getPosition(node)
	kind := node.Kind().String()

	switch {
	case kind == "Table":
		// Table root
		return &TableElement{
			StartPosition: startPos,
			EndPosition:   endPos,
			Children:      c.convertChildren(node),
		}
	case kind == "TableHeader":
		// Table header - this is the key fix!
		return &TableHeaderElement{
			StartPosition: startPos,
			EndPosition:   endPos,
			Children:      c.convertChildren(node),
		}
	case kind == "TableRow":
		return &TableRowElementGM{
			StartPosition: startPos,
			EndPosition:   endPos,
			Children:      c.convertChildren(node),
		}
	case kind == "TableCell":
		return &TableCellElement{
			StartPosition: startPos,
			EndPosition:   endPos,
			Alignment:     "left", // Default alignment
			Children:      c.convertChildren(node),
		}
	default:
		return c.convertGenericNode(node)
	}
}

// Generic node conversion for unknown types
func (c *GoldmarkConverter) convertGenericNode(node ast.Node) Element {
	startPos, endPos := c.getPosition(node)

	// If it's a leaf node with text, convert to TextElement
	if node.FirstChild() == nil {
		if hasText := node.HasChildren(); !hasText {
			// Check if this is a block node that can have Lines
			kind := node.Kind()
			isBlockNode := kind == ast.KindDocument ||
				kind == ast.KindParagraph ||
				kind == ast.KindHeading ||
				kind == ast.KindCodeBlock ||
				kind == ast.KindFencedCodeBlock ||
				kind == ast.KindBlockquote ||
				kind == ast.KindList ||
				kind == ast.KindListItem ||
				kind == ast.KindThematicBreak

			// Try to extract text content only from block nodes
			if isBlockNode {
				var buf bytes.Buffer
				if node.Lines().Len() > 0 {
					for i := 0; i < node.Lines().Len(); i++ {
						line := node.Lines().At(i)
						buf.Write(line.Value(c.source))
					}
					return &TextElement{
						StartPosition: startPos,
						EndPosition:   endPos,
						Text:          buf.String(),
					}
				}
			}

			// For inline nodes without text, return empty text element
			return &TextElement{
				StartPosition: startPos,
				EndPosition:   endPos,
				Text:          "",
			}
		}
	}

	// For container nodes, create a generic paragraph-like element
	return &ParagraphElement{
		StartPosition: startPos,
		EndPosition:   endPos,
		Children:      c.convertChildren(node),
	}
}
