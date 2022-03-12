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
	dateInputView       *views.DateInputView
	noteSelectView      *views.NoteSelectView
	noteNameInputView   *views.NoteNameInputView

	dailyDataRep repositories.DailyDataRepositoryInterface
	noteRep      repositories.NoteRepositoryInterface

	selectedDate time.Time
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
		return r.displayDateSelectView(usecases.INSERT_OLDER_MODE)
	})
	dailyList.AddWhenPushUpperOKey(func() error {
		return r.displayDateSelectView(usecases.INSERT_NEWER_MODE)
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
	r.dateSelectView, err = r.createDateSelectView(r.dailyListInsertMode)
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
func (r *RshinMemo) createDateSelectView(mode usecases.InsertMode) (*views.DateSelectView, error) {
	dateSelectView := views.NewDateSelectView()
	dateSelectView.AddWhenPushEscapeKey(func() error {
		// dateSelectViewを削除してDailyListにフォーカスを戻す
		r.closeDateSelectView()
		return nil
	})

	dateSelectView.AddWhenPushEnterKeyOnDateLine(func(selectedDate time.Time) error {
		r.selectedDate = selectedDate
		// noteSelectViewを表示してフォーカスを移す
		var err error
		r.noteSelectView, err = r.createNoteSelectView()
		if err != nil {
			return err
		}
		r.layoutView.AddPage(r.noteSelectView)
		return nil
	})

	dateSelectView.AddWhenPushEnterKeyOnInputNewDateLine(func() error {
		// 日付入力viewを表示する
		r.dateInputView = r.createDateInputView()
		r.layoutView.AddPage(r.dateInputView)
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

func (r *RshinMemo) closeAllView() {
	if r.noteNameInputView != nil {
		r.closeNoteNameInputView()
	}
	if r.noteSelectView != nil {
		r.closeNoteSelectView()
	}
	if r.dateInputView != nil {
		r.closeDateInputView()
	}
	if r.dateSelectView != nil {
		r.closeDateSelectView()
	}
}

func (r *RshinMemo) closeNoteSelectView() {
	r.layoutView.RemovePage(r.noteSelectView)
	r.noteSelectView = nil
}

func (r *RshinMemo) closeDateInputView() {
	r.layoutView.RemovePage(r.dateInputView)
	r.dateInputView = nil
}

func (r *RshinMemo) closeDateSelectView() {
	r.layoutView.RemovePage(r.dateSelectView)
	r.dateSelectView = nil
}

func (r *RshinMemo) closeNoteNameInputView() {
	r.layoutView.RemovePage(r.noteNameInputView)
	r.noteNameInputView = nil
}

func (r *RshinMemo) createNoteSelectView() (*views.NoteSelectView, error) {
	noteSelectView := views.NewNoteSelectView()
	noteSelectView.AddWhenPushEscapeKey(func() error {
		r.closeNoteSelectView()
		return nil
	})
	noteSelectView.AddWhenPushEnterKeyOnNoteNameLine(func(noteName string) error {
		return r.createAndOpenNoteAndClose(noteName)
	})

	noteSelectView.AddWhenPushEnterKeyOnInputNewNoteNameLine(func() error {
		// NoteNameInputViewを表示する
		r.noteNameInputView = r.createNoteNameInputView()
		r.layoutView.AddPage(r.noteNameInputView)
		return nil
	})

	noteSelectView.AddWhenPushCtrlFKey(func() error {
		r.noteSelectView.SearchMode(r.layoutView)
		return nil
	})

	noteSelectView.AddWhenPushEnterKeyOnSearchInputField(func(searchWord string) error {
		searchUseCase := usecases.NewGetNotesBySearchTextUseCase(r.noteRep)
		notes, err := searchUseCase.Handle(searchWord)
		if err != nil {
			return err
		}
		noteSelectView.SetData(notes)
		noteSelectView.NoteSelectMode(r.layoutView)
		return nil
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

func (r *RshinMemo) saveDailyData(noteName string) error {
	newDailyData, err := r.createNewDailyData(r.selectedDate, noteName, r.dailyListInsertMode)
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

func (r *RshinMemo) createDateInputView() *views.DateInputView {
	dateInputView := views.NewDateInputView()
	dateInputView.AddWhenPushEscapeKey(func() error {
		r.closeDateInputView()
		return nil
	})
	dateInputView.AddWhenPushEnterKey(func(inputDateStr string) (*views.ValidationError, error) {
		// 入力値をパースしてtime型で渡す
		date, err := time.ParseInLocation("2006-01-02", inputDateStr, time.Local)
		if err != nil {
			return &views.ValidationError{No: 1, Message: views.VALIDATION_ERROR_1_MESSAGE}, nil
		}
		// dateが許容範囲内かチェック
		result, err := r.IsDateInRange(date)
		if err != nil {
			return nil, err
		}
		if !result {
			return &views.ValidationError{No: 2, Message: views.VALIDATION_ERROR_2_MESSAGE}, nil
		}
		r.selectedDate = date
		r.noteSelectView, err = r.createNoteSelectView()
		if err != nil {
			return nil, err
		}
		r.layoutView.AddPage(r.noteSelectView)
		return nil, nil
	})
	return dateInputView
}

// IsDateInRange dateが想定の範囲内かどうかチェックする
func (r *RshinMemo) IsDateInRange(date time.Time) (bool, error) {
	var from, to time.Time
	var err error
	switch r.dailyListInsertMode {
	case usecases.INSERT_NEWER_MODE:
		from, err = r.dailyListView.GetCursorDate(0)
		if err != nil {
			return false, err
		}
		to, err = r.dailyListView.GetCursorDate(-1)
		if err != nil {
			return false, err
		}
	case usecases.INSERT_OLDER_MODE:
		from, err = r.dailyListView.GetCursorDate(1)
		if err != nil {
			return false, err
		}
		to, err = r.dailyListView.GetCursorDate(0)
		if err != nil {
			return false, err
		}
	default:
		return false, errors.Errorf("想定外エラー。r.dailyListInsertMode: %v", r.dailyListView)
	}
	return IsDateInRange(date, from, to), nil
}

func (r *RshinMemo) createNoteNameInputView() *views.NoteNameInputView {
	view := views.NewNoteNameInputView()
	view.AddWhenPushEscapeKey(func() error {
		r.closeNoteNameInputView()
		return nil
	})
	view.AddWhenPushEnterKey(func(noteName string) error {
		return r.createAndOpenNoteAndClose(noteName)
	})
	return view
}

func (r *RshinMemo) createAndOpenNoteAndClose(noteName string) error {
	// 指定のnoteをデータに追加
	err := r.saveDailyData(noteName)
	if err != nil {
		return err
	}

	err = r.openVim(noteName)
	if err != nil {
		return err
	}
	// dailyList表示までもどす
	r.closeAllView()
	// データ再読込
	return r.loadDailyListAllData()
}

func IsDateInRange(date, from, to time.Time) bool {
	if !from.IsZero() {
		if !(from.Before(date) || from.Equal(date)) {
			return false
		}
	}
	if !to.IsZero() {
		if !(date.Before(to) || date.Equal(to)) {
			return false
		}
	}
	return true
}
