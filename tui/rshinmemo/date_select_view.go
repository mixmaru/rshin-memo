package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"time"
)

type DateSelectView struct {
	view              *tview.Table
	whenPushEscapeKey []func() error
	whenPushEnterKey  []func() error
}

func NewDateSelectView() *DateSelectView {
	dateSelectView := &DateSelectView{}
	dateSelectView.initView()
	return dateSelectView
}

func (d *DateSelectView) initView() {
	view := tview.NewTable().SetSelectable(true, false)
	// event設定
	view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			err := d.executeWhenPushEscapeKey()
			if err != nil {
				panic(err)
			}
			return nil
		case tcell.KeyEnter:
			err := d.executeWhenPushEnterKey()
			if err != nil {
				panic(err)
			}
		}
		return event
	})
	d.view = view
}

func (d *DateSelectView) AddWhenPushEnterKey(function func() error) {
	d.whenPushEnterKey = append(d.whenPushEnterKey, function)
}

func (d *DateSelectView) AddWhenPushEscapeKey(function func() error) {
	d.whenPushEscapeKey = append(d.whenPushEscapeKey, function)
}

func (d *DateSelectView) executeWhenPushEscapeKey() error {
	return executeFunctions(d.whenPushEscapeKey)
}

func (d *DateSelectView) executeWhenPushEnterKey() error {
	return executeFunctions(d.whenPushEnterKey)
}

func (d *DateSelectView) SetData(dates []time.Time) {
	d.view.Clear()
	d.view.SetCellSimple(0, 0, "手入力する")
	for i, date := range dates {
		d.view.SetCellSimple(i+1, 0, date.Format("2006-01-02"))
	}
}

//func (r *RshinMemo) createInitDailySelectView(mode usecases.InsertMode) (*tview.Table, error) {
//	dateSelectView := tview.NewTable().SetSelectable(true, false)
//	// event設定
//	dateSelectView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
//		var err error
//		switch event.Key() {
//		case tcell.KeyEscape:
//			// dateSelectViewを削除してDailyListにフォーカスを戻す
//			r.closeDateSelectView()
//			return nil
//		case tcell.KeyEnter:
//			// noteSelectViewを表示してフォーカスを移す
//			r.noteSelectView, err = r.createNoteSelectView()
//			if err != nil {
//				panic(err)
//			}
//			r.layoutView.AddPage("noteSelectView", r.noteSelectView)
//		}
//		return event
//	})
//	dateSelectView.SetCellSimple(0, 0, "手入力する")
//
//	// 表示する日付の範囲を決定する
//	overCurrentDate, err := r.getDailyListCursorDate(-1)
//	if err != nil {
//		return nil, err
//	}
//	currentDate, err := r.getDailyListCursorDate(0)
//	if err != nil {
//		return nil, err
//	}
//	underCurrentDate, err := r.getDailyListCursorDate(1)
//	if err != nil {
//		return nil, err
//	}
//
//	now := time.Now().In(time.Local)
//	useCase := usecases.NewGetDateSelectRangeUseCase(now)
//	dates, err := useCase.Handle(overCurrentDate, currentDate, underCurrentDate, mode)
//	if err != nil {
//		return nil, err
//	}
//	for i, date := range dates {
//		dateSelectView.SetCellSimple(i+1, 0, date.Format("2006-01-02"))
//	}
//	return dateSelectView, nil
//}
