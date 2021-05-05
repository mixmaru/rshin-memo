package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/mattn/go-runewidth"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type RshinMemo struct {
	memoDirPath string
	gui                *gocui.Gui
	alreadyInitialized bool
	getAllDailyListUsecase usecases.GetAllDailyListUsecaseInterface
}

func NewRshinMemo(
	getAllDailyListUsecase usecases.GetAllDailyListUsecaseInterface,
) *RshinMemo {
	homedir, err := os.UserHomeDir()
	if err != nil{
		log.Panicf("初期化失敗. %+v", err)
	}
	rshinMemo := &RshinMemo{}
	rshinMemo.memoDirPath = filepath.Join(homedir, "rshin_memo")
	rshinMemo.alreadyInitialized = false
	rshinMemo.getAllDailyListUsecase = getAllDailyListUsecase
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

type dailyData struct {
	Date  string
	Notes []string
}

func (r *RshinMemo) initViews() error {
	_, err := r.createDailyListView()
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
			_, err = fmt.Fprintln(v, dailyData.Date+"\t"+convertStringForView(note))
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
	_, height := r.gui.Size()
	_, err := r.createOrResizeView(DAILY_LIST_VIEW, 0, 0, 50, height-1)
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

	// daily_listでの新規list追加
	if err := r.gui.SetKeybinding(DAILY_LIST_VIEW, 'o', gocui.ModNone, r.addList); err != nil {
		return errors.Wrap(err, "Enterキーバインド失敗")
	}
	return nil
}

func (r *RshinMemo) addList(g *gocui.Gui, v *gocui.View) error {
	// note名入力viewの表示
	_, err := r.createNoteNameInputView()
	// フォーカスの移動
	_, err = r.gui.SetCurrentView(NOTE_NAME_INPUT_VIEW)
	if err != nil {
		return errors.Wrap(err, "フォーカス移動失敗")
	}
	return nil
}

const NOTE_NAME_INPUT_VIEW = "note_name_input"
func (r *RshinMemo) createNoteNameInputView() (*gocui.View, error) {
	width, height := r.gui.Size()
	_, err := r.createOrResizeView(NOTE_NAME_INPUT_VIEW, width/2-20, height/2-1, width/2+20, height/2+1)
	if err != nil{
		return nil, err
	}
	return nil, nil
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
