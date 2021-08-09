package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/rivo/tview"
)

type dailyListView struct {
	view              *tview.Table
	name              string
	whenPushLowerOkey []func() error
	whenPushUpperOkey []func() error
}

func (d *dailyListView) GetTviewTable() *tview.Table {
	return d.view
}

func (d *dailyListView) GetName() string {
	return d.name
}

func NewDailyListView() *dailyListView {
	dailyListView := &dailyListView{name: "dailyListView"}
	dailyListView.initView()
	return dailyListView
}

func (d *dailyListView) initView() {
	table := tview.NewTable()
	table.SetSelectable(true, false)
	// イベント設定
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
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

func (d *dailyListView) AddWhenPushLowerOKey(function func() error) {
	d.whenPushLowerOkey = append(d.whenPushLowerOkey, function)
}

func (d *dailyListView) AddWhenPushUpperOKey(function func() error) {
	d.whenPushUpperOkey = append(d.whenPushUpperOkey, function)
}

func (d *dailyListView) whenPushLowerOKey() error {
	return executeFunctions(d.whenPushLowerOkey)
}

func (d *dailyListView) whenPushUpperOKey() error {
	return executeFunctions(d.whenPushUpperOkey)
}

func (d *dailyListView) SetData(data []usecases.DailyData) {
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

func (d *dailyListView) GetSelection() (row, column int) {
	return d.view.GetSelection()
}

func (d *dailyListView) GetRowCount() int {
	return d.view.GetRowCount()
}

func (d *dailyListView) GetCell(row int, i int) *tview.TableCell {
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
