package model

// Goldmark AST nodes に対応するElement型群

// Document element
type DocumentElement struct {
	StartPosition Position
	EndPosition   Position
	Children      []Element
}

func (e *DocumentElement) GetStyle() *Style           { return nil }
func (e *DocumentElement) GetText() string            { return "" }
func (e *DocumentElement) GetStartPosition() Position { return e.StartPosition }
func (e *DocumentElement) GetEndPosition() Position   { return e.EndPosition }
func (e *DocumentElement) GetElement(index int) Element {
	if index >= 0 && index < len(e.Children) {
		return e.Children[index]
	}
	return nil
}
func (e *DocumentElement) GetElementCount() int { return len(e.Children) }

// Block elements

// Paragraph element
type ParagraphElement struct {
	StartPosition Position
	EndPosition   Position
	Children      []Element
}

func (e *ParagraphElement) GetStyle() *Style           { return nil }
func (e *ParagraphElement) GetText() string            { return "" }
func (e *ParagraphElement) GetStartPosition() Position { return e.StartPosition }
func (e *ParagraphElement) GetEndPosition() Position   { return e.EndPosition }
func (e *ParagraphElement) GetElement(index int) Element {
	if index >= 0 && index < len(e.Children) {
		return e.Children[index]
	}
	return nil
}
func (e *ParagraphElement) GetElementCount() int { return len(e.Children) }

// Heading element
type HeadingElement struct {
	StartPosition Position
	EndPosition   Position
	Level         int
	Children      []Element
}

func (e *HeadingElement) GetStyle() *Style           { return nil }
func (e *HeadingElement) GetText() string            { return "" }
func (e *HeadingElement) GetStartPosition() Position { return e.StartPosition }
func (e *HeadingElement) GetEndPosition() Position   { return e.EndPosition }
func (e *HeadingElement) GetElement(index int) Element {
	if index >= 0 && index < len(e.Children) {
		return e.Children[index]
	}
	return nil
}
func (e *HeadingElement) GetElementCount() int { return len(e.Children) }

// CodeBlock element
type CodeBlockElement struct {
	StartPosition Position
	EndPosition   Position
	Language      string
	Text          string
}

func (e *CodeBlockElement) GetStyle() *Style             { return nil }
func (e *CodeBlockElement) GetText() string              { return e.Text }
func (e *CodeBlockElement) GetStartPosition() Position   { return e.StartPosition }
func (e *CodeBlockElement) GetEndPosition() Position     { return e.EndPosition }
func (e *CodeBlockElement) GetElement(index int) Element { return nil }
func (e *CodeBlockElement) GetElementCount() int         { return 0 }

// Blockquote element
type BlockquoteElement struct {
	StartPosition Position
	EndPosition   Position
	Children      []Element
}

func (e *BlockquoteElement) GetStyle() *Style           { return nil }
func (e *BlockquoteElement) GetText() string            { return "" }
func (e *BlockquoteElement) GetStartPosition() Position { return e.StartPosition }
func (e *BlockquoteElement) GetEndPosition() Position   { return e.EndPosition }
func (e *BlockquoteElement) GetElement(index int) Element {
	if index >= 0 && index < len(e.Children) {
		return e.Children[index]
	}
	return nil
}
func (e *BlockquoteElement) GetElementCount() int { return len(e.Children) }

// List element
type ListElement struct {
	StartPosition Position
	EndPosition   Position
	Ordered       bool
	Tight         bool
	Start         int
	Children      []Element
}

func (e *ListElement) GetStyle() *Style           { return nil }
func (e *ListElement) GetText() string            { return "" }
func (e *ListElement) GetStartPosition() Position { return e.StartPosition }
func (e *ListElement) GetEndPosition() Position   { return e.EndPosition }
func (e *ListElement) GetElement(index int) Element {
	if index >= 0 && index < len(e.Children) {
		return e.Children[index]
	}
	return nil
}
func (e *ListElement) GetElementCount() int { return len(e.Children) }

// ListItem element
type ListItemElement struct {
	StartPosition Position
	EndPosition   Position
	Children      []Element
}

func (e *ListItemElement) GetStyle() *Style           { return nil }
func (e *ListItemElement) GetText() string            { return "" }
func (e *ListItemElement) GetStartPosition() Position { return e.StartPosition }
func (e *ListItemElement) GetEndPosition() Position   { return e.EndPosition }
func (e *ListItemElement) GetElement(index int) Element {
	if index >= 0 && index < len(e.Children) {
		return e.Children[index]
	}
	return nil
}
func (e *ListItemElement) GetElementCount() int { return len(e.Children) }

// ThematicBreak element
type ThematicBreakElement struct {
	StartPosition Position
	EndPosition   Position
}

func (e *ThematicBreakElement) GetStyle() *Style             { return nil }
func (e *ThematicBreakElement) GetText() string              { return "" }
func (e *ThematicBreakElement) GetStartPosition() Position   { return e.StartPosition }
func (e *ThematicBreakElement) GetEndPosition() Position     { return e.EndPosition }
func (e *ThematicBreakElement) GetElement(index int) Element { return nil }
func (e *ThematicBreakElement) GetElementCount() int         { return 0 }

// Table elements
type TableElement struct {
	StartPosition Position
	EndPosition   Position
	Children      []Element
}

func (e *TableElement) GetStyle() *Style           { return nil }
func (e *TableElement) GetText() string            { return "" }
func (e *TableElement) GetStartPosition() Position { return e.StartPosition }
func (e *TableElement) GetEndPosition() Position   { return e.EndPosition }
func (e *TableElement) GetElement(index int) Element {
	if index >= 0 && index < len(e.Children) {
		return e.Children[index]
	}
	return nil
}
func (e *TableElement) GetElementCount() int { return len(e.Children) }

type TableHeaderElement struct {
	StartPosition Position
	EndPosition   Position
	Children      []Element
}

func (e *TableHeaderElement) GetStyle() *Style           { return nil }
func (e *TableHeaderElement) GetText() string            { return "" }
func (e *TableHeaderElement) GetStartPosition() Position { return e.StartPosition }
func (e *TableHeaderElement) GetEndPosition() Position   { return e.EndPosition }
func (e *TableHeaderElement) GetElement(index int) Element {
	if index >= 0 && index < len(e.Children) {
		return e.Children[index]
	}
	return nil
}
func (e *TableHeaderElement) GetElementCount() int { return len(e.Children) }

type TableRowElement struct {
	StartPosition Position
	EndPosition   Position
	Children      []Element
}

func (e *TableRowElement) GetStyle() *Style           { return nil }
func (e *TableRowElement) GetText() string            { return "" }
func (e *TableRowElement) GetStartPosition() Position { return e.StartPosition }
func (e *TableRowElement) GetEndPosition() Position   { return e.EndPosition }
func (e *TableRowElement) GetElement(index int) Element {
	if index >= 0 && index < len(e.Children) {
		return e.Children[index]
	}
	return nil
}
func (e *TableRowElement) GetElementCount() int { return len(e.Children) }

type TableCellElement struct {
	StartPosition Position
	EndPosition   Position
	Alignment     string
	Children      []Element
}

func (e *TableCellElement) GetStyle() *Style           { return nil }
func (e *TableCellElement) GetText() string            { return "" }
func (e *TableCellElement) GetStartPosition() Position { return e.StartPosition }
func (e *TableCellElement) GetEndPosition() Position   { return e.EndPosition }
func (e *TableCellElement) GetElement(index int) Element {
	if index >= 0 && index < len(e.Children) {
		return e.Children[index]
	}
	return nil
}
func (e *TableCellElement) GetElementCount() int { return len(e.Children) }

// Inline elements

// Text element
type TextElement struct {
	StartPosition Position
	EndPosition   Position
	Text          string
}

func (e *TextElement) GetStyle() *Style             { return nil }
func (e *TextElement) GetText() string              { return e.Text }
func (e *TextElement) GetStartPosition() Position   { return e.StartPosition }
func (e *TextElement) GetEndPosition() Position     { return e.EndPosition }
func (e *TextElement) GetElement(index int) Element { return nil }
func (e *TextElement) GetElementCount() int         { return 0 }

// Emphasis element
type EmphasisElement struct {
	StartPosition Position
	EndPosition   Position
	Children      []Element
}

func (e *EmphasisElement) GetStyle() *Style           { return nil }
func (e *EmphasisElement) GetText() string            { return "" }
func (e *EmphasisElement) GetStartPosition() Position { return e.StartPosition }
func (e *EmphasisElement) GetEndPosition() Position   { return e.EndPosition }
func (e *EmphasisElement) GetElement(index int) Element {
	if index >= 0 && index < len(e.Children) {
		return e.Children[index]
	}
	return nil
}
func (e *EmphasisElement) GetElementCount() int { return len(e.Children) }

// Strong element
type StrongElement struct {
	StartPosition Position
	EndPosition   Position
	Children      []Element
}

func (e *StrongElement) GetStyle() *Style           { return nil }
func (e *StrongElement) GetText() string            { return "" }
func (e *StrongElement) GetStartPosition() Position { return e.StartPosition }
func (e *StrongElement) GetEndPosition() Position   { return e.EndPosition }
func (e *StrongElement) GetElement(index int) Element {
	if index >= 0 && index < len(e.Children) {
		return e.Children[index]
	}
	return nil
}
func (e *StrongElement) GetElementCount() int { return len(e.Children) }

// Code element
type CodeElement struct {
	StartPosition Position
	EndPosition   Position
	Text          string
}

func (e *CodeElement) GetStyle() *Style             { return nil }
func (e *CodeElement) GetText() string              { return e.Text }
func (e *CodeElement) GetStartPosition() Position   { return e.StartPosition }
func (e *CodeElement) GetEndPosition() Position     { return e.EndPosition }
func (e *CodeElement) GetElement(index int) Element { return nil }
func (e *CodeElement) GetElementCount() int         { return 0 }

// Link element
type LinkElement struct {
	StartPosition Position
	EndPosition   Position
	URL           string
	Title         string
	Children      []Element
}

func (e *LinkElement) GetStyle() *Style           { return nil }
func (e *LinkElement) GetText() string            { return "" }
func (e *LinkElement) GetStartPosition() Position { return e.StartPosition }
func (e *LinkElement) GetEndPosition() Position   { return e.EndPosition }
func (e *LinkElement) GetElement(index int) Element {
	if index >= 0 && index < len(e.Children) {
		return e.Children[index]
	}
	return nil
}
func (e *LinkElement) GetElementCount() int { return len(e.Children) }

// Image element
type ImageElement struct {
	StartPosition Position
	EndPosition   Position
	URL           string
	Title         string
	Alt           string
}

func (e *ImageElement) GetStyle() *Style             { return nil }
func (e *ImageElement) GetText() string              { return e.Alt }
func (e *ImageElement) GetStartPosition() Position   { return e.StartPosition }
func (e *ImageElement) GetEndPosition() Position     { return e.EndPosition }
func (e *ImageElement) GetElement(index int) Element { return nil }
func (e *ImageElement) GetElementCount() int         { return 0 }

// LineBreak elements
type SoftBreakElement struct {
	StartPosition Position
	EndPosition   Position
}

func (e *SoftBreakElement) GetStyle() *Style             { return nil }
func (e *SoftBreakElement) GetText() string              { return " " }
func (e *SoftBreakElement) GetStartPosition() Position   { return e.StartPosition }
func (e *SoftBreakElement) GetEndPosition() Position     { return e.EndPosition }
func (e *SoftBreakElement) GetElement(index int) Element { return nil }
func (e *SoftBreakElement) GetElementCount() int         { return 0 }

type HardBreakElement struct {
	StartPosition Position
	EndPosition   Position
}

func (e *HardBreakElement) GetStyle() *Style             { return nil }
func (e *HardBreakElement) GetText() string              { return "\n" }
func (e *HardBreakElement) GetStartPosition() Position   { return e.StartPosition }
func (e *HardBreakElement) GetEndPosition() Position     { return e.EndPosition }
func (e *HardBreakElement) GetElement(index int) Element { return nil }
func (e *HardBreakElement) GetElementCount() int         { return 0 }
