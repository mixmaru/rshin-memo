package views

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/mixmaru/rshin-memo/cmd/rshinmemo/dto"
	"github.com/mixmaru/rshin-memo/cmd/rshinmemo/utils"
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/pkg/errors"
	"sort"
	"time"
)

const DATE_SELECT_VIEW = "date_select"

type DateSelectView struct {
	gui       *gocui.Gui
	view      *gocui.View
	dates     []time.Time
	dateRange DateRange

	insertData  dto.InsertData
	memoDirPath string
	addRowMode  AddRowMode

	dailyDataRepository repositories.DailyDataRepositoryInterface
	noteRepository      repositories.NoteRepositoryInterface

	openViews    []Deletable
	WhenFinished func() error
}

func NewDateSelectView(
	gui *gocui.Gui,
	openView []Deletable,
	insertData dto.InsertData,
	dateRange DateRange,
	memoDirPath string,
	dailyDataRepository repositories.DailyDataRepositoryInterface,
	noteRepository repositories.NoteRepositoryInterface,
) *DateSelectView {
	retObj := &DateSelectView{
		gui:                 gui,
		openViews:           openView,
		insertData:          insertData,
		dateRange:           dateRange,
		memoDirPath:         memoDirPath,
		dailyDataRepository: dailyDataRepository,
		noteRepository:      noteRepository,
	}
	return retObj
}

// 新規作成
func (n *DateSelectView) Create() error {
	dates, err := n.dateRange.GetSomeDateInRange(30)
	if err != nil {
		return err
	}
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].After(dates[j])
	})
	n.dates = dates
	width, height := n.gui.Size()
	v, err := createOrResizeView(n.gui, DATE_SELECT_VIEW, width/2-25, 0, width/2+25, height-1)
	if err != nil {
		return err
	}
	n.view = v

	n.view.Highlight = true
	n.view.SelBgColor = gocui.ColorGreen
	n.view.SelFgColor = gocui.ColorBlack

	n.setContents()
	err = n.setEvents()
	if err != nil {
		return err
	}

	n.openViews = append(n.openViews, n)
	return nil
}

func (n *DateSelectView) setEvents() error {
	// dateSelectView
	if err := n.gui.SetKeybinding(DATE_SELECT_VIEW, gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return errors.Wrap(err, "キーバインド失敗")
	}
	if err := n.gui.SetKeybinding(DATE_SELECT_VIEW, 'j', gocui.ModNone, cursorDown); err != nil {
		return errors.Wrap(err, "キーバインド失敗")
	}
	if err := n.gui.SetKeybinding(DATE_SELECT_VIEW, gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return errors.Wrap(err, "キーバインド失敗")
	}
	if err := n.gui.SetKeybinding(DATE_SELECT_VIEW, 'k', gocui.ModNone, cursorUp); err != nil {
		return errors.Wrap(err, "キーバイーンド失敗")
	}
	if err := n.gui.SetKeybinding(DATE_SELECT_VIEW, gocui.KeyEnter, gocui.ModNone, n.decisionDate); err != nil {
		return errors.Wrap(err, "キーバイーンド失敗")
	}
	return nil
}

func (n *DateSelectView) setContents() {
	fmt.Fprintln(n.view, utils.ConvertStringForView("手入力する"))
	for _, date := range n.dates {
		fmt.Fprintln(n.view, utils.ConvertStringForView(date.Format("2006-01-02")))
	}
}

func (n *DateSelectView) Focus() error {
	_, err := n.gui.SetCurrentView(DATE_SELECT_VIEW)
	if err != nil {
		return errors.Wrap(err, "フォーカス移動失敗")
	}
	return nil
}

func (n *DateSelectView) GetDateOnCursor() (string, error) {
	_, y := n.view.Cursor()
	line, err := n.view.Line(y)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return line, nil
}

func (n *DateSelectView) Delete() error {
	err := n.gui.DeleteView(DATE_SELECT_VIEW)
	if err != nil {
		return errors.Wrapf(err, "Viewの削除に失敗。%v", NOTE_NAME_INPUT_VIEW)
	}
	return nil
}

func (n *DateSelectView) isSelectedHandInput() bool {
	_, y := n.view.Cursor()
	if y == 0 {
		return true
	} else {
		return false
	}
}

func (n *DateSelectView) decisionDate(g *gocui.Gui, v *gocui.View) error {
	if n.isSelectedHandInput() {
		// dateInputViewの表示
		err := n.displayDateInputView()
		if err != nil {
			return err
		}
	} else {
		var err error
		n.insertData.DateStr, err = n.GetDateOnCursor()
		if err != nil {
			return err
		}
		// noteSelectViewの表示
		noteSelectView := NewNoteSelectView(
			n.gui,
			n.insertData,
			n.openViews,
			n.memoDirPath,
			n.dailyDataRepository,
			n.noteRepository,
		)
		useCase := usecases.NewGetAllNotesUseCase(n.noteRepository)
		allNotes, err := useCase.Handle()
		err = noteSelectView.Create(allNotes)
		if err != nil {
			return err
		}
		noteSelectView.WhenFinished = n.WhenFinished
		err = noteSelectView.Focus()
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *DateSelectView) displayDateInputView() error {
	// note名入力viewの表示
	dateInputView := NewDateInputView(
		n.gui,
		n.insertData,
		n.dateRange,
		n.memoDirPath,
		n.openViews,
		n.dailyDataRepository,
		n.noteRepository,
	)
	err := dateInputView.Create()
	if err != nil {
		return err
	}
	dateInputView.WhenFinished = n.WhenFinished
	// フォーカスの移動
	err = dateInputView.Focus()
	if err != nil {
		return errors.Wrap(err, "フォーカス移動失敗")
	}
	return nil
}