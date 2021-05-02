package main

import (
	"log"

	"github.com/jroimartin/gocui"
)

type RshinMemo struct {
	gui *gocui.Gui
}

func NewRshinMemo() *RshinMemo {
	rshinMemo := &RshinMemo{}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

    rshinMemo.gui = g
	rshinMemo.init()
	return rshinMemo
}

func (r *RshinMemo) init() {
	// 画面の設定
	r.gui.Cursor = true

	// viewの設定
	r.gui.SetManager(NewDailyListViewManager())

	// キーバインディング設定
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

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
