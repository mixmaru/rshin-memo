package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type NoteSelectView struct {
	view              *tview.Table
	name              string
	whenPushEscapeKey []func() error
	whenPushEnterKey  []func() error
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
			err := n.executeWhenPushEnterKey()
			if err != nil {
				panic(err)
			}
		}
		return event
	})
	n.view = view
	return
}

func (n *NoteSelectView) SetData(notes []string) {
	n.view.Clear()
	n.view.SetCellSimple(0, 0, "新規追加")
	for i, note := range notes {
		n.view.SetCellSimple(i+1, 0, note)
	}
}

func (n *NoteSelectView) AddWhenPushEnterKey(function func() error) {
	n.whenPushEnterKey = append(n.whenPushEnterKey, function)
}

func (n *NoteSelectView) AddWhenPushEscapeKey(function func() error) {
	n.whenPushEscapeKey = append(n.whenPushEscapeKey, function)
}

func (n *NoteSelectView) executeWhenPushEscapeKey() error {
	return executeFunctions(n.whenPushEscapeKey)
}

func (n *NoteSelectView) executeWhenPushEnterKey() error {
	return executeFunctions(n.whenPushEnterKey)
}

func (n *NoteSelectView) GetSelection() (row, column int) {
	return n.view.GetSelection()
}

func (n *NoteSelectView) GetCell(row int, column int) *tview.TableCell {
	return n.view.GetCell(row, column)
}
