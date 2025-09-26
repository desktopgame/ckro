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

func (d *DocumentTextView) Layout(textViewResolver TextViewResolver, e model.Element, x, y, w, h int) *TextLayout {
	children := []*TextLayout{}
	totalHeight := 0
	maxWidth := 0
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)

		mw, mh := childView.MinimumSize(textViewResolver, childElement, w, 9999)
		childTextLayout := childView.Layout(textViewResolver, childElement, x, y, mw, mh)
		children = append(children, childTextLayout)
		y += mh

		totalHeight += childTextLayout.Height
		if mw > maxWidth {
			maxWidth = childTextLayout.Width
		}
	}

	return &TextLayout{
		Element:   e,
		Children:  children,
		RelativeX: x,
		RelativeY: y,
		Width:     maxWidth,
		Height:    totalHeight,
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

func (d *DocumentTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) (Width int, Height int) {
	totalHeight := 0
	maxWidth := 0
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)

		mw, mh := childView.MinimumSize(textViewResolver, childElement, width, 9999)

		totalHeight += mh
		if mw > maxWidth {
			maxWidth = mw
		}
	}
	return max(maxWidth, width), totalHeight
}

// ParagraphTextView renders ParagraphElement
type ParagraphTextView struct{}

func (p *ParagraphTextView) Layout(textViewResolver TextViewResolver, e model.Element, x, y, w, h int) *TextLayout {
	children := []*TextLayout{}
	totalWidth := 0
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)

		mw, _ := childView.MinimumSize(textViewResolver, childElement, w, 1)
		childTextLayout := childView.Layout(textViewResolver, childElement, x, y, mw, 1)
		children = append(children, childTextLayout)
		x += childTextLayout.Width

		totalWidth += childTextLayout.Width
	}
	return &TextLayout{
		Element:   e,
		Children:  children,
		RelativeX: x,
		RelativeY: y,
		Width:     totalWidth,
		Height:    1,
	}
}

func (p *ParagraphTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	for i, child := range textLayout.Children {
		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		childView.Draw(textViewResolver, child, renderer.Translate(child.RelativeX, child.RelativeY))
	}
}

func (p *ParagraphTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) (Width int, Height int) {
	totalWidth := 0
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)

		mw, _ := childView.MinimumSize(textViewResolver, childElement, width, 1)
		totalWidth += mw
	}
	return totalWidth, 1
}

// HeadingTextView renders HeadingElement
type HeadingTextView struct{}

func (ht *HeadingTextView) Layout(textViewResolver TextViewResolver, e model.Element, x, y, w, h int) *TextLayout {
	children := []*TextLayout{}
	headingElement := e.(*model.HeadingElement)
	prefixWidth := headingElement.Level + 1 // "# " or "## " etc.
	contentWidth := 0
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)

		mw, mh := childView.MinimumSize(textViewResolver, childElement, w, 1)
		childTextLayout := childView.Layout(textViewResolver, childElement, x+prefixWidth, y, mw, mh)
		children = append(children, childTextLayout)
		contentWidth += childTextLayout.Width
	}

	return &TextLayout{
		Element:   e,
		Children:  children,
		Width:     prefixWidth + contentWidth,
		Height:    1,
		RelativeX: x,
		RelativeY: y,
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
		base:  renderer.Translate(0, 0),
		style: style,
	}

	for i, child := range textLayout.Children {
		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		childView.Draw(textViewResolver, child, headingRenderer.Translate(child.RelativeX, child.RelativeY))
	}
}

func (h *HeadingTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) (Width int, Height int) {
	totalWidth := 0
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)

		mw, _ := childView.MinimumSize(textViewResolver, childElement, width, 1)

		totalWidth += mw
	}
	return totalWidth, 1
}

// CodeBlockTextView renders CodeBlockElement
type CodeBlockTextView struct{}

func (c *CodeBlockTextView) Layout(textViewResolver TextViewResolver, e model.Element, x, y, w, h int) *TextLayout {
	mw, mh := c.MinimumSize(textViewResolver, e, w, h)

	return &TextLayout{
		Element:   e,
		Children:  nil,
		Width:     mw,
		Height:    mh,
		RelativeX: x,
		RelativeY: y,
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

func (c *CodeBlockTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) (Width int, Height int) {
	codeBlock := e.(*model.CodeBlockElement)
	lines := strings.Split(codeBlock.Text, "\n")

	maxWidth := -1
	for _, line := range lines {
		w := text.DisplayWidth(line)
		if w > maxWidth {
			maxWidth = w
		}
	}
	return maxWidth, len(lines)
}

// BlockquoteTextView renders BlockquoteElement
type BlockquoteTextView struct{}

func (b *BlockquoteTextView) Layout(textViewResolver TextViewResolver, e model.Element, x, y, w, h int) *TextLayout {
	children := []*TextLayout{}
	maxWidth := -1
	totalHeight := 0
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		mw, _ := childView.MinimumSize(textViewResolver, childElement, w-2, 1)
		childTextLayout := childView.Layout(textViewResolver, childElement, x+2, y, mw, 1)
		children = append(children, childTextLayout)

		if mw > maxWidth {
			maxWidth = mw
		}
		totalHeight++
	}
	return &TextLayout{
		Element:  e,
		Children: children,
		Width:    maxWidth + 2,
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

func (b *BlockquoteTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) (Width int, Height int) {
	maxWidth := -1
	totalHeight := 0
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		mw, _ := childView.MinimumSize(textViewResolver, childElement, width, 1)

		if mw > maxWidth {
			maxWidth = mw
		}
		totalHeight++
	}
	return maxWidth, totalHeight
}

// ListTextView renders ListElement
type ListTextView struct{}

func (l *ListTextView) Layout(textViewResolver TextViewResolver, e model.Element, x, y, w, h int) *TextLayout {
	children := []*TextLayout{}
	maxWidth := -1
	totalHeight := 0
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		mw, mh := childView.MinimumSize(textViewResolver, childElement, w-2, 9999)
		childTextLayout := childView.Layout(textViewResolver, childElement, x+2, y, mw, mh)
		children = append(children, childTextLayout)

		if mw > maxWidth {
			maxWidth = mw
		}
		totalHeight += mh
	}

	return &TextLayout{
		Element:  e,
		Children: children,
		Width:    maxWidth,
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

func (l *ListTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) (Width int, Height int) {
	maxWidth := -1
	totalHeight := 0
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		mw, mh := childView.MinimumSize(textViewResolver, childElement, width-2, 9999)

		if mw > maxWidth {
			maxWidth = mw
		}
		totalHeight += mh
	}
	return maxWidth, totalHeight
}

// ListItemTextView renders ListItemElement
type ListItemTextView struct{}

func (l *ListItemTextView) Layout(textViewResolver TextViewResolver, e model.Element, x, y, w, h int) *TextLayout {
	children := []*TextLayout{}
	totalWidth := 0
	maxHeight := -1
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		mw, mh := childView.MinimumSize(textViewResolver, childElement, w, h)
		childTextLayout := childView.Layout(textViewResolver, childElement, x, y, mw, mh)
		children = append(children, childTextLayout)

		if mh > maxHeight {
			maxHeight = mh
		}
		totalWidth += mw
	}
	return &TextLayout{
		Element:  e,
		Children: children,
		Width:    totalWidth,
		Height:   maxHeight,
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

func (l *ListItemTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) (Width int, Height int) {
	totalWidth := 0
	maxHeight := -1
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		mw, mh := childView.MinimumSize(textViewResolver, childElement, width, height)

		if mh > maxHeight {
			maxHeight = mh
		}
		totalWidth += mw
	}
	return totalWidth, maxHeight
}

// ThematicBreakTextView renders ThematicBreakElement
type ThematicBreakTextView struct{}

func (t *ThematicBreakTextView) Layout(textViewResolver TextViewResolver, e model.Element, x, y, w, h int) *TextLayout {
	return &TextLayout{
		Element:  e,
		Children: nil,
		Width:    1,
		Height:   1,
	}
}

func (t *ThematicBreakTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	// Draw horizontal line
	for x := 0; x < textLayout.Width; x++ {
		renderer.SetContent(x, 0, '-', nil, tcell.StyleDefault)
	}
}

func (t *ThematicBreakTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) (Width int, Height int) {
	return 1, 1
}

// TextElementView renders TextElement
type TextElementView struct{}

func (t *TextElementView) Layout(textViewResolver TextViewResolver, e model.Element, x, y, w, h int) *TextLayout {
	return &TextLayout{
		Element:  e,
		Children: nil,
		Width:    text.DisplayWidth(e.GetText()),
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
			spaces := text.TabWidth - (x % text.TabWidth)
			for i := 0; i < spaces; i++ {
				renderer.SetContent(x+i, 0, ' ', nil, tcell.StyleDefault)
			}
			x += spaces
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

func (t *TextElementView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) (Width int, Height int) {
	return text.DisplayWidth(e.GetText()), 1
}

// EmphasisTextView renders EmphasisElement
type EmphasisTextView struct{}

func (em *EmphasisTextView) Layout(textViewResolver TextViewResolver, e model.Element, x, y, w, h int) *TextLayout {
	children := []*TextLayout{}
	totalWidth := 0
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		mw, _ := childView.MinimumSize(textViewResolver, childElement, w, 1)
		childTextLayout := childView.Layout(textViewResolver, childElement, x, y, mw, 1)
		children = append(children, childTextLayout)
		totalWidth += mw
	}

	return &TextLayout{
		Element:  e,
		Children: children,
		Width:    totalWidth,
		Height:   1,
	}
}

func (em *EmphasisTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
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
		x += child.Width
	}
}

func (em *EmphasisTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) (Width int, Height int) {
	totalWidth := 0
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		mw, _ := childView.MinimumSize(textViewResolver, childElement, width, 1)
		totalWidth += mw
	}
	return totalWidth, 1
}

// StrongTextView renders StrongElement
type StrongTextView struct{}

func (s *StrongTextView) Layout(textViewResolver TextViewResolver, e model.Element, x, y, w, h int) *TextLayout {
	children := []*TextLayout{}
	totalWidth := 0
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		mw, _ := childView.MinimumSize(textViewResolver, childElement, w, 1)
		childTextLayout := childView.Layout(textViewResolver, childElement, x, y, mw, 1)
		children = append(children, childTextLayout)

		totalWidth += mw
	}

	return &TextLayout{
		Element:  e,
		Children: children,
		Width:    totalWidth,
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
		x += child.Width
	}
}

func (s *StrongTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) (Width int, Height int) {
	totalWidth := 0
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		mw, _ := childView.MinimumSize(textViewResolver, childElement, width, 1)
		totalWidth += mw
	}
	return totalWidth, 1
}

// CodeElementView renders CodeElement
type CodeElementView struct{}

func (c *CodeElementView) Layout(textViewResolver TextViewResolver, e model.Element, x, y, w, h int) *TextLayout {
	return &TextLayout{
		Element:  e,
		Children: nil,
		Width:    text.DisplayWidth(e.GetText()),
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

func (c *CodeElementView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) (Width int, Height int) {
	return text.DisplayWidth(e.GetText()), 1
}

// LinkTextView renders LinkElement
type LinkTextView struct{}

func (l *LinkTextView) Layout(textViewResolver TextViewResolver, e model.Element, x, y, w, h int) *TextLayout {
	children := []*TextLayout{}
	totalWidth := 0
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		mw, _ := childView.MinimumSize(textViewResolver, childElement, w, 1)
		childTextLayout := childView.Layout(textViewResolver, childElement, x, y, mw, 1)
		children = append(children, childTextLayout)
		totalWidth += mw
	}

	return &TextLayout{
		Element:  e,
		Children: children,
		Width:    totalWidth,
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
		x += child.Width
	}
}

func (l *LinkTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) (Width int, Height int) {
	totalWidth := 0
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		mw, _ := childView.MinimumSize(textViewResolver, childElement, width, 1)
		totalWidth += mw
	}
	return totalWidth, 1
}

// ImageTextView renders ImageElement
type ImageTextView struct{}

func (i *ImageTextView) Layout(textViewResolver TextViewResolver, e model.Element, x, y, w, h int) *TextLayout {
	mw, mh := i.MinimumSize(textViewResolver, e, w, h)
	return &TextLayout{
		Element:  e,
		Children: nil,
		Width:    mw,
		Height:   mh,
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

func (i *ImageTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) (Width int, Height int) {
	imageElement := e.(*model.ImageElement)
	displayText := "[Image: " + imageElement.Alt + "]"
	return text.DisplayWidth(displayText), 1
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

func (s *SoftBreakTextView) Layout(textViewResolver TextViewResolver, e model.Element, x, y, w, h int) *TextLayout {
	return &TextLayout{
		Element:  e,
		Children: nil,
		Width:    1,
		Height:   1,
	}
}

func (s *SoftBreakTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	// Draw a space for soft break
	renderer.SetContent(0, 0, ' ', nil, tcell.StyleDefault)
}

func (s *SoftBreakTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) (Width int, Height int) {
	return 1, 1
}

// HardBreakTextView renders HardBreakElement
type HardBreakTextView struct{}

func (hb *HardBreakTextView) Layout(textViewResolver TextViewResolver, e model.Element, x, y, w, h int) *TextLayout {
	return &TextLayout{
		Element:  e,
		Children: nil,
		Width:    1,
		Height:   1,
	}
}

func (hb *HardBreakTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	// Hard break creates a line break - no visual content needed
}

func (hb *HardBreakTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) (Width int, Height int) {
	return 1, 1
}

// TableTextView renders TableElement
type TableTextView struct{}

func (t *TableTextView) Layout(textViewResolver TextViewResolver, e model.Element, x, y, w, h int) *TextLayout {
	children := []*TextLayout{}
	var heightTable []int
	for i := 0; i < e.GetElementCount(); i++ {
		row := e.GetElement(i)
		maxHeight := -1
		for j := 0; j < row.GetElementCount(); j++ {
			col := row.GetElement(j)
			colView := textViewResolver.Resolve(col)
			_, mh := colView.MinimumSize(textViewResolver, col, w, h)

			if mh > maxHeight {
				maxHeight = mh
			}
		}
		heightTable = append(heightTable, maxHeight)
	}

	var widthTable []int
	for j := 0; j < e.GetElement(0).GetElementCount(); j++ {
		maxWidth := -1
		for i := 0; i < e.GetElementCount(); i++ {
			cell := e.GetElement(i).GetElement(j)
			cellView := textViewResolver.Resolve(cell)
			mw, _ := cellView.MinimumSize(textViewResolver, cell, w, h)

			if mw > maxWidth {
				maxWidth = mw
			}
		}
		widthTable = append(widthTable, maxWidth)
	}

	totalWidth := 0
	for _, w := range widthTable {
		totalWidth += w
	}

	totalHeight := 0
	yy := 1
	for i, h := range heightTable {
		row := e.GetElement(i)
		rowView := textViewResolver.Resolve(row)
		rowElement := rowView.Layout(textViewResolver, row, 0, yy, w, h)
		children = append(children, rowElement)
		totalHeight += h
		yy += h
	}

	return &TextLayout{
		Element:   e,
		Children:  children,
		Width:     totalWidth + (len(children[0].Children) + 1),
		Height:    totalHeight + 3,
		RelativeX: x,
		RelativeY: y,
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

func (t *TableTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) (Width int, Height int) {
	var heightTable []int
	for i := 0; i < e.GetElementCount(); i++ {
		row := e.GetElement(i)
		maxHeight := -1
		for j := 0; j < row.GetElementCount(); j++ {
			col := row.GetElement(j)
			colView := textViewResolver.Resolve(col)
			_, mh := colView.MinimumSize(textViewResolver, col, width, height)

			if mh > maxHeight {
				maxHeight = mh
			}
		}
		heightTable = append(heightTable, maxHeight)
	}

	var widthTable []int
	for j := 0; j < e.GetElement(0).GetElementCount(); j++ {
		maxWidth := -1
		for i := 0; i < e.GetElementCount(); i++ {
			cell := e.GetElement(i).GetElement(j)
			cellView := textViewResolver.Resolve(cell)
			mw, _ := cellView.MinimumSize(textViewResolver, cell, width, height)

			if mw > maxWidth {
				maxWidth = mw
			}
		}
		widthTable = append(widthTable, maxWidth)
	}

	totalWidth := 0
	for _, w := range widthTable {
		totalWidth += w
	}

	totalHeight := 0
	for _, h := range heightTable {
		totalHeight += h
	}

	return totalWidth + (e.GetElement(0).GetElementCount() + 1), totalHeight + 3
}

// TableHeaderTextView renders TableHeaderElement
type TableHeaderTextView struct{}

func (t *TableHeaderTextView) Layout(textViewResolver TextViewResolver, e model.Element, x, y, w, h int) *TextLayout {
	children := []*TextLayout{}
	totalWidth := 0
	maxHeight := -1
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		mw, mh := childView.MinimumSize(textViewResolver, childElement, w, h)
		childTextLayout := childView.Layout(textViewResolver, childElement, x, y, mw, mh)
		children = append(children, childTextLayout)

		if mh > maxHeight {
			maxHeight = mh
		}
		totalWidth += mw
	}

	return &TextLayout{
		Element:  e,
		Children: children,
		Width:    totalWidth,
		Height:   maxHeight,
	}
}

func (t *TableHeaderTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	style := tcell.StyleDefault.Bold(true)
	cellCount := len(textLayout.Children)
	if cellCount == 0 {
		return
	}

	//availableWidth := textLayout.Width - (cellCount - 1)
	//cellWidth := availableWidth / cellCount

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
		x += child.Width
	}
}

func (t *TableHeaderTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) (Width int, Height int) {
	totalWidth := 0
	maxHeight := -1
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		mw, mh := childView.MinimumSize(textViewResolver, childElement, width, height)

		if mh > maxHeight {
			maxHeight = mh
		}
		totalWidth += mw
	}
	return totalWidth, maxHeight
}

// TableRowTextView renders TableRowElement
type TableRowTextView struct{}

func (t *TableRowTextView) Layout(textViewResolver TextViewResolver, e model.Element, x, y, w, h int) *TextLayout {
	children := []*TextLayout{}
	totalWidth := 0
	maxHeight := -1
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		mw, mh := childView.MinimumSize(textViewResolver, childElement, w, h)
		childTextLayout := childView.Layout(textViewResolver, childElement, x, y, mw, mh)
		children = append(children, childTextLayout)

		if mh > maxHeight {
			maxHeight = mh
		}
		totalWidth += mw
	}

	return &TextLayout{
		Element:  e,
		Children: children,
		Width:    totalWidth,
		Height:   maxHeight,
	}
}

func (t *TableRowTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	cellCount := len(textLayout.Children)
	if cellCount == 0 {
		return
	}

	//availableWidth := textLayout.Width - (cellCount - 1)
	//cellWidth := availableWidth / cellCount

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
		x += child.Width
	}
}

func (t *TableRowTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) (Width int, Height int) {
	totalWidth := 0
	maxHeight := -1
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		mw, mh := childView.MinimumSize(textViewResolver, childElement, width, height)

		if mh > maxHeight {
			maxHeight = mh
		}
		totalWidth += mw
	}
	return totalWidth, maxHeight
}

// TableCellTextView renders TableCellElement
type TableCellTextView struct{}

func (t *TableCellTextView) Layout(textViewResolver TextViewResolver, e model.Element, x, y, w, h int) *TextLayout {
	children := []*TextLayout{}
	totalWidth := 0
	maxHeight := -1
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		mw, mh := childView.MinimumSize(textViewResolver, childElement, w, h)
		childTextLayout := childView.Layout(textViewResolver, childElement, x, y, mw, mh)
		children = append(children, childTextLayout)

		if mh > maxHeight {
			maxHeight = mh
		}
		totalWidth += mw
	}

	return &TextLayout{
		Element:  e,
		Children: children,
		Width:    totalWidth,
		Height:   maxHeight,
	}
}

func (t *TableCellTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	x := 0
	for i, child := range textLayout.Children {
		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		childView.Draw(textViewResolver, child, renderer.Translate(x, 0))
		x += child.Width
	}
}

func (t *TableCellTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) (Width int, Height int) {
	totalWidth := 0
	maxHeight := -1
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		mw, mh := childView.MinimumSize(textViewResolver, childElement, width, height)

		if mh > maxHeight {
			maxHeight = mh
		}
		totalWidth += mw
	}
	return totalWidth, maxHeight
}
