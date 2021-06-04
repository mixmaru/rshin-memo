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
	dateSelectView     *views.DateSelectView
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
	rshinMemo.dateSelectView = views.NewDateSelectView(rshinMemo.gui)

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

	// dateSelectView
	if err := r.gui.SetKeybinding(views.DATE_SELECT_VIEW, gocui.KeyArrowDown, gocui.ModNone, r.cursorDown); err != nil {
		return errors.Wrap(err, "キーバインド失敗")
	}
	if err := r.gui.SetKeybinding(views.DATE_SELECT_VIEW, 'j', gocui.ModNone, r.cursorDown); err != nil {
		return errors.Wrap(err, "キーバインド失敗")
	}
	if err := r.gui.SetKeybinding(views.DATE_SELECT_VIEW, gocui.KeyArrowUp, gocui.ModNone, r.cursorUp); err != nil {
		return errors.Wrap(err, "キーバインド失敗")
	}
	if err := r.gui.SetKeybinding(views.DATE_SELECT_VIEW, 'k', gocui.ModNone, r.cursorUp); err != nil {
		return errors.Wrap(err, "キーバイーンド失敗")
	}
	if err := r.gui.SetKeybinding(views.DATE_SELECT_VIEW, gocui.KeyEnter, gocui.ModNone, r.decisionDate); err != nil {
		return errors.Wrap(err, "キーバイーンド失敗")
	}

	return nil
}

func (r *RshinMemo) decisionDate(g *gocui.Gui, v *gocui.View) error {
	if r.dateSelectView.IsSelectedHandInput() {
		// dateInputViewの表示
		err := r.displayDateInputView()
		if err != nil {
			return err
		}
	} else {
		var err error
		r.insertData.DateStr, err = r.dateSelectView.GetDateOnCursor()
		if err != nil {
			return err
		}
		// noteSelectViewの表示
		allNotes, err := r.getAllNotesUseCase.Handle()
		noteSelectView := views.NewNoteSelectView(
			r.gui,
			r.insertData,
			r.openViews,
			r.memoDirPath,
			r.getNoteUseCase,
			r.saveDailyDataUseCase,
		)
		err = noteSelectView.Create(allNotes)
		if err != nil {
			return err
		}
		noteSelectView.WhenFinished = func() error {
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
		err = noteSelectView.Focus()
		if err != nil {
			return err
		}
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

	dates, err := dateRange.GetSomeDateInRange(30)
	if err != nil {
		return err
	}
	err = r.dateSelectView.Create(dates)
	r.openViews = append(r.openViews, r.dateSelectView)
	if err != nil {
		return err
	}
	err = r.dateSelectView.Focus()
	if err != nil {
		return err
	}
	return nil
}

func (r *RshinMemo) displayDateInputView() error {
	// rangeを取得する
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

	// note名入力viewの表示
	dateInputView := views.NewDateInputView(
		r.gui,
		r.insertData,
		dateRange,
		r.getAllNotesUseCase,
		r.getNoteUseCase,
		r.saveDailyDataUseCase,
		r.memoDirPath,
		r.openViews,
	)
	err = dateInputView.Create()
	if err != nil {
		return err
	}
	dateInputView.WhenFinished = func() error {
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
	// フォーカスの移動
	err = dateInputView.Focus()
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

// vimで対象noteを開く
func (r *RshinMemo) openVim(noteName string) error {
	return utils.OpenVim(filepath.Join(r.memoDirPath, noteName+".txt"))
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
