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
)

const DATE_SELECT_VIEW = "date_select"

type DateSelectView struct {
	gui         *gocui.Gui
	view        *gocui.View
	memoDirPath string

	insertData dto.InsertData
	dateRange  DateRange

	dailyDataRepository repositories.DailyDataRepositoryInterface
	noteRepository      repositories.NoteRepositoryInterface

	*ViewBase
}

func (n *DateSelectView) deleteThisView(g *gocui.Gui, v *gocui.View) error {
	return deleteThisView(n, n.parentView)
}

func (n *DateSelectView) Delete() error {
	return deleteView(n.gui, n.viewName)
}

func (n *DateSelectView) Focus() error {
	return focus(n.gui, n.viewName)
}

func (n *DateSelectView) AllDelete() error {
	return allDelete(n, n.parentView)
}

func (n *DateSelectView) Resize() error {
	width, height := n.gui.Size()
	return resize(n.gui, n.viewName, width/2-25, 0, width/2+25, height-1, n.childView)
}

func NewDateSelectView(
	gui *gocui.Gui,
	memoDirPath string,
	insertData dto.InsertData,
	dateRange DateRange,
	parentView View,
	dailyDataRepository repositories.DailyDataRepositoryInterface,
	noteRepository repositories.NoteRepositoryInterface,
) *DateSelectView {
	retObj := &DateSelectView{
		gui:                 gui,
		insertData:          insertData,
		dateRange:           dateRange,
		memoDirPath:         memoDirPath,
		dailyDataRepository: dailyDataRepository,
		noteRepository:      noteRepository,
	}
	retObj.ViewBase = NewViewBase(DATE_SELECT_VIEW, gui, parentView)
	return retObj
}

// 新規作成
func (n *DateSelectView) Create() error {
	width, height := n.gui.Size()
	v, err := createOrResizeView(n.gui, DATE_SELECT_VIEW, width/2-25, 0, width/2+25, height-1)
	if err != nil {
		return err
	}
	n.view = v

	n.view.Highlight = true
	n.view.SelBgColor = gocui.ColorGreen
	n.view.SelFgColor = gocui.ColorBlack

	err = n.setContents()
	if err != nil {
		return err
	}
	err = n.setEvents()
	if err != nil {
		return err
	}

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
	if err := n.gui.SetKeybinding(DATE_SELECT_VIEW, gocui.KeyEsc, gocui.ModNone, n.deleteThisView); err != nil {
		return errors.Wrap(err, "キーバイーンド失敗")
	}
	return nil
}

func (n *DateSelectView) deleteEvents() {
	n.gui.DeleteKeybindings(DATE_SELECT_VIEW)
}

func (n *DateSelectView) setContents() error {
	fmt.Fprintln(n.view, utils.ConvertStringForView("手入力する"))

	dates, err := n.dateRange.GetSomeDateInRange(30)
	if err != nil {
		return err
	}
	// 日付の並びをdailyListと合わせる
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].After(dates[j])
	})

	for _, date := range dates {
		fmt.Fprintln(n.view, utils.ConvertStringForView(date.Format("2006-01-02")))
	}
	return nil
}

func (n *DateSelectView) getDateOnCursor() (string, error) {
	_, y := n.view.Cursor()
	line, err := n.view.Line(y)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return line, nil
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
		n.insertData.DateStr, err = n.getDateOnCursor()
		if err != nil {
			return err
		}
		// 全note名の取得
		useCase := usecases.NewGetAllNotesUseCase(n.noteRepository)
		allNotes, err := useCase.Handle()
		// noteSelectViewの表示
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
		n.childView = noteSelectView
	}
	return nil
}

func (n *DateSelectView) displayDateInputView() error {
	// note名入力viewの表示
	dateInputView := NewDateInputView(
		n.gui,
		n.memoDirPath,
		n.insertData,
		n.dateRange,
		n,
		n.dailyDataRepository,
		n.noteRepository,
	)
	err := dateInputView.Create()
	if err != nil {
		return err
	}
	// フォーカスの移動
	err = dateInputView.Focus()
	if err != nil {
		return errors.Wrap(err, "フォーカス移動失敗")
	}
	n.childView = dateInputView
	return nil
}
