package views

import (
	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"
)

type View interface {
	Delete() error
	Focus() error
}

type ViewBase struct {
	viewName string
	gui      *gocui.Gui
}

func NewViewBase(viewName string, gui *gocui.Gui) *ViewBase {
	return &ViewBase{viewName: viewName, gui: gui}
}

func (v *ViewBase) Focus() error {
	_, err := v.gui.SetCurrentView(v.viewName)
	if err != nil {
		return errors.Wrap(err, "フォーカス移動失敗")
	}
	return nil
}

func (v *ViewBase) Delete() error {
	v.gui.DeleteKeybindings(v.viewName)
	err := v.gui.DeleteView(v.viewName)
	if err != nil {
		return errors.Wrapf(err, "Viewの削除に失敗。%v", v.viewName)
	}
	return nil
}
