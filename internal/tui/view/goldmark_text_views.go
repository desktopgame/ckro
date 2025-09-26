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

func (d *DocumentTextView) Layout(textViewResolver TextViewResolver, textLayout *TextLayout, x, y, w, h int) {
	offsetY := 0
	for i := 0; i < len(textLayout.Children); i++ {
		childElement := textLayout.Children[i].Element
		childView := textViewResolver.Resolve(childElement)

		mw := textLayout.Children[i].MinimumWidth
		mh := textLayout.Children[i].MinimumHeight
		childView.Layout(textViewResolver, textLayout.Children[i], 0, offsetY, mw, mh)
		offsetY += textLayout.Children[i].Height
	}

	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
}

func (d *DocumentTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	y := 0
	for i, child := range textLayout.Children {
		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		childView.Draw(textViewResolver, child, renderer.Translate(child.RelativeX, child.RelativeY))
		y += child.Height
	}
}

func (d *DocumentTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) *TextLayout {
	totalHeight := 0
	maxWidth := 0
	children := []*TextLayout{}
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)

		child := childView.MinimumSize(textViewResolver, childElement, width, 9999)
		children = append(children, child)

		totalHeight += child.MinimumHeight
		if child.MinimumWidth > maxWidth {
			maxWidth = child.MinimumWidth
		}
	}
	return &TextLayout{
		Element:       e,
		MinimumWidth:  max(maxWidth, width),
		MinimumHeight: totalHeight,
	}
}

// ParagraphTextView renders ParagraphElement
type ParagraphTextView struct{}

func (p *ParagraphTextView) Layout(textViewResolver TextViewResolver, textLayout *TextLayout, x, y, w, h int) {
	offsetX := 0
	for i := 0; i < len(textLayout.Children); i++ {
		childElement := textLayout.Children[i].Element
		childView := textViewResolver.Resolve(childElement)

		mw := textLayout.Children[i].MinimumWidth
		childView.Layout(textViewResolver, textLayout.Children[i], offsetX, 0, mw, 1)
		offsetX += textLayout.Children[i].Width
	}

	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
}

func (p *ParagraphTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	for i, child := range textLayout.Children {
		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		childView.Draw(textViewResolver, child, renderer.Translate(child.RelativeX, child.RelativeY))
	}
}

func (p *ParagraphTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) *TextLayout {
	totalWidth := 0
	children := []*TextLayout{}
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)

		child := childView.MinimumSize(textViewResolver, childElement, width, 1)
		children = append(children, child)
		totalWidth += child.MinimumWidth
	}
	return &TextLayout{
		Element:       e,
		MinimumWidth:  totalWidth,
		MinimumHeight: 1,
		Children:      children,
	}
}

// HeadingTextView renders HeadingElement
type HeadingTextView struct{}

func (ht *HeadingTextView) Layout(textViewResolver TextViewResolver, textLayout *TextLayout, x, y, w, h int) {
	headingElement := textLayout.Element.(*model.HeadingElement)
	prefixWidth := headingElement.Level + 1 // "# " or "## " etc.
	offsetX := prefixWidth
	for i := 0; i < len(textLayout.Children); i++ {
		childElement := textLayout.Children[i].Element
		childView := textViewResolver.Resolve(childElement)

		mw := textLayout.Children[i].MinimumWidth
		mh := textLayout.Children[i].MinimumHeight
		childView.Layout(textViewResolver, textLayout.Children[i], offsetX, 0, mw, mh)
		offsetX += textLayout.Children[i].Width
	}

	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
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

func (h *HeadingTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) *TextLayout {
	totalWidth := 0
	headingElement := e.(*model.HeadingElement)
	prefixWidth := headingElement.Level + 1 // "# " or "## " etc.
	children := []*TextLayout{}
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)

		child := childView.MinimumSize(textViewResolver, childElement, width, 1)
		children = append(children, child)

		totalWidth += child.MinimumWidth
	}
	return &TextLayout{
		Element:       e,
		MinimumWidth:  totalWidth + prefixWidth,
		MinimumHeight: 1,
		Children:      children,
	}
}

// CodeBlockTextView renders CodeBlockElement
type CodeBlockTextView struct{}

func (c *CodeBlockTextView) Layout(textViewResolver TextViewResolver, textLayout *TextLayout, x, y, w, h int) {
	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
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

func (c *CodeBlockTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) *TextLayout {
	codeBlock := e.(*model.CodeBlockElement)
	lines := strings.Split(codeBlock.Text, "\n")

	maxWidth := -1
	for _, line := range lines {
		w := text.DisplayWidth(line)
		if w > maxWidth {
			maxWidth = w
		}
	}
	return &TextLayout{
		Element:       e,
		MinimumWidth:  maxWidth,
		MinimumHeight: len(lines),
	}
}

// BlockquoteTextView renders BlockquoteElement
type BlockquoteTextView struct{}

func (b *BlockquoteTextView) Layout(textViewResolver TextViewResolver, textLayout *TextLayout, x, y, w, h int) {
	for i := 0; i < len(textLayout.Children); i++ {
		childElement := textLayout.Children[i].Element
		childView := textViewResolver.Resolve(childElement)
		mw := textLayout.Children[i].MinimumWidth
		childView.Layout(textViewResolver, textLayout.Children[i], 2, 0, mw, 1)
	}

	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
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

func (b *BlockquoteTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) *TextLayout {
	maxWidth := -1
	totalHeight := 0
	children := []*TextLayout{}
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		child := childView.MinimumSize(textViewResolver, childElement, width, 1)
		children = append(children, child)

		if child.MinimumWidth > maxWidth {
			maxWidth = child.MinimumWidth
		}
		totalHeight++
	}
	return &TextLayout{
		Element:       e,
		MinimumWidth:  maxWidth,
		MinimumHeight: totalHeight,
		Children:      children,
	}
}

// ListTextView renders ListElement
type ListTextView struct{}

func (l *ListTextView) Layout(textViewResolver TextViewResolver, textLayout *TextLayout, x, y, w, h int) {
	listElement := textLayout.Element.(*model.ListElement)
	var prefix string
	if listElement.Ordered {
		prefix = strings.Repeat(" ", 2) + "1" + ". "
	} else {
		prefix = "  - "
	}

	offsetY := 0
	for i := 0; i < len(textLayout.Children); i++ {
		childElement := textLayout.Children[i].Element
		childView := textViewResolver.Resolve(childElement)
		mw := textLayout.Children[i].MinimumWidth
		mh := textLayout.Children[i].MinimumHeight
		childView.Layout(textViewResolver, textLayout.Children[i], text.DisplayWidth(prefix), offsetY, mw, mh)
		offsetY += textLayout.Children[i].Height
	}

	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
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

		x := 0
		for _, r := range prefix {
			renderer.SetContent(x, y, r, nil, tcell.StyleDefault)
			x++
		}

		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		childView.Draw(textViewResolver, child, renderer.Translate(child.RelativeX, child.RelativeY))
		y += child.Height
	}
}

func (l *ListTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) *TextLayout {
	listElement := e.(*model.ListElement)
	var prefix string
	if listElement.Ordered {
		prefix = strings.Repeat(" ", 2) + "1" + ". "
	} else {
		prefix = "  - "
	}
	offset := text.DisplayWidth(prefix)

	maxWidth := -1
	totalHeight := 0
	children := []*TextLayout{}
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		child := childView.MinimumSize(textViewResolver, childElement, width-offset, 9999)
		children = append(children, child)

		if child.MinimumWidth > maxWidth {
			maxWidth = child.MinimumWidth
		}
		totalHeight += child.MinimumHeight
	}
	return &TextLayout{
		Element:       e,
		MinimumWidth:  maxWidth + offset,
		MinimumHeight: totalHeight,
		Children:      children,
	}
}

// ListItemTextView renders ListItemElement
type ListItemTextView struct{}

func (l *ListItemTextView) Layout(textViewResolver TextViewResolver, textLayout *TextLayout, x, y, w, h int) {
	offsetX := 0
	offsetY := 0
	for i := 0; i < len(textLayout.Children); i++ {
		childElement := textLayout.Children[i].Element
		childView := textViewResolver.Resolve(childElement)
		mw := textLayout.Children[i].MinimumWidth
		mh := textLayout.Children[i].MinimumHeight
		if _, ok := childView.(*ListTextView); ok {
			offsetY++
			offsetX = 0
		}
		childView.Layout(textViewResolver, textLayout.Children[i], offsetX, offsetY, mw, mh)
		offsetX += textLayout.Children[i].Width
	}

	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
}

func (l *ListItemTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	for i, child := range textLayout.Children {
		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		childView.Draw(textViewResolver, child, renderer.Translate(child.RelativeX, child.RelativeY))
	}
}

func (l *ListItemTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) *TextLayout {
	totalWidth := 0
	maxHeight := -1
	lists := 0
	children := []*TextLayout{}
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		child := childView.MinimumSize(textViewResolver, childElement, width, height)
		children = append(children, child)

		if _, ok := childView.(*ListTextView); ok {
			lists++
		}

		if child.MinimumHeight > maxHeight {
			maxHeight = child.MinimumHeight
		}
		totalWidth += child.MinimumWidth
	}
	return &TextLayout{
		Element:       e,
		MinimumWidth:  totalWidth,
		MinimumHeight: maxHeight + lists,
		Children:      children,
	}
}

// ThematicBreakTextView renders ThematicBreakElement
type ThematicBreakTextView struct{}

func (t *ThematicBreakTextView) Layout(textViewResolver TextViewResolver, textLayout *TextLayout, x, y, w, h int) {
	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
}

func (t *ThematicBreakTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	// Draw horizontal line
	for x := 0; x < textLayout.Width; x++ {
		renderer.SetContent(x, 0, '-', nil, tcell.StyleDefault)
	}
}

func (t *ThematicBreakTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) *TextLayout {
	return &TextLayout{
		Element:       e,
		MinimumWidth:  1,
		MinimumHeight: 1,
	}
}

// TextElementView renders TextElement
type TextElementView struct{}

func (t *TextElementView) Layout(textViewResolver TextViewResolver, textLayout *TextLayout, x, y, w, h int) {
	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
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

func (t *TextElementView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) *TextLayout {
	return &TextLayout{
		Element:       e,
		MinimumWidth:  text.DisplayWidth(e.GetText()),
		MinimumHeight: 1,
	}
}

// EmphasisTextView renders EmphasisElement
type EmphasisTextView struct{}

func (em *EmphasisTextView) Layout(textViewResolver TextViewResolver, textLayout *TextLayout, x, y, w, h int) {
	for i := 0; i < len(textLayout.Children); i++ {
		childElement := textLayout.Children[i].Element
		childView := textViewResolver.Resolve(childElement)
		mw := textLayout.Children[i].MinimumWidth
		childView.Layout(textViewResolver, textLayout.Children[i], 0, 0, mw, 1)
	}

	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
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

func (em *EmphasisTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) *TextLayout {
	totalWidth := 0
	children := []*TextLayout{}
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		child := childView.MinimumSize(textViewResolver, childElement, width, 1)
		children = append(children, child)
		totalWidth += child.MinimumWidth
	}
	return &TextLayout{
		Element:       e,
		MinimumWidth:  totalWidth,
		MinimumHeight: 1,
		Children:      children,
	}
}

// StrongTextView renders StrongElement
type StrongTextView struct{}

func (s *StrongTextView) Layout(textViewResolver TextViewResolver, textLayout *TextLayout, x, y, w, h int) {
	for i := 0; i < len(textLayout.Children); i++ {
		childElement := textLayout.Children[i].Element
		childView := textViewResolver.Resolve(childElement)
		mw := textLayout.Children[i].MinimumWidth
		childView.Layout(textViewResolver, textLayout.Children[i], x, y, mw, 1)
	}

	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
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

func (s *StrongTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) *TextLayout {
	totalWidth := 0
	children := []*TextLayout{}
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		child := childView.MinimumSize(textViewResolver, childElement, width, 1)
		children = append(children, child)
		totalWidth += child.MinimumWidth
	}
	return &TextLayout{
		Element:       e,
		MinimumWidth:  totalWidth,
		MinimumHeight: 1,
		Children:      children,
	}
}

// CodeElementView renders CodeElement
type CodeElementView struct{}

func (c *CodeElementView) Layout(textViewResolver TextViewResolver, textLayout *TextLayout, x, y, w, h int) {
	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
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

func (c *CodeElementView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) *TextLayout {
	return &TextLayout{
		Element:       e,
		MinimumWidth:  text.DisplayWidth(e.GetText()),
		MinimumHeight: 1,
	}
}

// LinkTextView renders LinkElement
type LinkTextView struct{}

func (l *LinkTextView) Layout(textViewResolver TextViewResolver, textLayout *TextLayout, x, y, w, h int) {
	for i := 0; i < len(textLayout.Children); i++ {
		childElement := textLayout.Children[i].Element
		childView := textViewResolver.Resolve(childElement)
		mw := textLayout.Children[i].MinimumWidth
		childView.Layout(textViewResolver, textLayout.Children[i], x, y, mw, 1)
	}

	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
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

func (l *LinkTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) *TextLayout {
	totalWidth := 0
	children := []*TextLayout{}
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		child := childView.MinimumSize(textViewResolver, childElement, width, 1)
		children = append(children, child)
		totalWidth += child.MinimumWidth
	}
	return &TextLayout{
		Element:       e,
		MinimumWidth:  totalWidth,
		MinimumHeight: 1,
		Children:      children,
	}
}

// ImageTextView renders ImageElement
type ImageTextView struct{}

func (i *ImageTextView) Layout(textViewResolver TextViewResolver, textLayout *TextLayout, x, y, w, h int) {
	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
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

func (i *ImageTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) *TextLayout {
	imageElement := e.(*model.ImageElement)
	displayText := "[Image: " + imageElement.Alt + "]"
	return &TextLayout{
		Element:       e,
		MinimumWidth:  text.DisplayWidth(displayText),
		MinimumHeight: 1,
	}
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

func (s *SoftBreakTextView) Layout(textViewResolver TextViewResolver, textLayout *TextLayout, x, y, w, h int) {
	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
}

func (s *SoftBreakTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	// Draw a space for soft break
	renderer.SetContent(0, 0, ' ', nil, tcell.StyleDefault)
}

func (s *SoftBreakTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) *TextLayout {
	return &TextLayout{
		Element:       e,
		MinimumWidth:  1,
		MinimumHeight: 1,
	}
}

// HardBreakTextView renders HardBreakElement
type HardBreakTextView struct{}

func (hb *HardBreakTextView) Layout(textViewResolver TextViewResolver, textLayout *TextLayout, x, y, w, h int) {
	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
}

func (hb *HardBreakTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	// Hard break creates a line break - no visual content needed
}

func (hb *HardBreakTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) *TextLayout {
	return &TextLayout{
		Element:       e,
		MinimumWidth:  1,
		MinimumHeight: 1,
	}
}

// TableTextView renders TableElement
type TableTextView struct{}

func (t *TableTextView) Layout(textViewResolver TextViewResolver, textLayout *TextLayout, x, y, w, h int) {
	var heightTable []int
	for i := 0; i < len(textLayout.Children); i++ {
		row := textLayout.Children[i].Element
		maxHeight := -1
		for j := 0; j < row.GetElementCount(); j++ {
			mh := textLayout.Children[i].MinimumHeight

			if mh > maxHeight {
				maxHeight = mh
			}
		}
		heightTable = append(heightTable, maxHeight)
	}

	var widthTable []int
	for j := 0; j < len(textLayout.Children[0].Children); j++ {
		maxWidth := -1
		for i := 0; i < len(textLayout.Children); i++ {
			mw := textLayout.Children[i].Children[j].MinimumWidth

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
		row := textLayout.Children[i].Element
		rowView := textViewResolver.Resolve(row)

		textLayout.Children[i].WidthTable = widthTable
		rowView.Layout(textViewResolver, textLayout.Children[i], 1, yy, w, h)
		if _, ok := row.(*model.TableHeaderElement); ok {
			h++
		}
		totalHeight += h
		yy += h
	}

	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
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
		childView.Draw(textViewResolver, child, renderer.Translate(child.RelativeX, child.RelativeY)) // Offset by left border
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

func (t *TableTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) *TextLayout {
	children := []*TextLayout{}
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		child := childView.MinimumSize(textViewResolver, childElement, width, height)
		children = append(children, child)
	}

	var heightTable []int
	for i := 0; i < e.GetElementCount(); i++ {
		row := e.GetElement(i)
		maxHeight := -1
		for j := 0; j < row.GetElementCount(); j++ {
			mh := children[i].Children[j].MinimumHeight

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
			mw := children[i].Children[j].MinimumWidth

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

	return &TextLayout{
		Element:       e,
		MinimumWidth:  totalWidth + (e.GetElement(0).GetElementCount() + 1),
		MinimumHeight: totalHeight + 3,
		Children:      children,
	}
}

// TableHeaderTextView renders TableHeaderElement
type TableHeaderTextView struct{}

func (t *TableHeaderTextView) Layout(textViewResolver TextViewResolver, textLayout *TextLayout, x, y, w, h int) {
	offsetX := 0
	for i := 0; i < len(textLayout.Children); i++ {
		childElement := textLayout.Children[i].Element
		childView := textViewResolver.Resolve(childElement)
		// mw := textLayout.Children[i].MinimumWidth
		mh := textLayout.Children[i].MinimumHeight
		childView.Layout(textViewResolver, textLayout.Children[i], offsetX, 0, textLayout.WidthTable[i], mh)
		offsetX += textLayout.Children[i].Width + 1
	}

	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
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

		childView.Draw(textViewResolver, child, renderer.Translate(child.RelativeX, child.RelativeY))
		x += textLayout.WidthTable[i]
	}
}

func (t *TableHeaderTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) *TextLayout {
	totalWidth := 0
	maxHeight := -1
	children := []*TextLayout{}
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		child := childView.MinimumSize(textViewResolver, childElement, width, height)
		children = append(children, child)

		if child.MinimumHeight > maxHeight {
			maxHeight = child.MinimumHeight
		}
		totalWidth += child.MinimumWidth
	}
	return &TextLayout{
		Element:       e,
		MinimumWidth:  totalWidth,
		MinimumHeight: maxHeight,
		Children:      children,
	}
}

// TableRowTextView renders TableRowElement
type TableRowTextView struct{}

func (t *TableRowTextView) Layout(textViewResolver TextViewResolver, textLayout *TextLayout, x, y, w, h int) {
	offsetX := 0
	for i := 0; i < len(textLayout.Children); i++ {
		childElement := textLayout.Children[i].Element
		childView := textViewResolver.Resolve(childElement)
		//mw := textLayout.Children[i].MinimumWidth
		mh := textLayout.Children[i].MinimumHeight
		childView.Layout(textViewResolver, textLayout.Children[i], offsetX, 0, textLayout.WidthTable[i], mh)
		offsetX += textLayout.WidthTable[i] + 1
	}

	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
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
		childView.Draw(textViewResolver, child, renderer.Translate(child.RelativeX, child.RelativeY))
		x += textLayout.WidthTable[i]
	}
}

func (t *TableRowTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) *TextLayout {
	totalWidth := 0
	maxHeight := -1
	children := []*TextLayout{}
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		child := childView.MinimumSize(textViewResolver, childElement, width, height)
		children = append(children, child)

		if child.MinimumHeight > maxHeight {
			maxHeight = child.MinimumHeight
		}
		totalWidth += child.MinimumWidth
	}
	return &TextLayout{
		Element:       e,
		MinimumWidth:  totalWidth,
		MinimumHeight: maxHeight,
		Children:      children,
	}
}

// TableCellTextView renders TableCellElement
type TableCellTextView struct{}

func (t *TableCellTextView) Layout(textViewResolver TextViewResolver, textLayout *TextLayout, x, y, w, h int) {
	offsetX := 0
	for i := 0; i < len(textLayout.Children); i++ {
		childElement := textLayout.Children[i].Element
		childView := textViewResolver.Resolve(childElement)
		mw := textLayout.Children[i].MinimumWidth
		mh := textLayout.Children[i].MinimumHeight
		childView.Layout(textViewResolver, textLayout.Children[i], offsetX, 0, mw, mh)
		offsetX += textLayout.Children[i].Width
	}

	textLayout.RelativeX = x
	textLayout.RelativeY = y
	textLayout.Width = w
	textLayout.Height = h
}

func (t *TableCellTextView) Draw(textViewResolver TextViewResolver, textLayout *TextLayout, renderer Renderer) {
	for i, child := range textLayout.Children {
		childElement := textLayout.Element.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		childView.Draw(textViewResolver, child, renderer.Translate(child.RelativeX, child.RelativeY))
	}
}

func (t *TableCellTextView) MinimumSize(textViewResolver TextViewResolver, e model.Element, width int, height int) *TextLayout {
	totalWidth := 0
	maxHeight := -1
	children := []*TextLayout{}
	for i := 0; i < e.GetElementCount(); i++ {
		childElement := e.GetElement(i)
		childView := textViewResolver.Resolve(childElement)
		child := childView.MinimumSize(textViewResolver, childElement, width, height)
		children = append(children, child)

		if child.MinimumHeight > maxHeight {
			maxHeight = child.MinimumHeight
		}
		totalWidth += child.MinimumWidth
	}
	return &TextLayout{
		Element:       e,
		MinimumWidth:  totalWidth,
		MinimumHeight: maxHeight,
		Children:      children,
	}
}
