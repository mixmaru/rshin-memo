package views

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/mixmaru/rshin-memo/cui_app/utils"
	"github.com/pkg/errors"
)

const NOTE_SELECT_VIEW = "note_select"

type NoteSelectView struct {
	gui   *gocui.Gui
	view  *gocui.View
	notes []string
}

func NewNoteSelectView(gui *gocui.Gui) *NoteSelectView {
	retObj := &NoteSelectView{
		gui: gui,
	}
	return retObj
}

// 新規作成
func (n *NoteSelectView) Create(notes []string) error {
	n.notes = notes
	width, height := n.gui.Size()
	v, err := createOrResizeView(n.gui, NOTE_SELECT_VIEW, width/2-25, 0, width/2+25, height-1)
	if err != nil {
		return err
	}
	n.view = v

	n.view.Highlight = true
	n.view.SelBgColor = gocui.ColorGreen
	n.view.SelFgColor = gocui.ColorBlack

	n.setContents()

	return nil
}

func (n *NoteSelectView) setContents() {
	fmt.Fprintln(n.view, utils.ConvertStringForView("新規追加"))
	for _, note := range n.notes {
		fmt.Fprintln(n.view, utils.ConvertStringForView(note))
	}
}

func (n *NoteSelectView) Focus() error {
	_, err := n.gui.SetCurrentView(NOTE_SELECT_VIEW)
	if err != nil {
		return errors.Wrap(err, "フォーカス移動失敗")
	}
	return nil
}

func (n *NoteSelectView) GetNoteNameOnCursor() string {
	_, y := n.view.Cursor()
	noteName := n.notes[y-1]
	return noteName
}

func (n *NoteSelectView) Delete() error {
	err := n.gui.DeleteView(NOTE_SELECT_VIEW)
	if err != nil {
		return errors.Wrapf(err, "Viewの削除に失敗。%v", NOTE_NAME_INPUT_VIEW)
	}
	return nil
}
