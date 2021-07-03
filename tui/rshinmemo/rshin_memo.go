package main

import (
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

type RshinMemo struct {
	app           *tview.Application
	layoutView    *tview.Flex
	dailyListView *tview.Table

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
	layoutView = tview.NewFlex().AddItem(dailyListView, 100, 0, true)
	return layoutView, dailyListView, nil
}

func (r *RshinMemo) createInitDailyListView() (*tview.Table, error) {
	table := tview.NewTable()
	table.SetSelectable(true, false)
	table.SetBorder(true)
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
