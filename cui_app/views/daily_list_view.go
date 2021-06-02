package views

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/mixmaru/rshin-memo/cui_app/utils"
	"github.com/pkg/errors"
	"strings"
	"time"
)

const DAILY_LIST_VIEW = "daily_list"

type DailyListView struct {
	gui                    *gocui.Gui
	getAllDailyListUsecase *usecases.GetAllDailyListUsecase
	view                   *gocui.View
	dailyList              []usecases.DailyData
}

func NewDailyListView(gui *gocui.Gui, getAllDailyListUsecase *usecases.GetAllDailyListUsecase) *DailyListView {
	retObj := &DailyListView{
		gui:                    gui,
		getAllDailyListUsecase: getAllDailyListUsecase,
	}
	return retObj
}

// dailyListViewの新規作成
func (d *DailyListView) Create() error {
	// あとでどうせリサイズされるので、ここではこまかな位置調整は行わない。
	v, err := createOrResizeView(d.gui, DAILY_LIST_VIEW, 0, 0, 1, 1)
	if err != nil {
		return err
	}
	d.view = v

	// viewへの設定
	d.view.Highlight = true
	d.view.SelBgColor = gocui.ColorGreen
	d.view.SelFgColor = gocui.ColorBlack

	err = d.Reload()
	if err != nil {
		return err
	}

	return nil
}

func (d *DailyListView) Resize() error {
	_, height := d.gui.Size()
	_, err := createOrResizeView(d.gui, DAILY_LIST_VIEW, 0, 0, 50, height-1)
	if err != nil {
		return err
	}
	return nil
}

func (d *DailyListView) Focus() error {
	_, err := d.gui.SetCurrentView(DAILY_LIST_VIEW)
	if err != nil {
		return errors.Wrap(err, "フォーカス失敗")
	}
	return nil
}

func createOrResizeView(gui *gocui.Gui, viewName string, x0, y0, x1, y1 int) (*gocui.View, error) {
	v, err := gui.SetView(viewName, x0, y0, x1, y1)
	if err != nil && err != gocui.ErrUnknownView {
		return nil, errors.Wrapf(err, "%vの初期化またはリサイズ失敗", DAILY_LIST_VIEW)
	}
	return v, nil
}

func (d *DailyListView) loadAllDailyList() ([]usecases.DailyData, error) {
	return d.getAllDailyListUsecase.Handle()
}

func (d *DailyListView) GetDateOnCursor() (string, error) {
	_, y := d.view.Cursor()
	text, err := d.view.Line(y)
	if err != nil {
		return "", errors.Wrap(err, "選択行のtextの取得に失敗")
	}
	return getDateString(text), nil
}

func (d *DailyListView) GetDateOnCursorNext() (string, error) {
	_, y := d.view.Cursor()
	text, err := d.view.Line(y + 1)
	if err != nil {
		return "", errors.Wrap(err, "選択行の次の行のtextの取得に失敗")
	}
	return getDateString(text), nil
}

func (d *DailyListView) GetDateOnCursorPrev() (string, error) {
	_, y := d.view.Cursor()
	text, err := d.view.Line(y - 1)
	if err != nil {
		return "", errors.Wrap(err, "選択行の前の行のtextの取得に失敗")
	}
	return getDateString(text), nil
}

func getDateString(text string) string {
	// \tで分割して日付を取得
	return strings.Split(text, "\t")[0]
}

func getNoteString(text string) string {
	// \tで分割してNote名を取得
	return strings.Split(text, "\t")[1]
}

const (
	prev_cursor = iota
	next_cursor
)

func (d *DailyListView) Reload() error {
	d.view.Clear()
	d.dailyList = nil

	dailyList, err := d.loadAllDailyList()
	if err != nil {
		return err
	}
	d.dailyList = dailyList

	for _, dailyData := range d.dailyList {
		for _, note := range dailyData.Notes {
			_, err = fmt.Fprintln(d.view, dailyData.Date+"\t"+utils.ConvertStringForView(note))
			if err != nil {
				return errors.Wrapf(err, "テキスト出力失敗。%+v", dailyData)
			}
		}
	}
	return nil
}

func (d *DailyListView) GetInsertDateRangeNextCursor() (DateRange, error) {
	retDateRange := DateRange{}
	// カーソル位置の日付を取得する
	toDateString, err := d.GetDateOnCursor()
	if err != nil {
		return DateRange{}, err
	}
	err = retDateRange.SetToByString(toDateString)
	if err != nil {
		return DateRange{}, err
	}

	// if カーソルがデータの末でなければ一つ次の日付を取得する
	_, y := d.view.Cursor()
	line, err := d.view.Line(y)
	if err != nil {
		return DateRange{}, errors.WithStack(err)
	}
	date := getDateString(line)
	note := getNoteString(utils.ConvertStringForLogic(line))
	if !isLastNote(d.dailyList, date, note) {
		fromDateString, err := d.GetDateOnCursorNext()
		if err != nil {
			return DateRange{}, err
		}
		err = retDateRange.SetFromByString(fromDateString)
		if err != nil {
			return DateRange{}, err
		}
	}
	return retDateRange, nil
}

func (d *DailyListView) GetInsertDateRangePrevCursor() (DateRange, error) {
	retDateRange := DateRange{}
	// カーソル位置の日付を取得する
	toDateString, err := d.GetDateOnCursor()
	if err != nil {
		return DateRange{}, err
	}
	err = retDateRange.SetFromByString(toDateString)
	if err != nil {
		return DateRange{}, err
	}

	// if カーソルがデータの先頭でなければ一つ前の日付を取得する
	if _, y := d.view.Cursor(); y > 0 {
		fromDateString, err := d.GetDateOnCursorPrev()
		if err != nil {
			return DateRange{}, err
		}
		err = retDateRange.SetToByString(fromDateString)
		if err != nil {
			return DateRange{}, err
		}
	}
	return retDateRange, nil
}

func (d *DailyListView) GetDailyDataByDate(dateStr string) usecases.DailyData {
	for _, dailyData := range d.dailyList {
		if dailyData.Date == dateStr {
			return dailyData
		}
	}
	return usecases.DailyData{}
}

func (d *DailyListView) OnCursorRowPosition() (int, error) {
	_, y := d.view.Cursor()
	lineStr, err := d.view.Line(y)
	if err != nil {
		return 0, err
	}
	selectedDateStr := getDateString(utils.ConvertStringForLogic(lineStr))
	selectedNoteName := getNoteString(utils.ConvertStringForLogic(lineStr))
	rowPosition := 0
	for _, dailyData := range d.dailyList {
		if dailyData.Date == selectedDateStr {
			for _, noteName := range dailyData.Notes {
				rowPosition += 1
				if noteName == selectedNoteName {
					return rowPosition, nil
				}
			}
		}
		rowPosition += len(dailyData.Notes)
	}
	return 0, errors.New("カーソル上のNoteNameが見当たらない")
}

func (d *DailyListView) GetDailyList() []usecases.DailyData {
	return d.dailyList
}

func isLastNote(dailyList []usecases.DailyData, date, note string) bool {
	lastDailyData := dailyList[len(dailyList)-1]
	if lastDailyData.Date == date {
		lastNote := lastDailyData.Notes[len(lastDailyData.Notes)-1]
		if lastNote == note {
			return true
		}
	}
	return false
}

type DateRange struct {
	From time.Time
	To   time.Time
}

// 日付の範囲から最大指定個数の日をスライスで返す
func (d *DateRange) GetSomeDateInRange(num int) ([]time.Time, error) {
	retDates := []time.Time{}
	var startData time.Time
	if !d.From.IsZero() {
		startData = d.From
	} else if !d.To.IsZero() {
		startData = d.To.AddDate(0, 0, -(num - 1))
	} else {
		return nil, errors.Errorf("このDateRangeはFromもToも設定されていませんので実行できません。%+v", d)
	}

	for i := 0; i < num; i++ {
		date := startData.AddDate(0, 0, i)
		if !d.To.IsZero() && date.After(d.To) {
			break
		}
		retDates = append(retDates, date)
	}
	return retDates, nil
}

func (d *DateRange) IsIn(targetDate time.Time) bool {
	if !d.From.IsZero() {
		if targetDate.Before(d.From) {
			return false
		}
	}
	if !d.To.IsZero() {
		if targetDate.After(d.To) {
			return false
		}
	}
	return true
}

func (d *DateRange) SetFromByString(dateStr string) error {
	var err error
	d.From, err = time.Parse("2006-01-02", dateStr)
	if err != nil {
		return errors.WithStack(err)
	}
	return err
}

func (d *DateRange) SetToByString(dateStr string) error {
	var err error
	d.To, err = time.Parse("2006-01-02", dateStr)
	if err != nil {
		return errors.WithStack(err)
	}
	return err
}
