package views

import (
	"github.com/jroimartin/gocui"
	"github.com/mixmaru/rshin-memo/cmd/rshinmemo/dto"
	"github.com/mixmaru/rshin-memo/cmd/rshinmemo/utils"
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/pkg/errors"
	"path/filepath"
)

const NOTE_NAME_INPUT_VIEW = "note_name_input"

type NoteNameInputView struct {
	gui         *gocui.Gui
	view        *gocui.View
	memoDirPath string

	insertData dto.InsertData

	dailyDataRepository repositories.DailyDataRepositoryInterface
	noteRepository      repositories.NoteRepositoryInterface

	*ViewBase
}

func (n *NoteNameInputView) Delete() error {
	return deleteView(n.gui, n.viewName)
}

func (n *NoteNameInputView) Focus() error {
	return focus(n.gui, n.viewName)
}

func (n *NoteNameInputView) AllDelete() error {
	return allDelete(n, n.parentView)
}

func (n *NoteNameInputView) deleteThisView(g *gocui.Gui, v *gocui.View) error {
	return deleteThisView(n, n.parentView)
}

func (n *NoteNameInputView) Resize() error {
	width, height := n.gui.Size()
	return resize(n.gui, n.viewName, width/2-50, height/2-1, width/2+50, height/2+1, n.childView)
}

func NewNoteNameInputView(
	gui *gocui.Gui,
	memoDirPath string,
	insertData dto.InsertData,
	parentView View,
	dailyDataRepository repositories.DailyDataRepositoryInterface,
	noteRepository repositories.NoteRepositoryInterface,
) *NoteNameInputView {
	retObj := &NoteNameInputView{
		gui:                 gui,
		memoDirPath:         memoDirPath,
		insertData:          insertData,
		dailyDataRepository: dailyDataRepository,
		noteRepository:      noteRepository,
	}
	retObj.ViewBase = NewViewBase(NOTE_NAME_INPUT_VIEW, gui, parentView)
	return retObj
}

// dailyListViewの新規作成
func (n *NoteNameInputView) Create() error {
	width, height := n.gui.Size()
	v, err := createOrResizeView(n.gui, NOTE_NAME_INPUT_VIEW, width/2-50, height/2-1, width/2+50, height/2+1)
	if err != nil {
		return err
	}
	n.view = v

	n.view.Editable = true
	n.view.Editor = &Editor{}

	err = n.setEvents()
	if err != nil {
		return err
	}
	return nil
}

func (n *NoteNameInputView) setEvents() error {
	// inputNoteNameViewでのEnterキー
	if err := n.gui.SetKeybinding(NOTE_NAME_INPUT_VIEW, gocui.KeyEnter, gocui.ModNone, n.createNote); err != nil {
		return errors.WithStack(err)
	}
	if err := n.gui.SetKeybinding(NOTE_NAME_INPUT_VIEW, gocui.KeyEsc, gocui.ModNone, n.deleteThisView); err != nil {
		return errors.Wrap(err, "キーバイーンド失敗")
	}
	return nil
}

func (n *NoteNameInputView) getInputNoteName() (string, error) {
	text, err := n.view.Line(0)
	if err != nil {
		return "", errors.Wrap(err, "入力データの取得に失敗しました")
	}
	inputText := utils.ConvertStringForLogic(text)
	return inputText, nil
}

func (n *NoteNameInputView) createNote(gui *gocui.Gui, view *gocui.View) error {
	// 入力内容を取得
	noteName, err := n.getInputNoteName()
	if err != nil {
		return err
	}
	n.insertData.NoteName = noteName

	// 同名Noteが存在しないかcheck
	useCase := usecases.NewGetNoteUseCase(n.noteRepository)
	_, notExist, err := useCase.Handle(noteName)
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

	err = n.AllDelete()
	if err != nil {
		return err
	}
	return nil
}

func (n *NoteNameInputView) createNewDailyList(insertData dto.InsertData) error {
	dailyData, err := insertData.GenerateNewDailyData()
	if err != nil {
		return err
	}
	// Note作成を依頼
	useCase := usecases.NewSaveDailyDataUseCase(n.noteRepository, n.dailyDataRepository)
	err = useCase.Handle(dailyData)
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
