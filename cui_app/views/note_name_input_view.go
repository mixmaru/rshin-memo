package views

import (
	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"
)

const NOTE_NAME_INPUT_VIEW = "note_name_input"

type NoteNameInputView struct {
	gui *gocui.Gui
	view *gocui.View
}

func NewNoteNameinputView(gui *gocui.Gui) *NoteNameInputView {
	retObj := &NoteNameInputView{
		gui: gui,
	}
	return retObj
}

// dailyListViewの新規作成
func (n *NoteNameInputView) Create() error{
	width, height := n.gui.Size()
	_, err := createOrResizeView(n.gui, NOTE_NAME_INPUT_VIEW, width/2-20, height/2-1, width/2+20, height/2+1)
	if err != nil{
		return err
	}
	return nil
}

func (n *NoteNameInputView) Focus() error {
	_, err := n.gui.SetCurrentView(NOTE_NAME_INPUT_VIEW)
	if err != nil {
		return errors.Wrap(err, "フォーカス移動失敗")
	}
	return nil
}