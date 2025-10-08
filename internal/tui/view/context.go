package view

import "github.com/desktopgame/ckro/internal/tui/model"

type Context struct {
	Resolver      TextViewResolver
	Document      model.Document
	FoldManager   FoldManager
	TextSelection TextSelection
}

func (ctx Context) GetSegment(e model.Element, rangeIndex int) model.Segment {
	r := e.GetRange(rangeIndex)
	return ctx.Document.Read(r)
}

func (ctx Context) GetText(e model.Element) string {
	sg := ctx.GetSegment(e, 0)
	if sg.GetLineCount() != 1 {
		panic("")
	}
	return sg.GetLine(0)
}
