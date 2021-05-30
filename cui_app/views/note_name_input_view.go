package views

import (
	"github.com/jroimartin/gocui"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/mixmaru/rshin-memo/cui_app/dto"
	"github.com/mixmaru/rshin-memo/cui_app/utils"
	"github.com/pkg/errors"
	"path/filepath"
)

const NOTE_NAME_INPUT_VIEW = "note_name_input"

type NoteNameInputView struct {
	gui          *gocui.Gui
	view         *gocui.View
	insertData   dto.InsertData
	memoDirPath  string
	WhenFinished func() error // call when finish

	getNoteUseCase       *usecases.GetNoteUseCase
	saveDailyDataUseCase *usecases.SaveDailyDataUseCase
}

func NewNoteNameinputView(
	gui *gocui.Gui,
	memoDirPath string,
	insertData dto.InsertData,
	getNoteUseCase *usecases.GetNoteUseCase,
	saveDailyDataUseCase *usecases.SaveDailyDataUseCase,
) *NoteNameInputView {
	retObj := &NoteNameInputView{
		gui:                  gui,
		insertData:           insertData,
		memoDirPath:          memoDirPath,
		getNoteUseCase:       getNoteUseCase,
		saveDailyDataUseCase: saveDailyDataUseCase,
	}
	return retObj
}

// dailyListViewの新規作成
func (n *NoteNameInputView) Create() error {
	width, height := n.gui.Size()
	v, err := createOrResizeView(n.gui, NOTE_NAME_INPUT_VIEW, width/2-20, height/2-1, width/2+20, height/2+1)
	if err != nil {
		return err
	}
	n.view = v

	n.view.Editable = true
	n.view.Editor = &Editor{}

	// inputNoteNameViewでのEnterキー
	if err := n.gui.SetKeybinding(NOTE_NAME_INPUT_VIEW, gocui.KeyEnter, gocui.ModNone, n.createNote); err != nil {
		return errors.WithStack(err)
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

func (n *NoteNameInputView) GetInputNoteName() (string, error) {
	text, err := n.view.Line(0)
	if err != nil {
		return "", errors.Wrap(err, "入力データの取得に失敗しました")
	}
	inputText := utils.ConvertStringForLogic(text)
	return inputText, nil
}

func (n *NoteNameInputView) Delete() error {
	err := n.gui.DeleteView(NOTE_NAME_INPUT_VIEW)
	if err != nil {
		return errors.Wrapf(err, "Viewの削除に失敗。%v", NOTE_NAME_INPUT_VIEW)
	}
	return nil
}

func (n *NoteNameInputView) createNote(gui *gocui.Gui, view *gocui.View) error {
	// 入力内容を取得
	noteName, err := n.GetInputNoteName()
	if err != nil {
		return err
	}
	n.insertData.SetNoteName(noteName)

	// 同名Noteが存在しないかcheck
	_, notExist, err := n.getNoteUseCase.Handle(noteName)
	if err != nil {
		return err
	} else if !notExist {
		// すでに同名のNoteが存在する
		// todo: エラーメッセージビューへメッセージを表示する
	} else {
		if err := n.createNewDailyList(n.insertData); err != nil {
			return err
		}

		err = utils.OpenVim(filepath.Join(n.memoDirPath, noteName+".txt"))
		if err != nil {
			return err
		}
	}

	return n.WhenFinished()
}

func (n *NoteNameInputView) createNewDailyList(insertData dto.InsertData) error {
	dailyData := insertData.GenerateNewDailyData()
	// Note作成を依頼
	err := n.saveDailyDataUseCase.Handle(dailyData)
	if err != nil {
		// todo: エラーメッセージビューへメッセージを表示する
		return err
	}
	return nil
}

type Editor struct {
}

func (e *Editor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
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
