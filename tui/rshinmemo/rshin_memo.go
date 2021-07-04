package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"time"
)

type RshinMemo struct {
	app            *tview.Application
	layoutView     *tview.Flex
	dailyListView  *tview.Table
	dateSelectView *tview.Table

	dailyDataRep repositories.DailyDataRepositoryInterface
}

func NewRshinMemo(dailyDataRep repositories.DailyDataRepositoryInterface) *RshinMemo {
	return &RshinMemo{
		dailyDataRep: dailyDataRep,
	}
}

func (r *RshinMemo) Run() error {
	var err error
	r.app = tview.NewApplication()
	r.layoutView, r.dailyListView, err = r.createInitViews()
	if err != nil {
		return err
	}

	if err := r.app.SetRoot(r.layoutView, true).Run(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (r *RshinMemo) createInitViews() (layoutView *tview.Flex, dailyListView *tview.Table, err error) {
	dailyListView, err = r.createInitDailyListView()
	if err != nil {
		return nil, nil, err
	}
	layoutView = tview.NewFlex().AddItem(dailyListView, 0, 1, true)
	return layoutView, dailyListView, nil
}

func (r *RshinMemo) createInitDailyListView() (*tview.Table, error) {
	table := tview.NewTable()
	table.SetSelectable(true, false)
	// イベント設定
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		var err error
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'o', 'O':
				// dataSelectViewを作る
				var mode usecases.InsertMode
				if event.Rune() == 'o' {
					mode = usecases.INSERT_UNDER_MODE
				} else {
					mode = usecases.INSERT_OVER_MODE
				}
				r.dateSelectView, err = r.createInitDailySelectView(mode)
				if err != nil {
					panic(errors.WithStack(err))
				}
				// 表示領域に挿入する
				r.layoutView.AddItem(r.dateSelectView, 0, 1, true)
				// フォーカスを移す
				r.app.SetFocus(r.dateSelectView)
				return nil
			}
		}
		return event
	})
	// データ取得
	useCase := usecases.NewGetAllDailyListUsecase(r.dailyDataRep)
	dailyList, err := useCase.Handle()
	if err != nil {
		return nil, err
	}

	// データをテーブルにセット
	for i, dailyData := range dailyList {
		for j, note := range dailyData.Notes {
			row := i + j
			table.SetCellSimple(row, 0, dailyData.Date)
			table.SetCellSimple(row, 1, note)
		}
	}
	return table, nil
}

func (r *RshinMemo) createInitDailySelectView(mode usecases.InsertMode) (*tview.Table, error) {
	dateSelectView := tview.NewTable().SetSelectable(true, false)
	dateSelectView.SetCellSimple(0, 0, "手入力する")

	// 表示する日付の範囲を決定する
	overCurrentDate, err := r.getDailyListCursorDate(-1)
	if err != nil {
		return nil, err
	}
	currentDate, err := r.getDailyListCursorDate(0)
	if err != nil {
		return nil, err
	}
	underCurrentDate, err := r.getDailyListCursorDate(1)
	if err != nil {
		return nil, err
	}

	now := time.Now().In(time.Local)
	useCase := usecases.NewGetDateSelectRangeUseCase(now)
	dates, err := useCase.Handle(overCurrentDate, currentDate, underCurrentDate, mode)
	if err != nil {
		return nil, err
	}
	for i, date := range dates {
		dateSelectView.SetCellSimple(i+1, 0, date.Format("2006-01-02"))
	}
	return dateSelectView, nil
}

// dailyListのカーソル位置の日付を取得する。
// cursorPointAdjustに数値を指定すると、指定分カーソル位置からずれた位置の日付を取得する
func (r *RshinMemo) getDailyListCursorDate(cursorPointAdjust int) (time.Time, error) {
	row, _ := r.dailyListView.GetSelection()
	targetRow := row + cursorPointAdjust
	if targetRow < 0 || targetRow+1 > r.dailyListView.GetRowCount() {
		// targetRow is out of range
		return time.Time{}, nil
	}
	dateStr := r.dailyListView.GetCell(targetRow, 0)
	date, err := time.ParseInLocation("2006-01-02", dateStr.Text, time.Local)
	if err != nil {
		return time.Time{}, errors.WithStack(err)
	}
	return date, nil
}

//list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
//	switch event.Key() {
//	case tcell.KeyEnter:
//		flex.AddItem(textView, 0, 1, false)
//		app.SetFocus(textView)
//		return nil
//	case tcell.KeyEsc:
//		flex.RemoveItem(textView)
//		return nil
//	}
//	return event
//})
//
//if err := app.SetRoot(flex, true).Run(); err != nil {
//	panic(err)
//}
