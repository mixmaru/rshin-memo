package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"time"
)

type DateSelectView struct {
	view                               *tview.Table
	name                               string
	whenPushEscapeKey                  []func() error
	whenPushEnterKeyOnDateLine         []func(selectedDate time.Time) error
	whenPushEnterKeyOnInputNewDateLine []func() error
}

func (d *DateSelectView) GetTviewTable() *tview.Table {
	return d.view
}

func (d *DateSelectView) GetName() string {
	return d.name
}

func NewDateSelectView() *DateSelectView {
	dateSelectView := &DateSelectView{name: "date_select_view"}
	dateSelectView.initView()
	return dateSelectView
}

func (d *DateSelectView) initView() {
	view := tview.NewTable().SetSelectable(true, false)
	// event設定
	view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			err := d.executeWhenPushEscapeKey()
			if err != nil {
				panic(err)
			}
			return nil
		case tcell.KeyEnter:
			if d.isSelectedInputNewDate() {
				err := d.executeWhenPushEnterKeyOnInputNewDateLine()
				if err != nil {
					panic(err)
				}
			} else {
				selectedDate, err := d.getSelectedDate()
				if err != nil {
					panic(err)
				}
				err = d.executeWhenPushEnterKeyOnDateLine(selectedDate)
				if err != nil {
					panic(err)
				}
			}
		}
		return event
	})
	d.view = view
}

func (d *DateSelectView) getSelectedDate() (time.Time, error) {
	row, _ := d.view.GetSelection()
	cell := d.view.GetCell(row, 0)
	dateStr := cell.Text
	date, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	if err != nil {
		return time.Time{}, errors.WithStack(err)
	}
	return date, nil
}

func (d *DateSelectView) AddWhenPushEnterKeyOnDateLine(function func(selectedDate time.Time) error) {
	d.whenPushEnterKeyOnDateLine = append(d.whenPushEnterKeyOnDateLine, function)
}

func (d *DateSelectView) AddWhenPushEnterKeyOnInputNewDateLine(function func() error) {
	d.whenPushEnterKeyOnInputNewDateLine = append(d.whenPushEnterKeyOnInputNewDateLine, function)
}

func (d *DateSelectView) AddWhenPushEscapeKey(function func() error) {
	d.whenPushEscapeKey = append(d.whenPushEscapeKey, function)
}

func (d *DateSelectView) executeWhenPushEscapeKey() error {
	return executeFunctions(d.whenPushEscapeKey)
}

func (d *DateSelectView) executeWhenPushEnterKeyOnDateLine(selectedDate time.Time) error {
	for _, function := range d.whenPushEnterKeyOnDateLine {
		err := function(selectedDate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DateSelectView) executeWhenPushEnterKeyOnInputNewDateLine() error {
	return executeFunctions(d.whenPushEnterKeyOnInputNewDateLine)
}

const NEW_INPUT_DATE_TITLE = "手入力する"

func (d *DateSelectView) SetData(dates []time.Time) {
	d.view.Clear()
	d.view.SetCellSimple(0, 0, NEW_INPUT_DATE_TITLE)
	for i, date := range dates {
		d.view.SetCellSimple(i+1, 0, date.Format("2006-01-02"))
	}
}

func (d *DateSelectView) isSelectedInputNewDate() bool {
	row, _ := d.view.GetSelection()
	return d.view.GetCell(row, 0).Text == NEW_INPUT_DATE_TITLE
}
