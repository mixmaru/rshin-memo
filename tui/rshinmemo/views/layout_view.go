package views

import (
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

type LayoutView struct {
	app        *tview.Application
	layoutView *tview.Pages
}

func NewLayoutView() *LayoutView {
	return &LayoutView{
		app:        tview.NewApplication(),
		layoutView: tview.NewPages(),
	}
}

func (l *LayoutView) AddPage(view viewInterface) {
	l.layoutView.AddPage(view.GetName(), view.GetTviewPrimitive(), true, true)
}

func (l *LayoutView) RemovePage(view viewInterface) {
	l.layoutView.RemovePage(view.GetName())
}

func (l *LayoutView) SetRoot(view viewInterface) {
	l.app.SetRoot(l.layoutView, true)
}

func (l *LayoutView) Run() error {
	err := l.app.Run()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (l *LayoutView) Refresh() {
	l.app.Sync()
	return
}

func (l *LayoutView) Suspend(f func()) bool {
	return l.app.Suspend(f)
}

func (l *LayoutView) Stop() {
	l.app.Stop()
}
