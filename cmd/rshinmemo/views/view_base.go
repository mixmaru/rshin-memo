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
	viewName  string
	gui       *gocui.Gui
	openViews []View
}

func NewViewBase(viewName string, gui *gocui.Gui, openViews []View) *ViewBase {
	return &ViewBase{
		viewName:  viewName,
		gui:       gui,
		openViews: openViews,
	}
}

func (vb *ViewBase) Focus() error {
	_, err := vb.gui.SetCurrentView(vb.viewName)
	if err != nil {
		return errors.Wrap(err, "フォーカス移動失敗")
	}
	return nil
}

func (vb *ViewBase) Delete() error {
	vb.gui.DeleteKeybindings(vb.viewName)
	err := vb.gui.DeleteView(vb.viewName)
	if err != nil {
		return errors.Wrapf(err, "Viewの削除に失敗。%vb", vb.viewName)
	}
	return nil
}

func (vb *ViewBase) deleteThisView(g *gocui.Gui, v *gocui.View) error {
	err := vb.Delete()
	if err != nil {
		return err
	}
	err = vb.openViews[len(vb.openViews)-2].Focus()
	if err != nil {
		return err
	}
	return nil
}
