package view

import "github.com/desktopgame/ckro/internal/tui/model"

// Context is rendering parameter set of basic.
type Context struct {
	Resolver      TextViewResolver
	Document      model.Document
	FoldManager   FoldManager
	TextSelection TextSelection
}

// GetSegment returns Range at specified index.
func (ctx Context) GetSegment(e model.Element, rangeIndex int) model.Segment {
	r := e.GetRange(rangeIndex)
	return ctx.Document.Read(r)
}

// GetText returns text of Element
func (ctx Context) GetText(e model.Element) string {
	sg := ctx.GetSegment(e, 0)
	if sg.GetLineCount() != 1 {
		panic("text is many lines")
	}
	return sg.GetLine(0)
}
