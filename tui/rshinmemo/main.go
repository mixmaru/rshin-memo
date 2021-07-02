package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	//box := tview.NewBox().SetBorder(true).SetTitle("Hello, world!")
	textView := tview.NewTextView().SetTitle(" どうですか ").SetBorder(true)
	// dailyListを用意する。
	list := tview.NewList().ShowSecondaryText(false).
		AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
		AddItem("2行目のコンテンツ", "", '2', nil).
		AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
		AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
		AddItem("2行目のコンテンツ", "", '2', nil).
		AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
		AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
		AddItem("2行目のコンテンツ", "", '2', nil).
		AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
		AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
		AddItem("2行目のコンテンツ", "", '2', nil).
		AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
		AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
		AddItem("2行目のコンテンツ", "", '2', nil).
		AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
		AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
		AddItem("2行目のコンテンツ", "", '2', nil).
		AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
		AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
		AddItem("2行目のコンテンツ", "", '2', nil).
		AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
		AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
		AddItem("2行目のコンテンツ", "", '2', nil).
		AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
		AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
		AddItem("2行目のコンテンツ", "", '2', nil).
		AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
		AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
		AddItem("2行目のコンテンツ", "", '2', nil).
		AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
		AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
		AddItem("2行目のコンテンツ", "", '2', nil).
		AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
		AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
		AddItem("2行目のコンテンツ", "", '2', nil).
		AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
		AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
		AddItem("2行目のコンテンツ", "", '2', nil).
		AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil).
		AddItem("1行目のコンテンツ", "セカンダリーテキスト", '1', nil).
		AddItem("2行目のコンテンツ", "", '2', nil).
		AddItem("3行目のコンテンツ", "セカンダリーテキスト3", '3', nil)
	// ボタン

	flex := tview.NewFlex().
		AddItem(list, 0, 1, true)

	//AddItem(textView, 0, 1, false)
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			flex.AddItem(textView, 0, 1, false)
			app.SetFocus(textView)
			return nil
		case tcell.KeyEsc:
			flex.RemoveItem(textView)
			return nil
		}
		return event
	})

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
