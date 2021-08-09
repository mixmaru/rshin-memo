package views

import (
	"github.com/rivo/tview"
)

type LayoutView struct {
	layoutView *tview.Pages
}

func newLayoutView() *LayoutView {
	return &LayoutView{layoutView: tview.NewPages()}
}

func (l *LayoutView) AddPage(view viewInterface) {
	l.layoutView.AddPage(view.GetName(), view.GetTviewTable(), true, true)
}

func (l *LayoutView) RemovePage(view viewInterface) {
	l.layoutView.RemovePage(view.GetName())
}
