package main

import (
	"github.com/rivo/tview"
)

func main() {
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
	treeView := tview.NewTreeView()

	flex := tview.NewFlex().
		AddItem(list, 0, 1, true).
		AddItem(textView, 0, 1, false)

	if err := tview.NewApplication().SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
