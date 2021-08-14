package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type NoteSelectView struct {
	grid                                   *tview.Grid
	table                                  *tview.Table
	message                                *tview.TextView
	searchInputField                       *tview.InputField
	name                                   string
	whenPushEscapeKey                      []func() error
	whenPushEnterKeyOnNoteNameLine         []func(noteName string) error
	whenPushEnterKeyOnInputNewNoteNameLine []func() error
	whenPushCtrlFKey                       []func() error
	whenPushEnterKeyOnSearchInputField     []func(searchWord string) error
}

func (n *NoteSelectView) GetTviewPrimitive() tview.Primitive {
	return n.grid
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
	table := tview.NewTable()
	table.SetSelectable(true, false)
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			err := n.executeWhenPushEscapeKey()
			if err != nil {
				panic(err)
			}
			return nil
		case tcell.KeyEnter:
			if n.isSelectedInputNewNoteNameLine() {
				err := n.executeWhenPushEnterKeyOnInputNewNoteNameLine()
				if err != nil {
					panic(err)
				}
			} else {
				row, _ := n.table.GetSelection()
				noteName := n.table.GetCell(row, 0).Text
				err := n.executeWhenPushEnterKeyOnNoteNameLine(noteName)
				if err != nil {
					panic(err)
				}
			}
			return nil
		case tcell.KeyCtrlF:
			err := n.executeWhenPushCtrlFKey()
			if err != nil {
				panic(err)
			}
			return nil
		}
		return event
	})
	message := tview.NewTextView().SetText("[esc]:back [j]:up [k]:down [enter]:open memo [ctrl-f]:search memo")
	grid := tview.NewGrid().SetRows(0, 1)
	grid.AddItem(table, 0, 0, 1, 1, 0, 0, true)
	grid.AddItem(message, 1, 0, 1, 1, 0, 0, false)
	n.table = table
	n.message = message
	n.grid = grid
	return
}

const INPUT_NEW_NOTE_NAME = "新規追加"

func (n *NoteSelectView) SetData(notes []string) {
	n.table.Clear()
	n.table.SetCellSimple(0, 0, INPUT_NEW_NOTE_NAME)
	for i, note := range notes {
		n.table.SetCellSimple(i+1, 0, note)
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

func (n *NoteSelectView) AddWhenPushEnterKeyOnSearchInputField(function func(searchWord string) error) {
	n.whenPushEnterKeyOnSearchInputField = append(n.whenPushEnterKeyOnSearchInputField, function)
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
	row, _ := n.table.GetSelection()
	return n.table.GetCell(row, 0).Text == INPUT_NEW_NOTE_NAME
}

func (n *NoteSelectView) executeWhenPushCtrlFKey() error {
	return executeFunctions(n.whenPushCtrlFKey)
}

func (n *NoteSelectView) AddWhenPushCtrlFKey(function func() error) {
	n.whenPushCtrlFKey = append(n.whenPushCtrlFKey, function)
}

func (n *NoteSelectView) SearchMode(layout *LayoutView) {
	n.searchInputField = n.createSearchInputField(layout)
	n.grid.RemoveItem(n.message)
	n.grid.AddItem(n.searchInputField, 1, 0, 1, 1, 0, 0, true)
	layout.app.SetFocus(n.searchInputField)
}
func (n *NoteSelectView) NoteSelectMode(layout *LayoutView) {
	n.grid.RemoveItem(n.searchInputField)
	n.grid.AddItem(n.message, 1, 0, 1, 1, 0, 0, true)
	n.searchInputField = nil
	layout.app.SetFocus(n.table)
}

func (n *NoteSelectView) createSearchInputField(layout *LayoutView) *tview.InputField {
	searchInputField := tview.NewInputField().SetLabel("Search:")
	searchInputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			n.NoteSelectMode(layout)
			return nil
		case tcell.KeyEnter:
			searchWord := n.searchInputField.GetText()
			err := n.executeWhenPushEnterKeyOnSearchInputField(searchWord)
			if err != nil {
				panic(err)
			}
			return nil
		}
		return event
	})
	return searchInputField
}

func (n *NoteSelectView) executeWhenPushEnterKeyOnSearchInputField(searchWord string) error {
	for _, function := range n.whenPushEnterKeyOnSearchInputField {
		err := function(searchWord)
		if err != nil {
			return err
		}
	}
	return nil
}
