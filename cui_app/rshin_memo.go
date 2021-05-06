package main

import (
	"github.com/jroimartin/gocui"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/mixmaru/rshin-memo/cui_app/views"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type RshinMemo struct {
	memoDirPath        string
	gui                *gocui.Gui
	dailyListView      *views.DailyListView
	noteNameInputView  *views.NoteNameInputView
	alreadyInitialized bool
}

func NewRshinMemo(
	getAllDailyListUsecase usecases.GetAllDailyListUsecaseInterface,
) *RshinMemo {

	homedir, err := os.UserHomeDir()
	if err != nil{
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
	rshinMemo.dailyListView = views.NewDailyListView(rshinMemo.gui, getAllDailyListUsecase)
	rshinMemo.noteNameInputView = views.NewNoteNameinputView(rshinMemo.gui)
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
	if err != nil{
		return err
	}

	// 起動時のフォーカス設定
	err = r.dailyListView.Focus()
	if err != nil {
		return err
	}
	return nil
}

func (r *RshinMemo) createOrResizeView(viewName string, x0, y0, x1, y1 int) (*gocui.View, error) {
	v, err := r.gui.SetView(viewName, x0, y0, x1, y1)
	if err != nil && err != gocui.ErrUnknownView {
		return nil, errors.Wrap(err, "初期化またはリサイズ失敗")
	}
	return v, nil
}

// イベントに対してのアクションを設定する
func (r *RshinMemo) setEventActions() error {
	// CtrlC
	if err := r.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return errors.Wrap(err, "CtrlCキーバインド失敗")
	}

	// daily_listのカーソル移動
	if err := r.gui.SetKeybinding(views.DAILY_LIST_VIEW, gocui.KeyArrowDown, gocui.ModNone, r.cursorDown); err != nil {
		return errors.Wrap(err, "KeyArrowDownキーバインド失敗")
	}
	if err := r.gui.SetKeybinding(views.DAILY_LIST_VIEW, 'j', gocui.ModNone, r.cursorDown); err != nil {
		return errors.Wrap(err, "jキーバインド失敗")
	}
	if err := r.gui.SetKeybinding(views.DAILY_LIST_VIEW, gocui.KeyArrowUp, gocui.ModNone, r.cursorUp); err != nil {
		return errors.Wrap(err, "KeyArrowUpキーバインド失敗")
	}
	if err := r.gui.SetKeybinding(views.DAILY_LIST_VIEW, 'k', gocui.ModNone, r.cursorUp); err != nil {
		return errors.Wrap(err, "kキーバイーンド失敗")
	}

	// daily_listでのエンターキー
	if err := r.gui.SetKeybinding(views.DAILY_LIST_VIEW, gocui.KeyEnter, gocui.ModNone, r.openNote); err != nil {
		return errors.Wrap(err, "Enterキーバインド失敗")
	}

	// daily_listでの新規list追加
	if err := r.gui.SetKeybinding(views.DAILY_LIST_VIEW, 'o', gocui.ModNone, r.addList); err != nil {
		return errors.Wrap(err, "Enterキーバインド失敗")
	}
	return nil
}

func (r *RshinMemo) addList(g *gocui.Gui, v *gocui.View) error {
	// note名入力viewの表示
	err := r.noteNameInputView.Create()
	if err != nil {
		return err
	}
	// フォーカスの移動
	err = r.noteNameInputView.Focus()
	if err != nil {
		return errors.Wrap(err, "フォーカス移動失敗")
	}
	return nil
}

func (r *RshinMemo) cursorDown(g *gocui.Gui, v *gocui.View) error {
	v.MoveCursor(0, 1, false)
	return nil
}

func (r *RshinMemo) cursorUp(g *gocui.Gui, v *gocui.View) error {
	v.MoveCursor(0, -1, false)
	return nil
}

func (r *RshinMemo) Close() {
	r.gui.Close()
}

// 指定NoteをVimで起動する
func (r *RshinMemo) openNote(g *gocui.Gui, v *gocui.View) error {
	// 選択行のテキストを取得
	_, y := v.Cursor()
	text, err := v.Line(y)
	if err != nil{
		return errors.Wrap(err, "選択行のtextの取得に失敗")
	}
	// \tで分割してノート名を取得
	noteName := strings.Split(text, "\t")[1]
	// 取得したテキストは表示のために半角スペースがはいってるので除去
	noteName = strings.ReplaceAll(noteName, " ", "")

	// vimで対象noteを開く
	c := exec.Command("vim", filepath.Join(r.memoDirPath, noteName+".txt"))
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err = c.Run()
	if err != nil {
		return errors.Wrap(err, "vim起動エラー")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
