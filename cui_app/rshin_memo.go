package main

import (
	"github.com/jroimartin/gocui"
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/mixmaru/rshin-memo/cui_app/utils"
	"github.com/mixmaru/rshin-memo/cui_app/views"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type RshinMemo struct {
	memoDirPath        string
	gui                *gocui.Gui
	dailyListView      *views.DailyListView
	noteNameInputView  *views.NoteNameInputView
	dateInputView      *views.DateInputView
	noteSelectView     *views.NoteSelectView
	alreadyInitialized bool

	getNoteUseCase       *usecases.GetNoteUseCase
	getAllNoteUseCase    *usecases.GetAllNoteUseCase
	saveDailyDataUseCase *usecases.SaveDailyDataUseCase

	addRowMode AddRowMode
}

type AddRowMode int

const (
	ADD_ROW_PREV_MODE = iota
	ADD_ROW_NEXT_MODE
)

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
	rshinMemo.dailyListView = views.NewDailyListView(rshinMemo.gui, usecases.NewGetAllDailyListUsecase(dailyDataRepository))
	rshinMemo.noteNameInputView = views.NewNoteNameinputView(rshinMemo.gui)
	rshinMemo.dateInputView = views.NewDateInputView(rshinMemo.gui)
	rshinMemo.noteSelectView = views.NewNoteSelectView(rshinMemo.gui)
	rshinMemo.getNoteUseCase = usecases.NewGetNoteUseCase(noteRepository)
	rshinMemo.saveDailyDataUseCase = usecases.NewSaveDailyDataUseCase(noteRepository, dailyDataRepository)
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

	// daily_listでカーソルの下行に新規list追加
	if err := r.gui.SetKeybinding(views.DAILY_LIST_VIEW, 'o', gocui.ModNone, r.displayDateInputViewForNext); err != nil {
		return errors.Wrap(err, "oキーバインド失敗")
	}
	// daily_listでカーソルの上行に新規list追加
	if err := r.gui.SetKeybinding(views.DAILY_LIST_VIEW, 'O', gocui.ModNone, r.displayDataInputViewForPrev); err != nil {
		return errors.Wrap(err, "Oキーバインド失敗")
	}
	// DateInputViewでのEnterキー
	if err := r.gui.SetKeybinding(views.DATE_INPUT_VIEW, gocui.KeyEnter, gocui.ModNone, r.displayNoteNameInputView); err != nil {
		return errors.Wrap(err, "Enterキーバインド失敗")
	}

	// inputNoteNameViewでのEnterキー
	if err := r.gui.SetKeybinding(views.NOTE_NAME_INPUT_VIEW, gocui.KeyEnter, gocui.ModNone, r.createNote); err != nil {
		return errors.Wrap(err, "Enterキーバインド失敗")
	}

	return nil
}
func (r *RshinMemo) displayDateInputViewForNext(g *gocui.Gui, v *gocui.View) error {
	// inputViewを表示する
	r.addRowMode = ADD_ROW_NEXT_MODE
	return r.displayDateInputView()
}

func (r *RshinMemo) displayDataInputViewForPrev(g *gocui.Gui, v *gocui.View) error {
	r.addRowMode = ADD_ROW_PREV_MODE
	return r.displayDateInputView()
}

func (r *RshinMemo) displayDateInputView() error {
	// note名入力viewの表示
	err := r.dateInputView.Create()
	if err != nil {
		return err
	}
	// フォーカスの移動
	err = r.dateInputView.Focus()
	if err != nil {
		return errors.Wrap(err, "フォーカス移動失敗")
	}
	return nil
}

func (r *RshinMemo) displayNoteNameInputView(g *gocui.Gui, v *gocui.View) error {
	// 日付入力値の取得
	dateString, err := r.dateInputView.GetInputString()
	// 日付入力値のバリデーション
	result, err := r.valid(dateString)
	if err != nil {
		return err
	}
	if !result {
		return nil
	}

	// noteSelectViewの表示
	allNotes, err := r.getAllNoteUseCase.Handle()
	err = r.noteSelectView.Create(allNotes)
	if err != nil {
		return err
	}
	err = r.noteSelectView.Focus()
	if err != nil {
		return err
	}

	// note名入力viewの表示
	//err = r.noteNameInputView.Create()
	//if err != nil {
	//	return err
	//}
	//// フォーカスの移動
	//err = r.noteNameInputView.Focus()
	//if err != nil {
	//	return errors.Wrap(err, "フォーカス移動失敗")
	//}
	return nil
}

func (r *RshinMemo) addNote() error {
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
	if err != nil {
		return errors.Wrap(err, "選択行のtextの取得に失敗")
	}
	// \tで分割してノート名を取得
	noteName := strings.Split(text, "\t")[1]
	// 取得したテキストは表示のために半角スペースがはいってるので除去
	noteName = utils.ConvertStringForLogic(noteName)

	err = r.openVim(noteName)
	if err != nil {
		return err
	}
	return nil
}

func (r *RshinMemo) createNote(gui *gocui.Gui, view *gocui.View) error {
	// 入力内容を取得
	date, err := r.dateInputView.GetInputString()
	if err != nil {
		return err
	}
	noteName, err := r.noteNameInputView.GetInputNoteName()
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	// 同名Noteが存在しないかcheck
	_, notExist, err := r.getNoteUseCase.Handle(noteName)
	if err != nil {
		return err
	} else if !notExist {
		// すでに同名のNoteが存在する
		// todo: エラーメッセージビューへメッセージを表示する
	} else {
		// 対象日付のdailyListを取得作成
		var dailyData usecases.DailyData
		if r.addRowMode == ADD_ROW_PREV_MODE {
			dailyData, err = r.dailyListView.GenerateNewDailyDataToPrevCursor(noteName, date)
			if err != nil {
				return err
			}
		} else {
			dailyData, err = r.dailyListView.GenerateNewDailyDataToNextCursor(noteName, date)
			if err != nil {
				return err
			}
		}
		// Note作成を依頼
		err = r.saveDailyDataUseCase.Handle(dailyData)
		if err != nil {
			// todo: エラーメッセージビューへメッセージを表示する
			return err
		}

		err = r.openVim(noteName)
		if err != nil {
			return err
		}

		// 追加されたNoteが表示されるようにDailyListをリフレッシュ
		err = r.dailyListView.Reload()
		if err != nil {
			return err
		}
	}

	err = r.dailyListView.Focus()
	if err != nil {
		return err
	}

	err = r.noteNameInputView.Delete()
	if err != nil {
		return err
	}
	return nil
}

func (r *RshinMemo) valid(dateString string) (bool, error) {
	// rangeを取得する
	var dateRange views.DateRange
	var err error
	switch r.addRowMode {
	case ADD_ROW_PREV_MODE:
		dateRange, err = r.dailyListView.GetInsertDateRangePrevCursor()
		if err != nil {
			return false, err
		}
	case ADD_ROW_NEXT_MODE:
		dateRange, err = r.dailyListView.GetInsertDateRangeNextCursor()
		if err != nil {
			return false, err
		}
	default:
		return false, errors.Errorf("考慮外の値が使われた。addRowMode: %v", r.addRowMode)
	}

	// 指定のdate文字列がRangeの範囲にとどまっているかをチェック

	targetDate, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		// パース失敗（入力フォーマットが違う）
		return false, nil
	}
	if !dateRange.IsIn(targetDate) {
		return false, nil
	}
	return true, nil
}

// vimで対象noteを開く
func (r *RshinMemo) openVim(noteName string) error {
	c := exec.Command("vim", filepath.Join(r.memoDirPath, noteName+".txt"))
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		return errors.Wrap(err, "vim起動エラー")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
