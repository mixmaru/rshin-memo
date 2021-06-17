package views

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type ExplainView struct {
	gui  *gocui.Gui
	view *gocui.View
}

func NewExplainView(
	gui *gocui.Gui,
) *ExplainView {
	return &ExplainView{
		gui: gui,
	}
}

const EXPLAIN_VIEW = "explain_view"

func (e *ExplainView) Create(text string) error {
	v, err := e.createOrResizeView()
	if err != nil {
		return err
	}
	e.view = v
	fmt.Fprintln(e.view, text)
	return nil
}

func (e *ExplainView) Set(text string) {
	e.view.Clear()
	fmt.Fprintln(e.view, text)
}

func (e *ExplainView) Resize() error {
	_, err := e.createOrResizeView()
	if err != nil {
		return err
	}
	return nil
}

func (e *ExplainView) createOrResizeView() (*gocui.View, error) {
	width, height := e.gui.Size()
	v, err := createOrResizeView(e.gui, EXPLAIN_VIEW, 0, height-2, width-1, height)
	if err != nil {
		return nil, err
	}
	v.Frame = false
	return v, nil
}
