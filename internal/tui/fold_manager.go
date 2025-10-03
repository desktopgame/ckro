package tui

import (
	"github.com/desktopgame/ckro/internal/optional"
	"github.com/desktopgame/ckro/internal/tui/model"
)

type FoldItem struct {
	prevLine optional.Optional[string]
	nextLine optional.Optional[string]
	element  model.Element
}

type FoldManager struct {
	items []FoldItem
}

func (fm *FoldManager) AddFold(doc model.Document, e model.Element) {
	prevLine, nextLine := fm.scan(doc, e)

	fm.items = append(fm.items, FoldItem{
		prevLine: prevLine,
		nextLine: nextLine,
		element:  e,
	})
}

func (fm *FoldManager) IsFolded(doc model.Document, e model.Element) bool {
	prevLine, nextLine := fm.scan(doc, e)

	for _, item := range fm.items {
		if item.prevLine.IsSame(prevLine) && item.nextLine.IsSame(nextLine) {
			return true
		}
	}

	return false
}

func (fm *FoldManager) AutoFold(doc model.Document, elements []model.Element) {
	fm.items = nil

	for _, element := range elements {
		if fold, ok := element.(*model.FoldBlockElement); ok {
			r := fold.GetRange(0)
			if r.EndPosition.Row-r.StartPosition.Row > 10 {
				fm.AddFold(doc, element)
			}
		}
	}
}

func (fm *FoldManager) scan(doc model.Document, e model.Element) (PrevLine optional.Optional[string], NextLine optional.Optional[string]) {
	r := e.GetRange(0)

	prevLine := optional.None[string]()
	if r.StartPosition.Row > 0 {
		prevRange := model.Range{
			StartPosition: model.Position{
				Row:    r.StartPosition.Row - 1,
				Column: 0,
			},
			EndPosition: model.Position{
				Row:    r.StartPosition.Row - 1,
				Column: doc.GetLineBytes(r.StartPosition.Row - 1),
			},
		}
		sg := doc.Read(prevRange)
		prevLine = optional.Some(sg.GetLine(0))
	}

	nextLine := optional.None[string]()
	if r.EndPosition.Row+1 < doc.GetLineCount() {
		prevRange := model.Range{
			StartPosition: model.Position{
				Row:    r.EndPosition.Row + 1,
				Column: 0,
			},
			EndPosition: model.Position{
				Row:    r.EndPosition.Row + 1,
				Column: doc.GetLineBytes(r.EndPosition.Row + 1),
			},
		}
		sg := doc.Read(prevRange)
		nextLine = optional.Some(sg.GetLine(0))
	}

	return prevLine, nextLine
}
