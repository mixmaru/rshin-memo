package main

import (
	"github.com/jroimartin/gocui"
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/mixmaru/rshin-memo/cui_app/dto"
	"github.com/mixmaru/rshin-memo/cui_app/utils"
	"github.com/mixmaru/rshin-memo/cui_app/views"
	"github.com/pkg/errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type RshinMemo struct {
	memoDirPath        string
	gui                *gocui.Gui
	dailyListView      *views.DailyListView
	alreadyInitialized bool

	getNoteUseCase       *usecases.GetNoteUseCase
	getAllNotesUseCase   *usecases.GetAllNotesUseCase
	saveDailyDataUseCase *usecases.SaveDailyDataUseCase

	addRowMode views.AddRowMode

	selectedDate string
	insertData   dto.InsertData
	openViews    []views.Deletable
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
	rshinMemo.dailyListView = views.NewDailyListView(rshinMemo.gui, usecases.NewGetAllDailyListUsecase(dailyDataRepository))

	rshinMemo.getNoteUseCase = usecases.NewGetNoteUseCase(noteRepository)
	rshinMemo.getAllNotesUseCase = usecases.NewGetAllNotesUseCase(noteRepository)
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
		return errors.Wrap(err, "キーバインド失敗")
	}
	if err := r.gui.SetKeybinding(views.DAILY_LIST_VIEW, 'j', gocui.ModNone, r.cursorDown); err != nil {
		return errors.Wrap(err, "キーバインド失敗")
	}
	if err := r.gui.SetKeybinding(views.DAILY_LIST_VIEW, gocui.KeyArrowUp, gocui.ModNone, r.cursorUp); err != nil {
		return errors.Wrap(err, "キーバインド失敗")
	}
	if err := r.gui.SetKeybinding(views.DAILY_LIST_VIEW, 'k', gocui.ModNone, r.cursorUp); err != nil {
		return errors.Wrap(err, "キーバイーンド失敗")
	}
	// daily_listでのエンターキー
	if err := r.gui.SetKeybinding(views.DAILY_LIST_VIEW, gocui.KeyEnter, gocui.ModNone, r.openNote); err != nil {
		return errors.Wrap(err, "キーバインド失敗")
	}
	// daily_listでカーソルの下行に新規list追加
	if err := r.gui.SetKeybinding(views.DAILY_LIST_VIEW, 'o', gocui.ModNone, r.displayDateInputViewForNext); err != nil {
		return errors.Wrap(err, "キーバインド失敗")
	}
	// daily_listでカーソルの上行に新規list追加
	if err := r.gui.SetKeybinding(views.DAILY_LIST_VIEW, 'O', gocui.ModNone, r.displayDataInputViewForPrev); err != nil {
		return errors.Wrap(err, "キーバインド失敗")
	}

	return nil
}

func (r *RshinMemo) displayDateInputViewForNext(g *gocui.Gui, v *gocui.View) error {
	r.insertData = dto.InsertData{}
	insertNum, err := r.dailyListView.OnCursorRowPosition()
	if err != nil {
		return err
	}
	r.insertData.InsertNum = insertNum + 1
	r.addRowMode = views.ADD_ROW_NEXT_MODE
	return r.displayDateSelectView()
}

func (r *RshinMemo) displayDataInputViewForPrev(g *gocui.Gui, v *gocui.View) error {
	r.insertData = dto.InsertData{}
	insertNum, err := r.dailyListView.OnCursorRowPosition()
	if err != nil {
		return err
	}
	r.insertData.InsertNum = insertNum
	r.addRowMode = views.ADD_ROW_PREV_MODE
	return r.displayDateSelectView()
}

func (r *RshinMemo) displayDateSelectView() error {
	r.insertData.TargetDailyData = r.dailyListView.GetDailyList()

	var dateRange views.DateRange
	var err error
	switch r.addRowMode {
	case views.ADD_ROW_PREV_MODE:
		dateRange, err = r.dailyListView.GetInsertDateRangePrevCursor()
		if err != nil {
			return err
		}
	case views.ADD_ROW_NEXT_MODE:
		dateRange, err = r.dailyListView.GetInsertDateRangeNextCursor()
		if err != nil {
			return err
		}
	default:
		return errors.Errorf("考慮外の値が使われた。addRowMode: %v", r.addRowMode)
	}

	dateSelectView := views.NewDateSelectView(
		r.gui,
		r.openViews,
		r.insertData,
		dateRange,
		r.memoDirPath,
		r.getAllNotesUseCase,
		r.getNoteUseCase,
		r.saveDailyDataUseCase,
	)
	err = dateSelectView.Create()
	if err != nil {
		return err
	}
	dateSelectView.WhenFinished = func() error {
		err := r.dailyListView.Reload()
		if err != nil {
			return err
		}
		err = r.dailyListView.Focus()
		if err != nil {
			return err
		}
		return nil
	}
	err = dateSelectView.Focus()
	if err != nil {
		return err
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

// vimで対象noteを開く
func (r *RshinMemo) openVim(noteName string) error {
	return utils.OpenVim(filepath.Join(r.memoDirPath, noteName+".txt"))
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
