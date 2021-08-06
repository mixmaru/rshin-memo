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

func (l *layoutView) AddPage(name string, view *tview.Table) {
	l.layoutView.AddPage(name, view, true, true)
}

func (l *layoutView) RemovePage(name string, view *tview.Table) {
	l.layoutView.RemovePage(name)
}
