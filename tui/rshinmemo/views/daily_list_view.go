package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"time"
)

type DailyListView struct {
	view              *tview.Table
	name              string
	whenPushEnterKey  []func() error
	whenPushLowerOkey []func() error
	whenPushUpperOkey []func() error
}

func (d *DailyListView) GetTviewPrimitive() tview.Primitive {
	return d.view
}

func (d *DailyListView) GetName() string {
	return d.name
}

func NewDailyListView() *DailyListView {
	dailyListView := &DailyListView{name: "DailyListView"}
	dailyListView.initView()
	return dailyListView
}

func (d *DailyListView) initView() {
	table := tview.NewTable()
	table.SetSelectable(true, false)
	// イベント設定
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			err := d.whenPushEnterkey()
			if err != nil {
				panic(err)
			}
		case tcell.KeyRune:
			switch event.Rune() {
			case 'o':
				err := d.whenPushLowerOKey()
				if err != nil {
					panic(err)
				}
			case 'O':
				err := d.whenPushUpperOKey()
				if err != nil {
					panic(err)
				}
			}
		}
		return event
	})
	d.view = table

	return
}

func (d *DailyListView) AddWhenPushEnterKey(function func() error) {
	d.whenPushEnterKey = append(d.whenPushEnterKey, function)
}

func (d *DailyListView) AddWhenPushLowerOKey(function func() error) {
	d.whenPushLowerOkey = append(d.whenPushLowerOkey, function)
}

func (d *DailyListView) AddWhenPushUpperOKey(function func() error) {
	d.whenPushUpperOkey = append(d.whenPushUpperOkey, function)
}

func (d *DailyListView) whenPushEnterkey() error {
	return executeFunctions(d.whenPushEnterKey)
}

func (d *DailyListView) whenPushLowerOKey() error {
	return executeFunctions(d.whenPushLowerOkey)
}

func (d *DailyListView) whenPushUpperOKey() error {
	return executeFunctions(d.whenPushUpperOkey)
}

func (d *DailyListView) SetData(data []usecases.DailyData) {
	d.view.Clear()
	// データをテーブルにセット
	row := 0
	for _, data := range data {
		for _, note := range data.Notes {
			d.view.SetCellSimple(row, 0, data.Date)
			d.view.SetCellSimple(row, 1, note)
			row++
		}
	}
}

// GetCursorDate dailyListのカーソル位置の日付を取得する。
// cursorPointAdjustに数値を指定すると、指定分カーソル位置からずれた位置の日付を取得する
func (d *DailyListView) GetCursorDate(cursorPointAdjust int) (time.Time, error) {
	row, _ := d.view.GetSelection()
	targetRow := row + cursorPointAdjust
	if targetRow < 0 || targetRow+1 > d.view.GetRowCount() {
		// targetRow is out of range
		return time.Time{}, nil
	}
	dateStr := d.view.GetCell(targetRow, 0)
	date, err := time.ParseInLocation("2006-01-02", dateStr.Text, time.Local)
	if err != nil {
		return time.Time{}, errors.WithStack(err)
	}
	return date, nil
}

func (d *DailyListView) GetCursorNoteName() string {
	row, _ := d.view.GetSelection()
	note := d.view.GetCell(row, 1)
	return note.Text
}

func (d *DailyListView) GetInsertPoint(mode usecases.InsertMode) (int, error) {
	cursorRow, _ := d.view.GetSelection()
	switch mode {
	case usecases.INSERT_OVER_DATE_MODE:
		return cursorRow, nil
	case usecases.INSERT_UNDER_DATE_MODE:
		return cursorRow + 1, nil
	default:
		return 0, errors.Errorf("考慮外の値. mode: %v", mode)
	}
}

func (d *DailyListView) GetRowCount() int {
	return d.view.GetRowCount()
}

func (d *DailyListView) GetCell(row int, i int) *tview.TableCell {
	return d.view.GetCell(row, i)

}

func executeFunctions(functions []func() error) error {
	for _, function := range functions {
		err := function()
		if err != nil {
			return err
		}
	}
	return nil
}
