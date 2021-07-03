package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
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
			case 'o':
				// dataSelectViewを作る
				r.dateSelectView, err = r.createInitDailySelectView()
				if err != nil {
					panic(errors.WithStack(err))
				}
				// 表示領域に挿入する
				r.layoutView.AddItem(r.dateSelectView, 0, 1, true)
				// フォーカスを移す
				r.app.SetFocus(r.dateSelectView)
				return nil
			case 'O':
				panic("noteImplement")
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

func (r *RshinMemo) createInitDailySelectView() (*tview.Table, error) {
	dateSelectView := tview.NewTable().SetSelectable(true, false)
	dateSelectView.SetCellSimple(0, 0, "2021-01-01")
	dateSelectView.SetCellSimple(1, 0, "2021-01-01")
	dateSelectView.SetCellSimple(2, 0, "2021-01-01")
	dateSelectView.SetCellSimple(3, 0, "2021-01-01")
	return dateSelectView, nil
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
