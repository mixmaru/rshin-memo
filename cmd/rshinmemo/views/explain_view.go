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
	width, height := e.gui.Size()
	v, err := createOrResizeView(e.gui, EXPLAIN_VIEW, 0, height-3, width-1, height-1)
	if err != nil {
		return err
	}
	e.view = v
	fmt.Fprintln(e.view, text)
	return nil
}
