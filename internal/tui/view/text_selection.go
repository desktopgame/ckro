package view

import "github.com/desktopgame/ckro/internal/tui/model"

type TextSelection struct {
	FromPos model.Position
	ToPos   model.Position
}

func (ts TextSelection) IsZero() bool {
	return ts.FromPos.Row == ts.ToPos.Row && ts.FromPos.Column == ts.ToPos.Column
}

func (ts TextSelection) Ordered() (First model.Position, Last model.Position) {
	if ts.FromPos.Row < ts.ToPos.Row {
		return ts.FromPos, ts.ToPos
	} else if ts.FromPos.Row > ts.ToPos.Row {
		return ts.ToPos, ts.FromPos
	}
	if ts.FromPos.Column < ts.ToPos.Column {
		return ts.FromPos, ts.ToPos
	}
	return ts.ToPos, ts.FromPos
}

func (ts TextSelection) Contain(at model.Position) bool {
	if ts.IsZero() {
		return false
	}
	first, last := ts.Ordered()
	if at.Row < first.Row || at.Row > last.Row {
		return false
	}
	if at.Row == first.Row {
		if at.Column < first.Column {
			return false
		}
	}
	if at.Row == last.Row {
		if at.Column >= last.Column {
			return false
		}
	}
	return true
}
