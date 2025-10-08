package view

import "github.com/desktopgame/ckro/internal/tui/model"

type TextSelection struct {
	FromPos CharacterReference
	ToPos   CharacterReference
}

func (ts TextSelection) IsZero() bool {
	return ts.FromPos.StartPosition.Row == ts.ToPos.StartPosition.Row && ts.FromPos.StartPosition.Column == ts.ToPos.StartPosition.Column
}

func (ts TextSelection) Ordered() (First CharacterReference, Last CharacterReference) {
	if ts.FromPos.StartPosition.Row < ts.ToPos.StartPosition.Row {
		return ts.FromPos, ts.ToPos
	} else if ts.FromPos.StartPosition.Row > ts.ToPos.StartPosition.Row {
		return ts.ToPos, ts.FromPos
	}
	if ts.FromPos.StartPosition.Column < ts.ToPos.StartPosition.Column {
		return ts.FromPos, ts.ToPos
	}
	return ts.ToPos, ts.FromPos
}

func (ts TextSelection) Contain(at model.Position) bool {
	if ts.IsZero() {
		return false
	}
	first, last := ts.Ordered()
	if at.Row < first.StartPosition.Row || at.Row > last.StartPosition.Row {
		return false
	}
	if at.Row == first.StartPosition.Row {
		if at.Column < first.StartPosition.Column {
			return false
		}
	}
	if at.Row == last.StartPosition.Row {
		if at.Column >= last.StartPosition.Column {
			return false
		}
	}
	return true
}
