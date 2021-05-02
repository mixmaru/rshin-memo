package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

type RshinMemo struct {
	gui *gocui.Gui
}

func NewRshinMemo() *RshinMemo {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	rshinMemo := &RshinMemo{gui: g}
	rshinMemo.init()
	return rshinMemo
}

func (r *RshinMemo) init() {
	// 画面の設定
	r.gui.Cursor = true

	// viewの設定
	r.gui.SetManagerFunc(dailyListViewLayout)

	// キーバインディングの初期化
	if err := r.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
}

func (r *RshinMemo) Run() {
	if err := r.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (r *RshinMemo) Close() {
	r.gui.Close()
}

// 表示をデータ状態に合わせて再描画する
func (r *RshinMemo) FlushView() {

}

// daily_listのview描画設定
func dailyListViewLayout(g *gocui.Gui) error {
	_, hight := g.Size()
	if v, err := g.SetView("dailyList", 0, 0, 30, hight-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "2021-05-02")
		fmt.Fprintln(v, " な ん ら か メ モ 1")
		fmt.Fprintln(v, " な ん ら か メ モ 2")
		fmt.Fprintln(v, " な ん ら か メ モ 3")
		fmt.Fprintln(v, "")
		fmt.Fprintln(v, "2021-05-01")
		fmt.Fprintln(v, " な ん ら か メ モ 1")
		fmt.Fprintln(v, " な ん ら か メ モ 2")
		fmt.Fprintln(v, " な ん ら か メ モ 3")
		fmt.Fprintln(v, "")
		fmt.Fprintln(v, "2021-04-30")
		fmt.Fprintln(v, " な ん ら か メ モ 1")
		fmt.Fprintln(v, " な ん ら か メ モ 2")
		fmt.Fprintln(v, " な ん ら か メ モ 3")
	}
	g.SetCurrentView("dailyList")
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
