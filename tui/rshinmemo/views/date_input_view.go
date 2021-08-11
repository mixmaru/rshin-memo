package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type DateInputView struct {
	view              *tview.InputField
	name              string
	whenPushEscapeKey []func() error
	whenPushEnterKey  []func(inputDateStr string) (*ValidationError, error)
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
		case tcell.KeyEnter:
			validErr, err := d.executeWhenPushEnterKey()
			if err != nil {
				panic(err)
			}
			if validErr != nil {
				// バリデーションエラー。何もしない。
				return nil
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

type ValidationError struct {
	No      int
	Message string
}

const VALIDATION_ERROR_1_MESSAGE = "formatがYYYY-MM-DDではない"

func (d *DateInputView) executeWhenPushEnterKey() (*ValidationError, error) {

	for _, function := range d.whenPushEnterKey {
		validErr, err := function(d.view.GetText())
		if err != nil {
			return nil, err
		}
		if validErr != nil {
			return validErr, nil
		}
	}
	return nil, nil
}

func (d *DateInputView) AddWhenPushEnterKey(function func(inputDateStr string) (*ValidationError, error)) {
	d.whenPushEnterKey = append(d.whenPushEnterKey, function)
}
