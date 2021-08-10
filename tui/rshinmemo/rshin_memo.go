package main

import (
	"github.com/mixmaru/rshin-memo/cmd/rshinmemo/utils"
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/mixmaru/rshin-memo/tui/rshinmemo/views"
	"github.com/pkg/errors"
	"path/filepath"
	"time"
)

type RshinMemo struct {
	memoDirPath         string // memoファイルをおいているDirPath
	layoutView          *views.LayoutView
	dailyListView       *views.DailyListView
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
	r.layoutView, r.dailyListView, err = r.createInitViews()
	if err != nil {
		return err
	}

	r.layoutView.SetRoot(r.dailyListView)
	return r.layoutView.Run()
}

func (r *RshinMemo) createInitViews() (layoutView *views.LayoutView, dailyListView *views.DailyListView, err error) {
	dailyListView, err = r.createDailyListView()
	if err != nil {
		return nil, nil, err
	}

	layoutView = views.NewLayoutView()
	layoutView.AddPage(dailyListView)
	return layoutView, dailyListView, nil
}

func (r *RshinMemo) createDailyListView() (*views.DailyListView, error) {
	dailyList := views.NewDailyListView()
	// イベント設定
	dailyList.AddWhenPushEnterKey(func() error {
		noteName := r.dailyListView.GetCursorNoteName()
		return r.openVim(noteName)
	})
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
	overCurrentDate, err := r.dailyListView.GetCursorDate(-1)
	if err != nil {
		return nil, err
	}
	currentDate, err := r.dailyListView.GetCursorDate(0)
	if err != nil {
		return nil, err
	}
	underCurrentDate, err := r.dailyListView.GetCursorDate(1)
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

		noteName := r.getNoteSelectCursorNoteName()
		err = r.openVim(noteName)
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

// vimでひらく
func (r *RshinMemo) openVim(noteName string) error {
	err := utils.OpenVim(filepath.Join(r.memoDirPath, noteName+".txt"))
	if err != nil {
		return nil
	}
	r.layoutView.Refresh()
	return nil
}

func (r *RshinMemo) closeNoteSelectView() {
	r.layoutView.RemovePage(r.noteSelectView)
	r.noteSelectView = nil
}

func (r *RshinMemo) saveDailyData() error {
	// 選択note名を取得する
	noteName := r.getNoteSelectCursorNoteName()

	// 選択した日付を取得
	selectedDate, err := r.dailyListView.GetCursorDate(0)
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
	insertPoint, err := r.dailyListView.GetInsertPoint(mode)
	rowCount := r.dailyListView.GetRowCount()
	for i := 0; i < rowCount; i++ {
		// 挿入位置であれば新規noteを追加
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
	// insert to end of list
	if insertPoint == rowCount {
		retData.Notes = append(retData.Notes, noteName)
	}
	return retData, nil
}
