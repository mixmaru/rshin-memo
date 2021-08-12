package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type NoteSelectView struct {
	view                                   *tview.Table
	name                                   string
	whenPushEscapeKey                      []func() error
	whenPushEnterKeyOnNoteNameLine         []func(noteName string) error
	whenPushEnterKeyOnInputNewNoteNameLine []func() error
}

func (n *NoteSelectView) GetTviewPrimitive() tview.Primitive {
	return n.view
}

func (n *NoteSelectView) GetName() string {
	return n.name
}

func NewNoteSelectView() *NoteSelectView {
	noteSelectView := &NoteSelectView{name: "note_select_view"}
	noteSelectView.initView()
	return noteSelectView
}

func (n *NoteSelectView) initView() {
	view := tview.NewTable()
	view.SetSelectable(true, false)
	view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			err := n.executeWhenPushEscapeKey()
			if err != nil {
				panic(err)
			}
		case tcell.KeyEnter:
			if n.isSelectedInputNewNoteNameLine() {
				err := n.executeWhenPushEnterKeyOnInputNewNoteNameLine()
				if err != nil {
					panic(err)
				}
			} else {
				row, _ := n.view.GetSelection()
				noteName := n.view.GetCell(row, 0).Text
				err := n.executeWhenPushEnterKeyOnNoteNameLine(noteName)
				if err != nil {
					panic(err)
				}
			}
		}
		return event
	})
	n.view = view
	return
}

const INPUT_NEW_NOTE_NAME = "新規追加"

func (n *NoteSelectView) SetData(notes []string) {
	n.view.Clear()
	n.view.SetCellSimple(0, 0, INPUT_NEW_NOTE_NAME)
	for i, note := range notes {
		n.view.SetCellSimple(i+1, 0, note)
	}
}

func (n *NoteSelectView) AddWhenPushEnterKeyOnNoteNameLine(function func(noteName string) error) {
	n.whenPushEnterKeyOnNoteNameLine = append(n.whenPushEnterKeyOnNoteNameLine, function)
}

func (n *NoteSelectView) AddWhenPushEnterKeyOnInputNewNoteNameLine(function func() error) {
	n.whenPushEnterKeyOnInputNewNoteNameLine = append(n.whenPushEnterKeyOnInputNewNoteNameLine, function)
}

func (n *NoteSelectView) AddWhenPushEscapeKey(function func() error) {
	n.whenPushEscapeKey = append(n.whenPushEscapeKey, function)
}

func (n *NoteSelectView) executeWhenPushEscapeKey() error {
	return executeFunctions(n.whenPushEscapeKey)
}

func (n *NoteSelectView) executeWhenPushEnterKeyOnNoteNameLine(noteName string) error {
	for _, function := range n.whenPushEnterKeyOnNoteNameLine {
		err := function(noteName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *NoteSelectView) executeWhenPushEnterKeyOnInputNewNoteNameLine() error {
	return executeFunctions(n.whenPushEnterKeyOnInputNewNoteNameLine)
}

func (n *NoteSelectView) isSelectedInputNewNoteNameLine() bool {
	row, _ := n.view.GetSelection()
	return n.view.GetCell(row, 0).Text == INPUT_NEW_NOTE_NAME
}
