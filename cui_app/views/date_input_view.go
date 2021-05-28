package views

import (
	"github.com/jroimartin/gocui"
	"github.com/mixmaru/rshin-memo/cui_app/utils"
	"github.com/pkg/errors"
)

const DATE_INPUT_VIEW = "date_input"

type DateInputView struct {
	gui  *gocui.Gui
	view *gocui.View
}

func NewDateInputView(gui *gocui.Gui) *DateInputView {
	retObj := &DateInputView{
		gui: gui,
	}
	return retObj
}

// 新規作成
func (n *DateInputView) Create() error {
	width, height := n.gui.Size()
	v, err := createOrResizeView(n.gui, DATE_INPUT_VIEW, width/2-20, height/2-1, width/2+20, height/2+1)
	if err != nil {
		return err
	}
	n.view = v

	n.view.Editable = true
	n.view.Editor = &Editor{}
	return nil
}

func (n *DateInputView) Focus() error {
	_, err := n.gui.SetCurrentView(DATE_INPUT_VIEW)
	if err != nil {
		return errors.Wrap(err, "フォーカス移動失敗")
	}
	return nil
}

func (n *DateInputView) GetInputString() (string, error) {
	text, err := n.view.Line(0)
	if err != nil {
		return "", errors.Wrap(err, "入力データの取得に失敗しました")
	}
	inputText := utils.ConvertStringForLogic(text)
	return inputText, nil
}

func (n *DateInputView) Delete() error {
	err := n.gui.DeleteView(DATE_INPUT_VIEW)
	if err != nil {
		return errors.Wrapf(err, "Viewの削除に失敗。%v", NOTE_NAME_INPUT_VIEW)
	}
	return nil
}
