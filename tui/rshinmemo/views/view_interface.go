package views

import "github.com/rivo/tview"

type viewInterface interface {
	GetTviewPrimitive() tview.Primitive
	GetName() string
}
