package views

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/mixmaru/rshin-memo/cui_app/utils"
	"github.com/pkg/errors"
	"strings"
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
	// \tで分割してノート名を取得
	dateText := strings.Split(text, "\t")[0]
	return dateText, nil
}

func (d *DailyListView) GenerateNewDailyDataToNextCursor(newNoteName string) (usecases.DailyData, error) {
	return d.generateNewDailyData(newNoteName, next_cursor)
}

func (d *DailyListView) GenerateNewDailyDataToPrevCursor(newNoteName string) (usecases.DailyData, error) {
	return d.generateNewDailyData(newNoteName, prev_cursor)
}

const (
	prev_cursor = iota
	next_cursor
)

func (d *DailyListView) generateNewDailyData(newNoteName string, insertPlace int) (usecases.DailyData, error) {
	// カーソル位置の日付を取得する
	date, err := d.GetDateOnCursor()
	if err != nil {
		return usecases.DailyData{}, err
	}
	// dateでdailyDateを取得して返す
	count := 0
	for _, dailyDate := range d.dailyList {
		if dailyDate.Date == date {
			// カーソル位置の下にnewNoteNameを追加する
			var insertNum int
			_, cursorNum := d.view.Cursor()
			switch insertPlace {
			case prev_cursor:
				insertNum = cursorNum - count - 1
			case next_cursor:
				insertNum = cursorNum - count
			default:
				return usecases.DailyData{}, errors.Errorf("想定外の値が使われた。insertPlace: %v", insertPlace)
			}
			newNotes := []string{}
			newNotes = append(newNotes, dailyDate.Notes[:insertNum+1]...)
			newNotes = append(newNotes, newNoteName)
			newNotes = append(newNotes, dailyDate.Notes[insertNum+1:]...)
			dailyDate.Notes = newNotes
			return dailyDate, nil
		}
		count += len(dailyDate.Notes)
	}
	return usecases.DailyData{}, errors.New("カーソル位置の日付のdailydataが取得できなかった。想定外エラー")

}

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
// numは0始まりでカウント
func IsEndOfDateList(num int, dailyList []usecases.DailyData) bool {
	num++ // 比較簡略化のため1追加しておく
	for _, dailyData := range dailyList {
		if len(dailyData.Notes) == num {
			return true
		} else if len(dailyData.Notes) < num {
			num -= len(dailyData.Notes)
		} else {
			return false
		}
	}
	return false
}

