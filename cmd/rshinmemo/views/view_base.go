package views

import (
	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"
)

type View interface {
	Delete() error
	Focus() error
	AllDelete() error
	deleteThisView(g *gocui.Gui, v *gocui.View) error
	Resize() error
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

func deleteThisView(view View, parentView View) error {
	err := view.Delete()
	if err != nil {
		return err
	}
	err = parentView.Focus()
	if err != nil {
		return err
	}
	return nil
}

func focus(gui *gocui.Gui, viewName string) error {
	_, err := gui.SetCurrentView(viewName)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func deleteView(gui *gocui.Gui, viewName string) error {
	gui.DeleteKeybindings(viewName)
	err := gui.DeleteView(viewName)
	if err != nil {
		return errors.Wrapf(err, "Viewの削除に失敗。%+v", viewName)
	}
	return nil
}

func allDelete(view, parentView View) error {
	if parentView != nil {
		if err := view.Delete(); err != nil {
			return err
		}
		return parentView.AllDelete()
	} else {
		return view.Focus()
	}
}

func resize(gui *gocui.Gui, currentViewName string, x0, y0, x1, y1 int, childView View) error {
	_, err := gui.View(currentViewName)
	if err != nil {
		if err == gocui.ErrUnknownView {
			// viewが存在しなければ(既にdeleteされているとかで)なにもしない
			return nil
		} else {
			return errors.WithStack(err)
		}
	}
	_, err = createOrResizeView(gui, currentViewName, x0, y0, x1, y1)
	if err != nil {
		return err
	}
	if childView != nil {
		return childView.Resize()
	}
	return nil
}
