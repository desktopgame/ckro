package view

import "github.com/desktopgame/ckro/internal/tui/model"

type Context struct {
	Resolver TextViewResolver
	Document model.Document
}
