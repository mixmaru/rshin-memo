package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/mattn/go-runewidth"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/pkg/errors"
	"time"
)

/*
vimを開く方法メモ
c := exec.Command("vim", "main.go")
c.Stdin = os.Stdin
c.Stdout = os.Stdout
c.Stderr = os.Stderr
err := c.Run()
if err != nil {
    return errors.Wrap(err, "実行エラー")
}
*/

type RshinMemo struct {
	gui                *gocui.Gui
	alreadyInitialized bool
	getAllDailyListUsecase usecases.GetAllDailyListUsecaseInterface
}

func NewRshinMemo(getAllDailyListUsecase usecases.GetAllDailyListUsecaseInterface) *RshinMemo {
	rshinMemo := &RshinMemo{}
	rshinMemo.alreadyInitialized = false
	rshinMemo.getAllDailyListUsecase = getAllDailyListUsecase
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
		if err := r.resizeViews(); err != nil {
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

const DAILY_LIST_VIEW = "daily_list"
const NOTE_VIEW = "note"

type dailyData struct {
	Date  time.Time
	Notes []string
}

func (r *RshinMemo) initViews() error {
	_, err := r.createDailyListView()
	if err != nil{
		return err
	}

	_, err = r.createNoteView()
	if err != nil{
		return err
	}

	// 起動時のフォーカス設定
	_, err = r.gui.SetCurrentView(DAILY_LIST_VIEW)
	if err != nil {
		return errors.Wrap(err, "起動時フォーカス失敗")
	}
	return nil
}

func (r * RshinMemo) createDailyListView() (*gocui.View, error) {
	// あとでどうせリサイズされるので、ここではこまかな位置調整は行わない。
	v, err := r.createOrResizeView(DAILY_LIST_VIEW, 0, 0, 1, 1)
	if err != nil {
		return nil, err
	}
	// viewへの設定
	v.Highlight = true
	v.SelBgColor = gocui.ColorGreen
	v.SelFgColor = gocui.ColorBlack

	dailyList, err := r.loadAllDailyList()
	if err != nil {
		return nil, err
	}

	for _, dailyData := range dailyList {
		for _, note := range dailyData.Notes {
			_, err = fmt.Fprintln(v, dailyData.Date.Format("2006-01-02")+"\t"+convertStringForView(note))
			if err != nil {
				return nil, errors.Wrapf(err, "テキスト出力失敗。%+v", dailyData)
			}
		}
	}
	return v, nil
}

func (r * RshinMemo) loadAllDailyList() ([]dailyData, error) {
	retList := []dailyData{}
	response, err := r.getAllDailyListUsecase.Handle()
	if err != nil{
		return nil, err
	}
	for _, oneDayList := range response.DailyList {
		dailyData := dailyData{
			Date: oneDayList.Date,
			Notes: oneDayList.Notes,
		}
		retList = append(retList, dailyData)
	}
	return retList, nil
}

func (r * RshinMemo) createNoteView() (*gocui.View, error) {
	// あとでどうせリサイズされて配置調整されるので、ここでは細かな位置調整は行わない
	v, err := r.createOrResizeView(NOTE_VIEW, 0, 0, 1, 1)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func convertStringForView(s string) string {
	runeArr := []rune{}
	for _, r := range s {
		runeArr = append(runeArr, r)
		// if もし全角文字だったら
		if runewidth.StringWidth(string(r)) == 2 {
			runeArr = append(runeArr, ' ')
		}
	}
	return string(runeArr)
}

func (r *RshinMemo) createOrResizeView(viewName string, x0, y0, x1, y1 int) (*gocui.View, error) {
	v, err := r.gui.SetView(viewName, x0, y0, x1, y1)
	if err != nil && err != gocui.ErrUnknownView {
		return nil, errors.Wrapf(err, "%vの初期化またはリサイズ失敗", DAILY_LIST_VIEW)
	}
	return v, nil
}

// viewのリサイズ
func (r *RshinMemo) resizeViews() error {
	width, height := r.gui.Size()
	dailyListView, err := r.createOrResizeView(DAILY_LIST_VIEW, 0, 0, 50, height-1)
	if err != nil {
		return err
	}

	x, _ := dailyListView.Size()
	_, err = r.createOrResizeView(NOTE_VIEW, x+2, 0, width-1, height-1)
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

	// daily_listのカーソル移動
	if err := r.gui.SetKeybinding(DAILY_LIST_VIEW, gocui.KeyArrowDown, gocui.ModNone, r.cursorDown); err != nil {
		return errors.Wrap(err, "KeyArrowDownキーバインド失敗")
	}
	if err := r.gui.SetKeybinding(DAILY_LIST_VIEW, 'j', gocui.ModNone, r.cursorDown); err != nil {
		return errors.Wrap(err, "jキーバインド失敗")
	}
	if err := r.gui.SetKeybinding(DAILY_LIST_VIEW, gocui.KeyArrowUp, gocui.ModNone, r.cursorUp); err != nil {
		return errors.Wrap(err, "KeyArrowUpキーバインド失敗")
	}
	if err := r.gui.SetKeybinding(DAILY_LIST_VIEW, 'k', gocui.ModNone, r.cursorUp); err != nil {
		return errors.Wrap(err, "kキーバイーンド失敗")
	}

	// daily_listでのエンターキー
	if err := r.gui.SetKeybinding(DAILY_LIST_VIEW, gocui.KeyEnter, gocui.ModNone, r.openNote); err != nil {
		return errors.Wrap(err, "Enterキーバインド失敗")
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

// ノートViewに指定のノートを表示する
func (r *RshinMemo) openNote(g *gocui.Gui, v *gocui.View) error {
	// 選択行のテキストを取得
	// \tで分割してノート名を取得
	// ノート名で清書データを取得
	// データをnoteViewへ流し込む
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
