package views

import (
	"github.com/jroimartin/gocui"
	"github.com/mixmaru/rshin-memo/cui_app/utils"
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
	v, err := createOrResizeView(n.gui, NOTE_NAME_INPUT_VIEW, width/2-20, height/2-1, width/2+20, height/2+1)
	if err != nil{
		return err
	}

	v.Editable = true
	v.Editor = &Editor{}
	return nil
}

func (n *NoteNameInputView) Focus() error {
	_, err := n.gui.SetCurrentView(NOTE_NAME_INPUT_VIEW)
	if err != nil {
		return errors.Wrap(err, "フォーカス移動失敗")
	}
	return nil
}

func (n *NoteNameInputView) GetInputNoteName() (string, error) {
	text, err := n.view.Line(0)
	if err != nil {
		return "", errors.Wrap(err, "入力データの取得に失敗しました")
	}
	inputText := utils.ConvertStringForLogic(text)
	return inputText, nil
}

type Editor struct {
}
func (e *Editor)Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		text := utils.ConvertStringForView(string(ch))
		for _, ch := range text {
			v.EditWrite(ch)
		}
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case key == gocui.KeyEnter:
	case key == gocui.KeyArrowDown:
	case key == gocui.KeyArrowUp:
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	}
}
