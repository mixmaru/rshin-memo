package views

import "github.com/rivo/tview"

type viewInterface interface {
	GetTviewTable() *tview.Table
	GetName() string
}
