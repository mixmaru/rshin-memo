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

	openViews    []Deletable
	WhenFinished func() error

	dailyDataRepository repositories.DailyDataRepositoryInterface
	noteRepository      repositories.NoteRepositoryInterface
}

func NewDateInputView(
	gui *gocui.Gui,
	memoDirPath string,
	insertData dto.InsertData,
	dateRange DateRange,
	openViews []Deletable,
	dailyDataRepository repositories.DailyDataRepositoryInterface,
	noteRepository repositories.NoteRepositoryInterface,
) *DateInputView {
	retObj := &DateInputView{
		gui:                 gui,
		insertData:          insertData,
		dateRange:           dateRange,
		openViews:           openViews,
		memoDirPath:         memoDirPath,
		dailyDataRepository: dailyDataRepository,
		noteRepository:      noteRepository,
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

	// DateInputViewでのEnterキー
	if err := n.gui.SetKeybinding(DATE_INPUT_VIEW, gocui.KeyEnter, gocui.ModNone, n.displayNoteNameInputView); err != nil {
		return errors.Wrap(err, "Enterキーバインド失敗")
	}

	n.openViews = append(n.openViews, n)
	return nil
}

func (n *DateInputView) Focus() error {
	_, err := n.gui.SetCurrentView(DATE_INPUT_VIEW)
	if err != nil {
		return errors.Wrap(err, "フォーカス移動失敗")
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

func (n *DateInputView) Delete() error {
	err := n.gui.DeleteView(DATE_INPUT_VIEW)
	if err != nil {
		return errors.Wrapf(err, "Viewの削除に失敗。%v", NOTE_NAME_INPUT_VIEW)
	}
	return nil
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
		n.insertData,
		n.openViews,
		n.memoDirPath,
		n.dailyDataRepository,
		n.noteRepository,
	)
	err = noteSelectView.Create(allNotes)
	if err != nil {
		return err
	}
	noteSelectView.WhenFinished = n.WhenFinished
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
