package views

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"
)

const NOTE_SELECT_VIEW = "note_select"

type NoteSelectView struct {
	gui  *gocui.Gui
	view *gocui.View
}

func NewNoteSelectView(gui *gocui.Gui) *NoteSelectView {
	retObj := &NoteSelectView{
		gui: gui,
	}
	return retObj
}

// 新規作成
func (n *NoteSelectView) Create(notes []string) error {
	width, height := n.gui.Size()
	v, err := createOrResizeView(n.gui, NOTE_SELECT_VIEW, width/2-25, 0, width/2+25, height-1)
	if err != nil {
		return err
	}
	n.view = v

	n.view.Highlight = true
	n.view.SelBgColor = gocui.ColorGreen
	n.view.SelFgColor = gocui.ColorBlack

	n.setContents(notes)

	return nil
}

func (n *NoteSelectView) setContents(notes []string) {
	fmt.Fprintln(n.view, "新規追加")
	for _, note := range notes {
		fmt.Fprintln(n.view, note)
	}
}

func (n *NoteSelectView) Focus() error {
	_, err := n.gui.SetCurrentView(NOTE_SELECT_VIEW)
	if err != nil {
		return errors.Wrap(err, "フォーカス移動失敗")
	}
	return nil
}
