package tui

import (
	"slices"

	"github.com/desktopgame/ckro/internal/tui/model"
)

type FoldItem struct {
	track   *model.Track
	element model.Element
}

type FoldManager struct {
	items []FoldItem
}

func (fm *FoldManager) addFold(doc model.Document, e model.Element) {
	r := e.GetRange(0)
	tr := doc.CreateTrack(r.StartPosition.Row, r.StartPosition.Column)

	fm.items = append(fm.items, FoldItem{
		track:   tr,
		element: e,
	})
}

func (fm *FoldManager) ToggleFold(doc model.Document, e model.Element) {
	for i, fold := range fm.items {
		if fold.track.Lost {
			continue
		}
		if fm.isMatchRange(fold.track, e) {
			fm.items = slices.Delete(fm.items, i, i+1)
			return
		}
	}
	fm.addFold(doc, e)
}

func (fm *FoldManager) Refresh(doc model.Document) {
	newItems := []FoldItem{}

	for _, fold := range fm.items {
		if !fold.track.Lost {
			newItems = append(newItems, fold)
		}
	}

	fm.items = newItems
}

func (fm *FoldManager) isMatchRange(tr *model.Track, e model.Element) bool {
	r := e.GetRange(0)
	if tr.Position.Row == r.StartPosition.Row && tr.Position.Column == r.StartPosition.Column {
		return true
	}
	return false
}

func (fm *FoldManager) IsFolded(doc model.Document, e model.Element) bool {
	for _, fold := range fm.items {
		if fold.track.Lost {
			continue
		}
		if fm.isMatchRange(fold.track, e) {
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
			if r.EndPosition.Row-r.StartPosition.Row >= 10 {
				fm.addFold(doc, element)
			}
		}
	}
}
