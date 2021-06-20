package views

import (
	"github.com/jroimartin/gocui"
)

type View interface {
	Delete() error
	Focus() error
	AllDelete() error
	deleteThisView(g *gocui.Gui, v *gocui.View) error
}

type ViewBase struct {
	viewName   string
	gui        *gocui.Gui
	parentView View
	childView  View
}

func NewViewBase(viewName string, gui *gocui.Gui, parentView View) *ViewBase {
	return &ViewBase{
		viewName:   viewName,
		gui:        gui,
		parentView: parentView,
	}
}
