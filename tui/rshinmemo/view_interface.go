package main

import "github.com/rivo/tview"

type viewInterface interface {
	GetTviewTable() *tview.Table
	GetName() string
}
