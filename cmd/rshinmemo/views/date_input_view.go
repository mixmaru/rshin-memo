package views

import (
	"github.com/jroimartin/gocui"
	"github.com/mixmaru/rshin-memo/cmd/rshinmemo/dto"
	"github.com/mixmaru/rshin-memo/cmd/rshinmemo/utils"
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/pkg/errors"
	"time"
)

const DATE_INPUT_VIEW = "date_input"

type DateInputView struct {
	gui         *gocui.Gui
	view        *gocui.View
	memoDirPath string

	insertData dto.InsertData
	dateRange  DateRange

	dailyDataRepository repositories.DailyDataRepositoryInterface
	noteRepository      repositories.NoteRepositoryInterface

	*ViewBase
}

func (n *DateInputView) Delete() error {
	return deleteView(n.gui, n.viewName)
}

func (n *DateInputView) Focus() error {
	return focus(n.gui, n.viewName)
}

func (n *DateInputView) AllDelete() error {
	return allDelete(n, n.parentView)
}

func (n *DateInputView) deleteThisView(g *gocui.Gui, v *gocui.View) error {
	return deleteThisView(n, n.parentView)
}

func NewDateInputView(
	gui *gocui.Gui,
	memoDirPath string,
	insertData dto.InsertData,
	dateRange DateRange,
	parentView View,
	dailyDataRepository repositories.DailyDataRepositoryInterface,
	noteRepository repositories.NoteRepositoryInterface,
) *DateInputView {
	retObj := &DateInputView{
		gui:                 gui,
		insertData:          insertData,
		dateRange:           dateRange,
		memoDirPath:         memoDirPath,
		dailyDataRepository: dailyDataRepository,
		noteRepository:      noteRepository,
	}
	retObj.ViewBase = NewViewBase(DATE_INPUT_VIEW, gui, parentView)
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

	err = n.setEvent()
	if err != nil {
		return err
	}

	return nil
}

func (n *DateInputView) setEvent() error {
	// DateInputViewでのEnterキー
	if err := n.gui.SetKeybinding(DATE_INPUT_VIEW, gocui.KeyEnter, gocui.ModNone, n.displayNoteNameInputView); err != nil {
		return errors.Wrap(err, "Enterキーバインド失敗")
	}
	if err := n.gui.SetKeybinding(DATE_INPUT_VIEW, gocui.KeyEsc, gocui.ModNone, n.deleteThisView); err != nil {
		return errors.Wrap(err, "Enterキーバインド失敗")
	}
	return nil
}

func (n *DateInputView) getInputString() (string, error) {
	text, err := n.view.Line(0)
	if err != nil {
		return "", errors.Wrap(err, "入力データの取得に失敗しました")
	}
	inputText := utils.ConvertStringForLogic(text)
	return inputText, nil
}

func (n *DateInputView) displayNoteNameInputView(g *gocui.Gui, v *gocui.View) error {
	// 日付入力値の取得
	dateString, err := n.getInputString()
	// 日付入力値のバリデーション
	result, err := n.valid(dateString)
	if err != nil {
		return err
	}
	if !result {
		return nil
	}
	n.insertData.DateStr = dateString

	// noteSelectViewの表示
	useCase := usecases.NewGetAllNotesUseCase(n.noteRepository)
	allNotes, err := useCase.Handle()
	noteSelectView := NewNoteSelectView(
		n.gui,
		n.memoDirPath,
		n.insertData,
		n,
		n.dailyDataRepository,
		n.noteRepository,
	)
	err = noteSelectView.Create(allNotes)
	if err != nil {
		return err
	}
	err = noteSelectView.Focus()
	if err != nil {
		return err
	}

	return nil
}

func (n *DateInputView) valid(dateString string) (bool, error) {
	// 指定のdate文字列がRangeの範囲にとどまっているかをチェック

	targetDate, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		// パース失敗（入力フォーマットが違う）
		return false, nil
	}
	if !n.dateRange.IsIn(targetDate) {
		return false, nil
	}
	return true, nil
}
