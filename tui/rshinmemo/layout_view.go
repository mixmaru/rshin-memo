package main

import (
	"github.com/rivo/tview"
)

type layoutView struct {
	layoutView *tview.Pages
}

func newLayoutView() *layoutView {
	return &layoutView{layoutView: tview.NewPages()}
}

func (l *layoutView) AddPage(view viewInterface) {
	l.layoutView.AddPage(view.GetName(), view.GetTviewTable(), true, true)
}

func (l *layoutView) RemovePage(view viewInterface) {
	l.layoutView.RemovePage(view.GetName())
}
