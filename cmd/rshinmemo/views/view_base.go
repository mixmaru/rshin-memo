package views

import (
	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"
)

type View interface {
	Delete() error
	Focus() error
	AllDelete() error
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

func (vb *ViewBase) Focus() error {
	_, err := vb.gui.SetCurrentView(vb.viewName)
	if err != nil {
		return errors.Wrapf(err, "フォーカス移動失敗。%+v", vb)
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

// 一番上の親のview以外、全ての親を削除し、フォーカスを設定する
func (vb *ViewBase) AllDelete() error {
	if vb.parentView != nil {
		if err := vb.Delete(); err != nil {
			return err
		}
		return vb.parentView.AllDelete()
	} else {
		return vb.Focus()
	}
}

func (vb *ViewBase) deleteThisView(g *gocui.Gui, v *gocui.View) error {
	err := vb.Delete()
	if err != nil {
		return err
	}
	err = vb.parentView.Focus()
	if err != nil {
		return err
	}
	return nil
}
