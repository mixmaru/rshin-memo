package main

import (
	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"
)

type RshinMemo struct {
	gui                *gocui.Gui
	alreadyInitialized bool
}

func NewRshinMemo() *RshinMemo {
	rshinMemo := &RshinMemo{}
	rshinMemo.alreadyInitialized = false
	return rshinMemo
}

func (r *RshinMemo) Run() error {
	// guiの初期化
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	g.SetManagerFunc(r.layout)
	r.gui = g

	// guiメインループの起動
	if err := r.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}

// layout is called for every screen re-render e.g. when the screen is resized
func (r *RshinMemo) layout(g *gocui.Gui) error {
	if !r.alreadyInitialized {
		// 初期化
		if err := r.init(); err != nil {
			return err
		}
		r.alreadyInitialized = true
	} else {
		// viewのリサイズ
		if err := r.initOrResizeViews(); err != nil {
			return err
		}
	}
	return nil
}

func (r *RshinMemo) init() error {
	// 画面の設定
	r.gui.Cursor = true

	// viewの設定
	if err := r.initOrResizeViews(); err != nil {
		return err
	}

	// キーバインディング設定
	if err := r.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}

const DAILY_LIST_VIEW = "daily_list"

// viewの初期化とリサイズは同じ処理なので使い回す
func (r *RshinMemo) initOrResizeViews() error {
	_, height := r.gui.Size()
	_, err := r.gui.SetView(DAILY_LIST_VIEW, 0, 0, 50, height-1)
	if err != nil && err != gocui.ErrUnknownView {
		return errors.Wrapf(err, "%vの初期化またはリサイズ失敗", DAILY_LIST_VIEW)
	}
	return nil
}

func (r *RshinMemo) Close() {
	r.gui.Close()
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
