package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

const VIEW_NAME = "dailyList"

type DailyListViewManager struct {
	view *gocui.View
}

func NewDailyListViewManager() *DailyListViewManager {
	return &DailyListViewManager{}
}

func (m *DailyListViewManager) Layout(g *gocui.Gui) error {
	_, hight := g.Size()
	if v, err := g.SetView(VIEW_NAME, 0, 0, 30, hight-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		m.view = v
		m.flushDailyListView()
	}
	g.SetCurrentView(VIEW_NAME)
	return nil
}

// 表示をデータ状態に合わせて描画する
func (m *DailyListViewManager) flushDailyListView() {
	fmt.Fprintln(m.view, "2021-05-02")
	fmt.Fprintln(m.view, " な ん ら か メ モ 1")
	fmt.Fprintln(m.view, " な ん ら か メ モ 2")
	fmt.Fprintln(m.view, " な ん ら か メ モ 3")
	fmt.Fprintln(m.view, "")
	fmt.Fprintln(m.view, "2021-05-01")
	fmt.Fprintln(m.view, " な ん ら か メ モ 1")
	fmt.Fprintln(m.view, " な ん ら か メ モ 2")
	fmt.Fprintln(m.view, " な ん ら か メ モ 3")
	fmt.Fprintln(m.view, "")
	fmt.Fprintln(m.view, "2021-04-30")
	fmt.Fprintln(m.view, " な ん ら か メ モ 1")
	fmt.Fprintln(m.view, " な ん ら か メ モ 2")
	fmt.Fprintln(m.view, " な ん ら か メ モ 3")
}

// func quit(g *gocui.Gui, v *gocui.View) error {
// 	return gocui.ErrQuit
// }
// 
// // カーソルを一つ下に動かす
// func down(g *gocui.Gui, v *gocui.View) error {
//     // 現在の位置を取得
//     _, y := v.Cursor()
//     v.SetCursor(0, y + 1)
//     return nil
// }
// 
// // カーソルを一つうえに動かす
// func up(g *gocui.Gui, v *gocui.View) error {
//     // 現在の位置を取得
//     _, y := v.Cursor()
//     v.SetCursor(0, y - 1)
//     return nil
// }
// 
// func enter(g *gocui.Gui, v *gocui.View) error {
//     // 現在の位置を取得
//     _, y := v.Cursor()
//     str, _ := v.Line(y)
// 	fmt.Fprintln(v, str)
//     return nil
// }
