package main

import (
	"github.com/jroimartin/gocui"
	"github.com/mixmaru/rshin-memo/cmd/rshinmemo/views"
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/pkg/errors"
	"log"
	"os"
	"path/filepath"
)

type RshinMemo struct {
	memoDirPath        string
	gui                *gocui.Gui
	dailyListView      *views.DailyListView
	alreadyInitialized bool
}

func NewRshinMemo(
	dailyDataRepository repositories.DailyDataRepositoryInterface,
	noteRepository repositories.NoteRepositoryInterface,
) *RshinMemo {

	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Panicf("初期化失敗. %+v", err)
	}

	rshinMemo := &RshinMemo{}
	// guiの初期化
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicf("初期化失敗。error: %+v", errors.Wrap(err, "初期化失敗"))
	}
	g.SetManagerFunc(rshinMemo.layout)
	rshinMemo.gui = g
	rshinMemo.memoDirPath = filepath.Join(homedir, "rshin_memo")
	rshinMemo.alreadyInitialized = false
	rshinMemo.dailyListView = views.NewDailyListView(
		rshinMemo.gui,
		rshinMemo.memoDirPath,
		dailyDataRepository,
		noteRepository,
	)

	return rshinMemo
}

func (r *RshinMemo) Run() error {
	// なければmemo用dirの作成
	if _, err := os.Stat(r.memoDirPath); os.IsNotExist(err) {
		err := os.Mkdir(r.memoDirPath, 0777)
		if err != nil {
			return errors.Wrap(err, "memo用dirの作成に失敗しました。")
		}
	}

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
		if err := r.dailyListView.Resize(); err != nil {
			return err
		}
	}
	return nil
}

func (r *RshinMemo) init() error {
	// 画面の設定
	r.gui.Cursor = true

	err := r.initViews()
	if err != nil {
		return err
	}

	if err := r.setEventActions(); err != nil {
		return err
	}
	return nil
}

func (r *RshinMemo) initViews() error {
	err := r.dailyListView.Create()
	if err != nil {
		return err
	}

	// 起動時のフォーカス設定
	err = r.dailyListView.Focus()
	if err != nil {
		return err
	}
	return nil
}

// イベントに対してのアクションを設定する
func (r *RshinMemo) setEventActions() error {
	// CtrlC
	if err := r.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return errors.Wrap(err, "CtrlCキーバインド失敗")
	}
	return nil
}

func (r *RshinMemo) Close() {
	r.gui.Close()
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
