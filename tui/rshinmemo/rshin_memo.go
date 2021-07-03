package main

import (
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

type RshinMemo struct {
	app           *tview.Application
	dailyListView *tview.Table
}

func NewRshinMemo() *RshinMemo {
	return &RshinMemo{}
}

func (r *RshinMemo) Run() error {
	var err error
	r.app = tview.NewApplication()
	// dailyListViewの初期化
	r.dailyListView = tview.NewTable().SetSelectable(true, false)
	r.dailyListView, err = setContents(r.dailyListView)
	if err != nil {
		return err
	}

	// レイアウト
	flex := tview.NewFlex().AddItem(r.dailyListView, 300, 0, true)

	if err := r.app.SetRoot(flex, true).Run(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func setContents(table *tview.Table) (*tview.Table, error) {
	table.SetCellSimple(0, 0, "2021-01-01")
	table.SetCellSimple(0, 1, "aaaaaaaaaa")
	table.SetCellSimple(1, 0, "2021-01-01")
	table.SetCellSimple(1, 1, "aaaaaaaaaa")
	return table, nil
}

//app := tview.NewApplication()
////box := tview.NewBox().SetBorder(true).SetTitle("Hello, world!")
//textView := tview.NewTextView().SetTitle(" どうですか ").SetBorder(true)
//// dailyListを用意する。
//list := tview.NewList().ShowSecondaryText(false).
//	AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
//	AddItem("2行目のコンテンツ", "", '2', nil).
//	AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
//	AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
//	AddItem("2行目のコンテンツ", "", '2', nil).
//	AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
//	AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
//	AddItem("2行目のコンテンツ", "", '2', nil).
//	AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
//	AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
//	AddItem("2行目のコンテンツ", "", '2', nil).
//	AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
//	AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
//	AddItem("2行目のコンテンツ", "", '2', nil).
//	AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
//	AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
//	AddItem("2行目のコンテンツ", "", '2', nil).
//	AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
//	AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
//	AddItem("2行目のコンテンツ", "", '2', nil).
//	AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
//	AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
//	AddItem("2行目のコンテンツ", "", '2', nil).
//	AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
//	AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
//	AddItem("2行目のコンテンツ", "", '2', nil).
//	AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
//	AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
//	AddItem("2行目のコンテンツ", "", '2', nil).
//	AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
//	AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
//	AddItem("2行目のコンテンツ", "", '2', nil).
//	AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
//	AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
//	AddItem("2行目のコンテンツ", "", '2', nil).
//	AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
//	AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
//	AddItem("2行目のコンテンツ", "", '2', nil).
//	AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
//	AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
//	AddItem("2行目のコンテンツ", "", '2', nil).
//	AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil)
//// ボタン
//
//flex := tview.NewFlex().
//	AddItem(list, 0, 1, true)
//
////AddItem(textView, 0, 1, false)
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
