package view

import (
	"strings"

	"github.com/desktopgame/ckro/internal/text"
	"github.com/desktopgame/ckro/internal/tui/model"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

// DocumentTextView renders DocumentElement
type DocumentTextView struct{}

func (d *DocumentTextView) Layout(textViewResolver TextViewResolver, e model.Element, width int) *TextLayout {
	children := []*TextLayout{}
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		children = append(children, childView.Layout(textViewResolver, childElement, width))
	}

	totalHeight := 0
	for _, child := range children {
		totalHeight += child.Height
	}

	return &TextLayout{
		Element:  e,
		Children: children,
		Width:    width,
		Height:   totalHeight,
	}
}

func (d *DocumentTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	y := 0
	for i, child := range textLayout.Children {
		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		childView.Draw(textViewResolver, child, renderer.Translate(0, y))
		y += child.Height
	}
}

func (d *DocumentTextView) Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int {
	return textLayout.Width
}

func (d *DocumentTextView) Height(textViewResolver TextViewResolver, textLayout *TextLayout) int {
	return textLayout.Height
}

// ParagraphTextViewGM renders ParagraphElement
type ParagraphTextViewGM struct{}

func (p *ParagraphTextViewGM) Layout(textViewResolver TextViewResolver, e model.Element, width int) *TextLayout {
	children := []*TextLayout{}
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		children = append(children, childView.Layout(textViewResolver, childElement, width))
	}

	// Calculate height based on text wrapping
	height := 1
	if len(children) > 0 {
		// For now, assume single line for inline elements
		height = 1
	}

	return &TextLayout{
		Element:  e,
		Children: children,
		Width:    width,
		Height:   height,
	}
}

func (p *ParagraphTextViewGM) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	x := 0
	for i, child := range textLayout.Children {
		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		childView.Draw(textViewResolver, child, renderer.Translate(x, 0))
		x += childView.Width(textViewResolver, child, 0)
	}
}

func (p *ParagraphTextViewGM) Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int {
	totalWidth := 0
	for i, child := range textLayout.Children {
		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		totalWidth += childView.Width(textViewResolver, child, 0)
	}
	return totalWidth
}

func (p *ParagraphTextViewGM) Height(textViewResolver TextViewResolver, textLayout *TextLayout) int {
	return textLayout.Height
}

// HeadingTextView renders HeadingElement
type HeadingTextView struct{}

func (h *HeadingTextView) Layout(textViewResolver TextViewResolver, e model.Element, width int) *TextLayout {
	children := []*TextLayout{}
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		children = append(children, childView.Layout(textViewResolver, childElement, width))
	}

	return &TextLayout{
		Element:  e,
		Children: children,
		Width:    width,
		Height:   1,
	}
}

func (h *HeadingTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	headingElement := textLayout.Element.(*model.HeadingElement)

	// Draw heading prefix (# ## ### etc.)
	prefix := strings.Repeat("#", headingElement.Level) + " "

	// Set color based on heading level
	var color tcell.Color
	switch headingElement.Level {
	case 1:
		color = tcell.ColorRed // H1: Red
	case 2:
		color = tcell.ColorBlue // H2: Blue
	case 3:
		color = tcell.ColorGreen // H3: Green
	case 4:
		color = tcell.ColorYellow // H4: Yellow
	case 5:
		color = tcell.ColorPurple // H5: Purple
	case 6:
		color = tcell.ColorTeal // H6: Teal
	default:
		color = tcell.ColorWhite // Default: White
	}

	style := tcell.StyleDefault.Bold(true).Foreground(color)

	x := 0
	for _, r := range prefix {
		renderer.SetContent(x, 0, r, nil, style)
		x++
	}

	// Draw heading content with the same style
	headingRenderer := &StyleRenderer{
		base:  renderer.Translate(x, 0),
		style: style,
	}

	for i, child := range textLayout.Children {
		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		childView.Draw(textViewResolver, child, headingRenderer)
		x += childView.Width(textViewResolver, child, 0)
		headingRenderer = &StyleRenderer{
			base:  renderer.Translate(x, 0),
			style: style,
		}
	}
}

func (h *HeadingTextView) Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int {
	headingElement := textLayout.Element.(*model.HeadingElement)
	prefixWidth := headingElement.Level + 1 // "# " or "## " etc.

	contentWidth := 0
	for i, child := range textLayout.Children {
		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		contentWidth += childView.Width(textViewResolver, child, 0)
	}

	return prefixWidth + contentWidth
}

func (h *HeadingTextView) Height(textViewResolver TextViewResolver, textLayout *TextLayout) int {
	return 1
}

// CodeBlockTextView renders CodeBlockElement
type CodeBlockTextView struct{}

func (c *CodeBlockTextView) Layout(textViewResolver TextViewResolver, e model.Element, width int) *TextLayout {
	codeBlock := e.(*model.CodeBlockElement)
	lines := strings.Split(codeBlock.Text, "\n")

	return &TextLayout{
		Element:  e,
		Children: nil,
		Width:    width,
		Height:   len(lines),
	}
}

func (c *CodeBlockTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	codeBlock := textLayout.Element.(*model.CodeBlockElement)
	lines := strings.Split(codeBlock.Text, "\n")

	style := tcell.StyleDefault.Background(tcell.ColorDarkGray)

	for y, line := range lines {
		x := 0
		for _, r := range line {
			renderer.SetContent(x, y, r, nil, style)
			x++
		}
	}
}

func (c *CodeBlockTextView) Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int {
	codeBlock := textLayout.Element.(*model.CodeBlockElement)
	lines := strings.Split(codeBlock.Text, "\n")

	if row >= 0 && row < len(lines) {
		return text.DisplayWidth(lines[row])
	}
	return 0
}

func (c *CodeBlockTextView) Height(textViewResolver TextViewResolver, textLayout *TextLayout) int {
	return textLayout.Height
}

// BlockquoteTextView renders BlockquoteElement
type BlockquoteTextView struct{}

func (b *BlockquoteTextView) Layout(textViewResolver TextViewResolver, e model.Element, width int) *TextLayout {
	children := []*TextLayout{}
	childWidth := width - 2 // Account for "> " prefix

	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		children = append(children, childView.Layout(textViewResolver, childElement, childWidth))
	}

	totalHeight := 0
	for _, child := range children {
		totalHeight += child.Height
	}

	return &TextLayout{
		Element:  e,
		Children: children,
		Width:    width,
		Height:   totalHeight,
	}
}

func (b *BlockquoteTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	style := tcell.StyleDefault.Foreground(tcell.ColorGray)

	y := 0
	for i, child := range textLayout.Children {
		// Draw "> " prefix
		renderer.SetContent(0, y, '>', nil, style)
		renderer.SetContent(1, y, ' ', nil, style)

		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		childView.Draw(textViewResolver, child, renderer.Translate(2, y))
		y += child.Height
	}
}

func (b *BlockquoteTextView) Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int {
	return textLayout.Width
}

func (b *BlockquoteTextView) Height(textViewResolver TextViewResolver, textLayout *TextLayout) int {
	return textLayout.Height
}

// ListTextView renders ListElement
type ListTextView struct{}

func (l *ListTextView) Layout(textViewResolver TextViewResolver, e model.Element, width int) *TextLayout {
	children := []*TextLayout{}
	childWidth := width - 4 // Account for list item prefix

	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		childTextLayout := childView.Layout(textViewResolver, childElement, childWidth)
		childTextLayout.Indent++
		children = append(children, childTextLayout)
	}

	totalHeight := 0
	for _, child := range children {
		totalHeight += child.Height
	}

	return &TextLayout{
		Element:  e,
		Children: children,
		Width:    width,
		Height:   totalHeight,
	}
}

func (l *ListTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	listElement := textLayout.Element.(*model.ListElement)

	y := 0
	for i, child := range textLayout.Children {
		// Draw list item prefix
		var prefix string
		if listElement.Ordered {
			prefix = strings.Repeat(" ", 2) + string(rune('1'+i)) + ". "
		} else {
			prefix = "  - "
		}

		if textLayout.Indent > 0 {
			prefix = strings.Repeat("  ", textLayout.Indent) + prefix
		}

		x := 0
		for _, r := range prefix {
			renderer.SetContent(x, y, r, nil, tcell.StyleDefault)
			x++
		}

		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		childView.Draw(textViewResolver, child, renderer.Translate(x, y))
		y += child.Height
	}
}

func (l *ListTextView) Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int {
	return textLayout.Width
}

func (l *ListTextView) Height(textViewResolver TextViewResolver, textLayout *TextLayout) int {
	return textLayout.Height
}

// ListItemTextView renders ListItemElement
type ListItemTextView struct{}

func (l *ListItemTextView) Layout(textViewResolver TextViewResolver, e model.Element, width int) *TextLayout {
	children := []*TextLayout{}

	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		children = append(children, childView.Layout(textViewResolver, childElement, width))
	}

	totalHeight := 0
	for _, child := range children {
		totalHeight += child.Height
	}

	return &TextLayout{
		Element:  e,
		Children: children,
		Width:    width,
		Height:   totalHeight,
	}
}

func (l *ListItemTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	y := 0
	for i, child := range textLayout.Children {
		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		childView.Draw(textViewResolver, child, renderer.Translate(0, y))
		y += child.Height
	}
}

func (l *ListItemTextView) Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int {
	return textLayout.Width
}

func (l *ListItemTextView) Height(textViewResolver TextViewResolver, textLayout *TextLayout) int {
	return textLayout.Height
}

// ThematicBreakTextView renders ThematicBreakElement
type ThematicBreakTextView struct{}

func (t *ThematicBreakTextView) Layout(textViewResolver TextViewResolver, e model.Element, width int) *TextLayout {
	return &TextLayout{
		Element:  e,
		Children: nil,
		Width:    width,
		Height:   1,
	}
}

func (t *ThematicBreakTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	// Draw horizontal line
	for x := 0; x < textLayout.Width; x++ {
		renderer.SetContent(x, 0, '-', nil, tcell.StyleDefault)
	}
}

func (t *ThematicBreakTextView) Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int {
	return textLayout.Width
}

func (t *ThematicBreakTextView) Height(textViewResolver TextViewResolver, textLayout *TextLayout) int {
	return 1
}

// TextElementView renders TextElement
type TextElementView struct{}

func (t *TextElementView) Layout(textViewResolver TextViewResolver, e model.Element, width int) *TextLayout {
	return &TextLayout{
		Element:  e,
		Children: nil,
		Width:    width,
		Height:   1,
	}
}

func (t *TextElementView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	textElement := textLayout.Element.(*model.TextElement)
	style := ConvertStyle(textElement.GetStyle())

	clusters := text.GraphemeClusters(textElement.Text)
	x := 0
	for _, cluster := range clusters {
		if cluster == "\t" {
			// Handle tab
		} else {
			runes := []rune(cluster)
			if len(runes) > 0 {
				mainRune := runes[0]
				var combining []rune
				if len(runes) > 1 {
					combining = runes[1:]
				}
				width := runewidth.RuneWidth(mainRune)

				renderer.SetContent(x, 0, mainRune, combining, style)
				if width == 2 {
					x++
					renderer.SetContent(x, 0, 0, nil, style)
				}
			}
			x++
		}
	}
}

func (t *TextElementView) Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int {
	textElement := textLayout.Element.(*model.TextElement)
	return text.DisplayWidth(textElement.Text)
}

func (t *TextElementView) Height(textViewResolver TextViewResolver, textLayout *TextLayout) int {
	return 1
}

// EmphasisTextView renders EmphasisElement
type EmphasisTextView struct{}

func (e *EmphasisTextView) Layout(textViewResolver TextViewResolver, el model.Element, width int) *TextLayout {
	children := []*TextLayout{}
	for i := 0; i < el.GetElementCount(); i++ {
		childElement := el.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		children = append(children, childView.Layout(textViewResolver, childElement, width))
	}

	return &TextLayout{
		Element:  el,
		Children: children,
		Width:    width,
		Height:   1,
	}
}

func (e *EmphasisTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	x := 0
	for i, child := range textLayout.Children {
		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)

		// Apply italic style to child renderer
		italicRenderer := &StyleRenderer{
			base:  renderer.Translate(x, 0),
			style: tcell.StyleDefault.Italic(true),
		}

		childView.Draw(textViewResolver, child, italicRenderer)
		x += childView.Width(textViewResolver, child, 0)
	}
}

func (e *EmphasisTextView) Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int {
	totalWidth := 0
	for i, child := range textLayout.Children {
		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		totalWidth += childView.Width(textViewResolver, child, 0)
	}
	return totalWidth
}

func (e *EmphasisTextView) Height(textViewResolver TextViewResolver, textLayout *TextLayout) int {
	return 1
}

// StrongTextView renders StrongElement
type StrongTextView struct{}

func (s *StrongTextView) Layout(textViewResolver TextViewResolver, e model.Element, width int) *TextLayout {
	children := []*TextLayout{}
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		children = append(children, childView.Layout(textViewResolver, childElement, width))
	}

	return &TextLayout{
		Element:  e,
		Children: children,
		Width:    width,
		Height:   1,
	}
}

func (s *StrongTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	x := 0
	for i, child := range textLayout.Children {
		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)

		// Apply bold style to child renderer
		boldRenderer := &StyleRenderer{
			base:  renderer.Translate(x, 0),
			style: tcell.StyleDefault.Bold(true),
		}

		childView.Draw(textViewResolver, child, boldRenderer)
		x += childView.Width(textViewResolver, child, 0)
	}
}

func (s *StrongTextView) Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int {
	totalWidth := 0
	for i, child := range textLayout.Children {
		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		totalWidth += childView.Width(textViewResolver, child, 0)
	}
	return totalWidth
}

func (s *StrongTextView) Height(textViewResolver TextViewResolver, textLayout *TextLayout) int {
	return 1
}

// CodeElementView renders CodeElement
type CodeElementView struct{}

func (c *CodeElementView) Layout(textViewResolver TextViewResolver, e model.Element, width int) *TextLayout {
	return &TextLayout{
		Element:  e,
		Children: nil,
		Width:    width,
		Height:   1,
	}
}

func (c *CodeElementView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	codeElement := textLayout.Element.(*model.CodeElement)
	style := tcell.StyleDefault.Background(tcell.ColorDarkGray)

	clusters := text.GraphemeClusters(codeElement.Text)
	x := 0
	for _, cluster := range clusters {
		runes := []rune(cluster)
		if len(runes) > 0 {
			mainRune := runes[0]
			var combining []rune
			if len(runes) > 1 {
				combining = runes[1:]
			}
			width := runewidth.RuneWidth(mainRune)

			renderer.SetContent(x, 0, mainRune, combining, style)
			if width == 2 {
				x++
				renderer.SetContent(x, 0, 0, nil, style)
			}
		}
		x++
	}
}

func (c *CodeElementView) Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int {
	codeElement := textLayout.Element.(*model.CodeElement)
	return text.DisplayWidth(codeElement.Text)
}

func (c *CodeElementView) Height(textViewResolver TextViewResolver, textLayout *TextLayout) int {
	return 1
}

// LinkTextView renders LinkElement
type LinkTextView struct{}

func (l *LinkTextView) Layout(textViewResolver TextViewResolver, e model.Element, width int) *TextLayout {
	children := []*TextLayout{}
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		children = append(children, childView.Layout(textViewResolver, childElement, width))
	}

	return &TextLayout{
		Element:  e,
		Children: children,
		Width:    width,
		Height:   1,
	}
}

func (l *LinkTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	x := 0
	for i, child := range textLayout.Children {
		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)

		// Apply underline style to child renderer
		linkRenderer := &StyleRenderer{
			base:  renderer.Translate(x, 0),
			style: tcell.StyleDefault.Underline(true).Foreground(tcell.ColorBlue),
		}

		childView.Draw(textViewResolver, child, linkRenderer)
		x += childView.Width(textViewResolver, child, 0)
	}
}

func (l *LinkTextView) Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int {
	totalWidth := 0
	for i, child := range textLayout.Children {
		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		totalWidth += childView.Width(textViewResolver, child, 0)
	}
	return totalWidth
}

func (l *LinkTextView) Height(textViewResolver TextViewResolver, textLayout *TextLayout) int {
	return 1
}

// ImageTextView renders ImageElement
type ImageTextView struct{}

func (i *ImageTextView) Layout(textViewResolver TextViewResolver, e model.Element, width int) *TextLayout {
	return &TextLayout{
		Element:  e,
		Children: nil,
		Width:    width,
		Height:   1,
	}
}

func (i *ImageTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	imageElement := textLayout.Element.(*model.ImageElement)

	// Display as [Image: alt text]
	displayText := "[Image: " + imageElement.Alt + "]"
	style := tcell.StyleDefault.Foreground(tcell.ColorGreen)

	x := 0
	for _, r := range displayText {
		renderer.SetContent(x, 0, r, nil, style)
		x++
	}
}

func (i *ImageTextView) Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int {
	imageElement := textLayout.Element.(*model.ImageElement)
	displayText := "[Image: " + imageElement.Alt + "]"
	return text.DisplayWidth(displayText)
}

func (i *ImageTextView) Height(textViewResolver TextViewResolver, textLayout *TextLayout) int {
	return 1
}

// StyleRenderer wraps a renderer to apply additional styles
type StyleRenderer struct {
	base  Renderer
	style tcell.Style
}

func (s *StyleRenderer) SetContent(x int, y int, primary rune, combining []rune, style tcell.Style) {
	// Merge styles
	mergedStyle := s.style

	// Get style components
	fg, bg, attr := style.Decompose()

	if attr != tcell.AttrNone {
		mergedStyle = mergedStyle.Attributes(attr)
	}
	if fg != tcell.ColorDefault {
		mergedStyle = mergedStyle.Foreground(fg)
	}
	if bg != tcell.ColorDefault {
		mergedStyle = mergedStyle.Background(bg)
	}

	s.base.SetContent(x, y, primary, combining, mergedStyle)
}

func (s *StyleRenderer) Translate(offsetX int, offsetY int) Renderer {
	return &StyleRenderer{
		base:  s.base.Translate(offsetX, offsetY),
		style: s.style,
	}
}

// SoftBreakTextView renders SoftBreakElement
type SoftBreakTextView struct{}

func (s *SoftBreakTextView) Layout(textViewResolver TextViewResolver, e model.Element, width int) *TextLayout {
	return &TextLayout{
		Element:  e,
		Children: nil,
		Width:    width,
		Height:   1,
	}
}

func (s *SoftBreakTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	// Draw a space for soft break
	renderer.SetContent(0, 0, ' ', nil, tcell.StyleDefault)
}

func (s *SoftBreakTextView) Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int {
	return 1 // Single space
}

func (s *SoftBreakTextView) Height(textViewResolver TextViewResolver, textLayout *TextLayout) int {
	return 1
}

// HardBreakTextView renders HardBreakElement
type HardBreakTextView struct{}

func (h *HardBreakTextView) Layout(textViewResolver TextViewResolver, e model.Element, width int) *TextLayout {
	return &TextLayout{
		Element:  e,
		Children: nil,
		Width:    width,
		Height:   1,
	}
}

func (h *HardBreakTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	// Hard break creates a line break - no visual content needed
}

func (h *HardBreakTextView) Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int {
	return 0 // Line break has no width
}

func (h *HardBreakTextView) Height(textViewResolver TextViewResolver, textLayout *TextLayout) int {
	return 1
}

// TableTextView renders TableElement
type TableTextView struct{}

func (t *TableTextView) Layout(textViewResolver TextViewResolver, e model.Element, width int) *TextLayout {
	children := []*TextLayout{}

	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		children = append(children, childView.Layout(textViewResolver, childElement, width-2)) // Account for left/right borders
	}

	totalHeight := 0
	for i, child := range children {
		totalHeight += child.Height

		// Add space for header separator line
		childElement := e.GetElement(i)
		if _, isHeader := childElement.(*model.TableHeaderElement); isHeader {
			totalHeight += 1 // Add 1 line for header separator
		}
	}

	// Add space for top and bottom borders
	totalHeight += 2

	return &TextLayout{
		Element:  e,
		Children: children,
		Width:    width,
		Height:   totalHeight,
	}
}

func (t *TableTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	tableWidth := textLayout.Width
	tableHeight := textLayout.Height

	// Draw top border
	renderer.SetContent(0, 0, '┌', nil, tcell.StyleDefault)
	for x := 1; x < tableWidth-1; x++ {
		renderer.SetContent(x, 0, '─', nil, tcell.StyleDefault)
	}
	renderer.SetContent(tableWidth-1, 0, '┐', nil, tcell.StyleDefault)

	// Draw bottom border
	renderer.SetContent(0, tableHeight-1, '└', nil, tcell.StyleDefault)
	for x := 1; x < tableWidth-1; x++ {
		renderer.SetContent(x, tableHeight-1, '─', nil, tcell.StyleDefault)
	}
	renderer.SetContent(tableWidth-1, tableHeight-1, '┘', nil, tcell.StyleDefault)

	// Draw left and right borders
	for y := 1; y < tableHeight-1; y++ {
		renderer.SetContent(0, y, '│', nil, tcell.StyleDefault)
		renderer.SetContent(tableWidth-1, y, '│', nil, tcell.StyleDefault)
	}

	// Draw table content and header separator
	y := 1 // Start after top border
	for i, child := range textLayout.Children {
		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		childView.Draw(textViewResolver, child, renderer.Translate(1, y)) // Offset by left border
		y += child.Height

		// Draw horizontal separator after header
		if _, isHeader := childElement.(*model.TableHeaderElement); isHeader {
			// Draw header separator line
			renderer.SetContent(0, y, '├', nil, tcell.StyleDefault)
			for x := 1; x < tableWidth-1; x++ {
				renderer.SetContent(x, y, '─', nil, tcell.StyleDefault)
			}
			renderer.SetContent(tableWidth-1, y, '┤', nil, tcell.StyleDefault)
			y++ // Move to next line after separator
		}
	}
}

func (t *TableTextView) Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int {
	return textLayout.Width
}

func (t *TableTextView) Height(textViewResolver TextViewResolver, textLayout *TextLayout) int {
	return textLayout.Height
}

// TableHeaderTextView renders TableHeaderElement
type TableHeaderTextView struct{}

func (t *TableHeaderTextView) Layout(textViewResolver TextViewResolver, e model.Element, width int) *TextLayout {
	children := []*TextLayout{}

	// Calculate cell width - distribute available width among cells
	cellCount := e.GetElementCount()
	if cellCount == 0 {
		return &TextLayout{
			Element:  e,
			Children: children,
			Width:    width,
			Height:   1,
		}
	}

	// Account for cell separators (│) between cells
	availableWidth := width - (cellCount - 1)
	cellWidth := availableWidth / cellCount

	for i := 0; i < cellCount; i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		children = append(children, childView.Layout(textViewResolver, childElement, cellWidth))
	}

	return &TextLayout{
		Element:  e,
		Children: children,
		Width:    width,
		Height:   1,
	}
}

func (t *TableHeaderTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	style := tcell.StyleDefault.Bold(true)
	cellCount := len(textLayout.Children)
	if cellCount == 0 {
		return
	}

	availableWidth := textLayout.Width - (cellCount - 1)
	cellWidth := availableWidth / cellCount

	x := 0
	for i, child := range textLayout.Children {
		// Draw cell separator
		if i > 0 {
			renderer.SetContent(x, 0, '│', nil, style)
			x++
		}

		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)

		// Apply bold style to header cells
		headerRenderer := &StyleRenderer{
			base:  renderer.Translate(x, 0),
			style: style,
		}

		childView.Draw(textViewResolver, child, headerRenderer)
		x += cellWidth
	}
}

func (t *TableHeaderTextView) Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int {
	return textLayout.Width
}

func (t *TableHeaderTextView) Height(textViewResolver TextViewResolver, textLayout *TextLayout) int {
	return 1
}

// TableRowTextViewGM renders TableRowElement
type TableRowTextViewGM struct{}

func (t *TableRowTextViewGM) Layout(textViewResolver TextViewResolver, e model.Element, width int) *TextLayout {
	children := []*TextLayout{}

	// Calculate cell width - distribute available width among cells
	cellCount := e.GetElementCount()
	if cellCount == 0 {
		return &TextLayout{
			Element:  e,
			Children: children,
			Width:    width,
			Height:   1,
		}
	}

	// Account for cell separators (|) between cells
	availableWidth := width - (cellCount - 1)
	cellWidth := availableWidth / cellCount

	for i := 0; i < cellCount; i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		children = append(children, childView.Layout(textViewResolver, childElement, cellWidth))
	}

	return &TextLayout{
		Element:  e,
		Children: children,
		Width:    width,
		Height:   1,
	}
}

func (t *TableRowTextViewGM) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	cellCount := len(textLayout.Children)
	if cellCount == 0 {
		return
	}

	availableWidth := textLayout.Width - (cellCount - 1)
	cellWidth := availableWidth / cellCount

	x := 0
	for i, child := range textLayout.Children {
		// Draw cell separator
		if i > 0 {
			renderer.SetContent(x, 0, '│', nil, tcell.StyleDefault)
			x++
		}

		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		childView.Draw(textViewResolver, child, renderer.Translate(x, 0))
		x += cellWidth
	}
}

func (t *TableRowTextViewGM) Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int {
	return textLayout.Width
}

func (t *TableRowTextViewGM) Height(textViewResolver TextViewResolver, textLayout *TextLayout) int {
	return 1
}

// TableCellTextView renders TableCellElement
type TableCellTextView struct{}

func (t *TableCellTextView) Layout(textViewResolver TextViewResolver, e model.Element, width int) *TextLayout {
	children := []*TextLayout{}

	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		children = append(children, childView.Layout(textViewResolver, childElement, width))
	}

	return &TextLayout{
		Element:  e,
		Children: children,
		Width:    width,
		Height:   1,
	}
}

func (t *TableCellTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	x := 0
	for i, child := range textLayout.Children {
		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		childView.Draw(textViewResolver, child, renderer.Translate(x, 0))
		x += childView.Width(textViewResolver, child, 0)
	}
}

func (t *TableCellTextView) Width(textViewResolver TextViewResolver, textLayout *TextLayout, row int) int {
	totalWidth := 0
	for i, child := range textLayout.Children {
		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		totalWidth += childView.Width(textViewResolver, child, 0)
	}
	return totalWidth
}

func (t *TableCellTextView) Height(textViewResolver TextViewResolver, textLayout *TextLayout) int {
	return 1
}
