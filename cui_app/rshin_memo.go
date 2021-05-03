package main

import (
	"fmt"

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

	if err := r.initOrResizeViews(); err != nil {
		return err
	}

	if err := r.setEventActions(); err != nil {
		return err
	}
	return nil
}

const DAILY_LIST_VIEW = "daily_list"

// viewの設定
// viewの初期化とリサイズは同じ処理なので使い回す
func (r *RshinMemo) initOrResizeViews() error {
    // viewの生成
	_, height := r.gui.Size()
	v, err := r.gui.SetView(DAILY_LIST_VIEW, 0, 0, 50, height-1)
	if err != nil && err != gocui.ErrUnknownView {
		return errors.Wrapf(err, "%vの初期化またはリサイズ失敗", DAILY_LIST_VIEW)
	}
    // viewへの設定
    v.Highlight = true
    v.SelBgColor = gocui.ColorGreen
    v.SelFgColor = gocui.ColorBlack
    // 初期ダミーデータ
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

    // 起動時のフォーカス設定
    r.gui.SetCurrentView(DAILY_LIST_VIEW)
	return nil
}

// イベントに対してのアクションを設定する
func (r *RshinMemo) setEventActions() error {
	// CtrlC
	if err := r.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	// daily_listのカーソル移動
	if err := r.gui.SetKeybinding(DAILY_LIST_VIEW, gocui.KeyArrowDown, gocui.ModNone, r.cursorDown); err != nil {
		return err
	}
	return nil
}

func (r *RshinMemo) cursorDown(g *gocui.Gui, v *gocui.View) error {
    v.MoveCursor(0, 1, false)
	return nil
}

func (r *RshinMemo) Close() {
	r.gui.Close()
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
