package views

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/mixmaru/rshin-memo/cui_app/utils"
	"github.com/pkg/errors"
)

const DAILY_LIST_VIEW = "daily_list"

type DailyListView struct {
	gui *gocui.Gui
	getAllDailyListUsecase usecases.GetAllDailyListUsecaseInterface
	view *gocui.View
}

func NewDailyListView(gui *gocui.Gui, getAllDailyListUsecase usecases.GetAllDailyListUsecaseInterface) *DailyListView {
	retObj := &DailyListView{
		gui: gui,
		getAllDailyListUsecase: getAllDailyListUsecase,
	}
	return retObj
}

// dailyListViewの新規作成
func (d *DailyListView) Create() error{
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

	dailyList, err := d.loadAllDailyList()
	if err != nil {
		return err
	}

	for _, dailyData := range dailyList {
		for _, note := range dailyData.Notes {
			_, err = fmt.Fprintln(d.view, dailyData.Date+"\t"+utils.ConvertStringForView(note))
			if err != nil {
				return errors.Wrapf(err, "テキスト出力失敗。%+v", dailyData)
			}
		}
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

func (d *DailyListView) loadAllDailyList() ([]DailyData, error) {
	retList := []DailyData{}
	response, err := d.getAllDailyListUsecase.Handle()
	if err != nil{
		return nil, err
	}
	for _, oneDayList := range response.DailyList {
		dailyData := DailyData{
			Date: oneDayList.Date,
			Notes: oneDayList.Notes,
		}
		retList = append(retList, dailyData)
	}
	return retList, nil
}

type DailyData struct {
	Date  string
	Notes []string
}