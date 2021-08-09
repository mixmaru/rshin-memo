package main

import (
	"github.com/mixmaru/rshin-memo/cmd/rshinmemo/utils"
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/mixmaru/rshin-memo/tui/rshinmemo/views"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"path/filepath"
	"time"
)

type RshinMemo struct {
	app                 *tview.Application
	memoDirPath         string // memoファイルをおいているDirPath
	layoutView          *views.LayoutView
	dailyListView       *views.dailyListView
	dailyListInsertMode usecases.InsertMode
	dateSelectView      *views.DateSelectView
	noteSelectView      *views.NoteSelectView

	dailyDataRep repositories.DailyDataRepositoryInterface
	noteRep      repositories.NoteRepositoryInterface
}

func NewRshinMemo(
	memoDirPath string,
	dailyDataRep repositories.DailyDataRepositoryInterface,
	noteRep repositories.NoteRepositoryInterface,
) *RshinMemo {
	return &RshinMemo{
		memoDirPath:  memoDirPath,
		dailyDataRep: dailyDataRep,
		noteRep:      noteRep,
	}
}

func (r *RshinMemo) Run() error {
	var err error
	r.app = tview.NewApplication()
	r.layoutView, r.dailyListView, err = r.createInitViews()
	if err != nil {
		return err
	}

	if err := r.app.SetRoot(r.layoutView.layoutView, true).Run(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (r *RshinMemo) createInitViews() (layoutView *views.layoutView, dailyListView *views.dailyListView, err error) {
	dailyListView, err = r.createDailyListView()
	if err != nil {
		return nil, nil, err
	}

	layoutView = views.newLayoutView()
	layoutView.AddPage(dailyListView)
	return layoutView, dailyListView, nil
}

func (r *RshinMemo) createDailyListView() (*views.dailyListView, error) {
	dailyList := views.NewDailyListView()
	// イベント設定
	dailyList.AddWhenPushLowerOKey(func() error {
		return r.displayDateSelectView(usecases.INSERT_UNDER_DATE_MODE)
	})
	dailyList.AddWhenPushUpperOKey(func() error {
		return r.displayDateSelectView(usecases.INSERT_OVER_DATE_MODE)
	})

	// データセット
	dailyListData, err := r.getDailyListAllData()
	if err != nil {
		return nil, err
	}
	dailyList.SetData(dailyListData)
	return dailyList, nil
}

func (r *RshinMemo) displayDateSelectView(mode usecases.InsertMode) error {
	var err error
	r.dailyListInsertMode = mode
	r.dateSelectView, err = r.createInitDailySelectView(r.dailyListInsertMode)
	if err != nil {
		panic(errors.WithStack(err))
	}
	// 表示領域に挿入する
	r.layoutView.AddPage(r.dateSelectView)
	return nil
}

func (r *RshinMemo) loadDailyListAllData() error {
	// データ取得
	dailyListData, err := r.getDailyListAllData()
	if err != nil {
		return err
	}

	// データをテーブルにセット
	r.dailyListView.SetData(dailyListData)
	return nil
}

func (r *RshinMemo) getDailyListAllData() ([]usecases.DailyData, error) {
	useCase := usecases.NewGetAllDailyListUsecase(r.dailyDataRep)
	return useCase.Handle()
}
func (r *RshinMemo) createInitDailySelectView(mode usecases.InsertMode) (*views.DateSelectView, error) {
	dateSelectView := views.NewDateSelectView()
	dateSelectView.AddWhenPushEscapeKey(func() error {
		// dateSelectViewを削除してDailyListにフォーカスを戻す
		r.closeDateSelectView()
		return nil
	})

	dateSelectView.AddWhenPushEnterKey(func() error {
		// noteSelectViewを表示してフォーカスを移す
		var err error
		r.noteSelectView, err = r.createNoteSelectView()
		if err != nil {
			return err
		}
		r.layoutView.AddPage(r.noteSelectView)
		return nil
	})

	// データセット
	dates, err := r.createDates(mode)
	if err != nil {
		return nil, err
	}
	dateSelectView.SetData(dates)
	return dateSelectView, nil
}

func (r *RshinMemo) createDates(mode usecases.InsertMode) ([]time.Time, error) {
	// 表示する日付の範囲を決定する
	overCurrentDate, err := r.getDailyListCursorDate(-1)
	if err != nil {
		return nil, err
	}
	currentDate, err := r.getDailyListCursorDate(0)
	if err != nil {
		return nil, err
	}
	underCurrentDate, err := r.getDailyListCursorDate(1)
	if err != nil {
		return nil, err
	}

	now := time.Now().In(time.Local)
	useCase := usecases.NewGetDateSelectRangeUseCase(now)
	return useCase.Handle(overCurrentDate, currentDate, underCurrentDate, mode)
}

func (r *RshinMemo) closeDateSelectView() {
	r.layoutView.RemovePage(r.dateSelectView)
	r.dateSelectView = nil
}

func (r *RshinMemo) createNoteSelectView() (*views.NoteSelectView, error) {
	noteSelectView := views.NewNoteSelectView()
	noteSelectView.AddWhenPushEscapeKey(func() error {
		r.closeNoteSelectView()
		return nil
	})
	noteSelectView.AddWhenPushEnterKey(func() error {
		// 指定のnoteをデータに追加
		err := r.saveDailyData()
		if err != nil {
			return err
		}

		// vimでひらく
		noteName := r.getNoteSelectCursorNoteName()
		err = utils.OpenVim(filepath.Join(r.memoDirPath, noteName+".txt"))
		if err != nil {
			return err
		}

		// dailyList表示までもどす
		r.closeNoteSelectView()
		r.closeDateSelectView()
		// データ再読込
		return r.loadDailyListAllData()
	})

	useCase := usecases.NewGetAllNotesUseCase(r.noteRep)
	notes, err := useCase.Handle()
	if err != nil {
		return nil, err
	}
	noteSelectView.SetData(notes)
	return noteSelectView, nil
}

func (r *RshinMemo) closeNoteSelectView() {
	r.layoutView.RemovePage(r.noteSelectView)
	r.noteSelectView = nil
}

func (r *RshinMemo) saveDailyData() error {
	// 選択note名を取得する
	noteName := r.getNoteSelectCursorNoteName()

	// 選択した日付を取得
	selectedDate, err := r.getDailyListCursorDate(0)
	if err != nil {
		return err
	}
	newDailyData, err := r.createNewDailyData(selectedDate, noteName, r.dailyListInsertMode)
	if err != nil {
		return err
	}

	usecase := usecases.NewSaveDailyDataUseCase(r.noteRep, r.dailyDataRep)
	err = usecase.Handle(newDailyData)
	if err != nil {
		return err
	}

	return nil
}

// dailyListのカーソル位置の日付を取得する。
// cursorPointAdjustに数値を指定すると、指定分カーソル位置からずれた位置の日付を取得する
func (r *RshinMemo) getDailyListCursorDate(cursorPointAdjust int) (time.Time, error) {
	row, _ := r.dailyListView.GetSelection()
	targetRow := row + cursorPointAdjust
	if targetRow < 0 || targetRow+1 > r.dailyListView.GetRowCount() {
		// targetRow is out of range
		return time.Time{}, nil
	}
	dateStr := r.dailyListView.GetCell(targetRow, 0)
	date, err := time.ParseInLocation("2006-01-02", dateStr.Text, time.Local)
	if err != nil {
		return time.Time{}, errors.WithStack(err)
	}
	return date, nil
}

// noteSelectViewのカーソル位置のnome名を取得する
func (r *RshinMemo) getNoteSelectCursorNoteName() (noteName string) {
	row, _ := r.noteSelectView.GetSelection()
	cell := r.noteSelectView.GetCell(row, 0)
	return cell.Text
}

// わたされた日付のnote一覧を取得
// 取得しながら、カーソル位置を基準にnoteを挿入する
func (r *RshinMemo) createNewDailyData(date time.Time, noteName string, mode usecases.InsertMode) (usecases.DailyData, error) {
	retData := usecases.DailyData{}
	retData.Date = date.Format("2006-01-02")
	for i := 0; i < r.dailyListView.GetRowCount(); i++ {
		// 挿入位置であれば新規noteを追加
		insertPoint, err := r.getInsertPoint(mode)
		if err != nil {
			return usecases.DailyData{}, err
		}
		if i == insertPoint {
			retData.Notes = append(retData.Notes, noteName)
		}

		tmpDateStr := r.dailyListView.GetCell(i, 0).Text
		tmpNoteName := r.dailyListView.GetCell(i, 1).Text
		if tmpDateStr == retData.Date {
			retData.Notes = append(retData.Notes, tmpNoteName)
		}
	}
	return retData, nil
}

func (r *RshinMemo) getInsertPoint(mode usecases.InsertMode) (int, error) {
	cursorRow, _ := r.dailyListView.GetSelection()
	switch mode {
	case usecases.INSERT_OVER_DATE_MODE:
		return cursorRow, nil
	case usecases.INSERT_UNDER_DATE_MODE:
		return cursorRow + 1, nil
	default:
		return 0, errors.Errorf("考慮外の値. mode: %v", mode)
	}
}

//list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
//	switch event.Key() {
//	case tcell.KeyEnter:
//		flex.AddItem(textView, 0, 1, false)
//		app.SetFocus(textView)
//		return nil
//	case tcell.KeyEsc:
//		flex.RemoveItem(textView)
//		return nil
//	}
//	return event
//})
//
//if err := app.SetRoot(flex, true).Run(); err != nil {
//	panic(err)
//}
