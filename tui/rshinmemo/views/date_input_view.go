package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type DateInputView struct {
	view              *tview.InputField
	name              string
	whenPushEscapeKey []func() error
}

func NewDateInputView() *DateInputView {
	dateInputView := &DateInputView{
		name: "date_input_view",
	}
	dateInputView.init()
	return dateInputView
}

func (d *DateInputView) GetTviewPrimitive() tview.Primitive {
	return d.view
}

func (d *DateInputView) GetName() string {
	return d.name
}

func (d *DateInputView) init() {
	view := tview.NewInputField()
	view.SetLabel("日付入力(YYYY-MM-DD)：")
	view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			err := d.executeWhenPushEscapeKey()
			if err != nil {
				panic(err)
			}
			return nil
		}
		return event
	})
	d.view = view
}

func (d *DateInputView) executeWhenPushEscapeKey() error {
	return executeFunctions(d.whenPushEscapeKey)
}

func (d *DateInputView) AddWhenPushEscapeKey(function func() error) {
	d.whenPushEscapeKey = append(d.whenPushEscapeKey, function)
}
