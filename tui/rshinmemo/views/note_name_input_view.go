package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type NoteNameInputView struct {
	view *tview.InputField
	name string

	whenPushEscapeKey []func() error
	whenPushEnterKey  []func(noteName string) error
}

func NewNoteNameInputView() *NoteNameInputView {
	noteNameInputView := &NoteNameInputView{
		name: "note_name_input_view",
	}
	noteNameInputView.init()
	return noteNameInputView
}

func (n *NoteNameInputView) GetTviewPrimitive() tview.Primitive {
	return n.view
}

func (n *NoteNameInputView) GetName() string {
	return n.name
}

func (n *NoteNameInputView) init() {
	view := tview.NewInputField()
	view.SetLabel("NoteName：")
	view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			err := n.executeWhenPushEscapeKey()
			if err != nil {
				panic(err)
			}
			return nil
		case tcell.KeyEnter:
			err := n.executeWhenPushEnterKey(n.view.GetText())
			if err != nil {
				panic(err)
			}
			return nil
		}
		return event
	})
	n.view = view
}

func (n *NoteNameInputView) executeWhenPushEscapeKey() error {
	return executeFunctions(n.whenPushEscapeKey)
}
func (n *NoteNameInputView) AddWhenPushEscapeKey(function func() error) {
	n.whenPushEscapeKey = append(n.whenPushEscapeKey, function)
}

func (n *NoteNameInputView) executeWhenPushEnterKey(noteName string) error {
	for _, function := range n.whenPushEnterKey {
		err := function(noteName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *NoteNameInputView) AddWhenPushEnterKey(function func(noteName string) error) {
	n.whenPushEnterKey = append(n.whenPushEnterKey, function)
}
