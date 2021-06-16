package views

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/mixmaru/rshin-memo/cmd/rshinmemo/dto"
	"github.com/mixmaru/rshin-memo/cmd/rshinmemo/utils"
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/pkg/errors"
	"path/filepath"
	"strings"
	"time"
)

const DAILY_LIST_VIEW = "daily_list"

type DailyListView struct {
	gui       *gocui.Gui
	view      *gocui.View
	dailyList []usecases.DailyData

	memoDirPath string

	dailyDataRepository repositories.DailyDataRepositoryInterface
	noteRepository      repositories.NoteRepositoryInterface

	explainView *ExplainView
	*ViewBase
}

func NewDailyListView(
	gui *gocui.Gui,
	memoDirPath string,
	dailyDataRepository repositories.DailyDataRepositoryInterface,
	noteRepository repositories.NoteRepositoryInterface,
	explainView *ExplainView,
) *DailyListView {
	retObj := &DailyListView{
		gui:                 gui,
		memoDirPath:         memoDirPath,
		dailyDataRepository: dailyDataRepository,
		noteRepository:      noteRepository,
		explainView:         explainView,
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

	err = d.setEvents()
	if err != nil {
		return err
	}

	d.ViewBase = NewViewBase(DAILY_LIST_VIEW, d.gui, []View{d})

	d.explainView.Set("[j]:up cursor [k]:down cursor [o]:insert memo under the cursor [O]:insert memo just the cursor [enter]:open memo")
	return nil
}

func (d *DailyListView) setEvents() error {
	// daily_listのカーソル移動
	if err := d.gui.SetKeybinding(DAILY_LIST_VIEW, gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return errors.Wrap(err, "キーバインド失敗")
	}
	if err := d.gui.SetKeybinding(DAILY_LIST_VIEW, 'j', gocui.ModNone, cursorDown); err != nil {
		return errors.Wrap(err, "キーバインド失敗")
	}
	if err := d.gui.SetKeybinding(DAILY_LIST_VIEW, gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return errors.Wrap(err, "キーバインド失敗")
	}
	if err := d.gui.SetKeybinding(DAILY_LIST_VIEW, 'k', gocui.ModNone, cursorUp); err != nil {
		return errors.Wrap(err, "キーバイーンド失敗")
	}
	// daily_listでのエンターキー
	if err := d.gui.SetKeybinding(DAILY_LIST_VIEW, gocui.KeyEnter, gocui.ModNone, d.openNote); err != nil {
		return errors.Wrap(err, "キーバインド失敗")
	}
	// daily_listでカーソルの下行に新規list追加
	if err := d.gui.SetKeybinding(DAILY_LIST_VIEW, 'o', gocui.ModNone, d.displayDateInputViewForInsertingNextCursorPosition); err != nil {
		return errors.Wrap(err, "キーバインド失敗")
	}
	// daily_listでカーソルの上行に新規list追加
	if err := d.gui.SetKeybinding(DAILY_LIST_VIEW, 'O', gocui.ModNone, d.displayDataInputViewForInsertingCurrentCursorPosition); err != nil {
		return errors.Wrap(err, "キーバインド失敗")
	}

	return nil
}

// 指定NoteをVimで起動する
func (d *DailyListView) openNote(g *gocui.Gui, v *gocui.View) error {
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

	err = d.openVim(noteName)
	if err != nil {
		return err
	}
	return nil
}

// vimで対象noteを開く
func (d *DailyListView) openVim(noteName string) error {
	return utils.OpenVim(filepath.Join(d.memoDirPath, noteName+".txt"))
}

func (d *DailyListView) Resize() error {
	_, height := d.gui.Size()
	_, err := createOrResizeView(d.gui, DAILY_LIST_VIEW, 0, 0, 50, height-3)
	if err != nil {
		return err
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
	useCase := usecases.NewGetAllDailyListUsecase(d.dailyDataRepository)
	return useCase.Handle()
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

func (d *DailyListView) GetInsertDateRangeCurrentCursor() (DateRange, error) {
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

func (d *DailyListView) OnCursorRowPosition() (lineNum int, emptyLine bool, err error) {
	_, y := d.view.Cursor()
	lineStr, err := d.view.Line(y)
	if err != nil {
		return 0, false, err
	}
	if lineStr == "" {
		return 0, true, nil
	}
	selectedDateStr := getDateString(utils.ConvertStringForLogic(lineStr))
	selectedNoteName := getNoteString(utils.ConvertStringForLogic(lineStr))
	rowPosition := 0
	for _, dailyData := range d.dailyList {
		if dailyData.Date == selectedDateStr {
			for _, noteName := range dailyData.Notes {
				if noteName == selectedNoteName {
					return rowPosition, false, nil
				}
				rowPosition += 1
			}
		}
		rowPosition += len(dailyData.Notes)
	}
	return 0, false, errors.New("カーソル上のNoteNameが見当たらない")
}

func (d *DailyListView) displayDateInputViewForInsertingNextCursorPosition(g *gocui.Gui, v *gocui.View) error {
	return d.displayDataInputView(true)
}

func (d *DailyListView) displayDataInputViewForInsertingCurrentCursorPosition(g *gocui.Gui, v *gocui.View) error {
	return d.displayDataInputView(false)
}

func (d *DailyListView) displayDataInputView(next bool) error {
	var dateRange DateRange
	insertData := dto.InsertData{}
	if len(d.dailyList) == 0 {
		// データがないときは無条件に新規登録
		insertData.InsertNum = 0
		now := time.Now().In(time.Local)
		dateRange = DateRange{
			From: time.Date(now.Year(), now.Month(), now.Day()-15, 0, 0, 0, 0, time.Local),
			To:   time.Date(now.Year(), now.Month(), now.Day()+15, 0, 0, 0, 0, time.Local),
		}
	} else {
		insertNum, emptyLine, err := d.OnCursorRowPosition()
		if err != nil {
			return err
		}
		if emptyLine {
			// データがあるのに何もない行でoかOを押したときはなにもしない
			return nil
		}

		if next {
			insertData.InsertNum = insertNum + 1
			dateRange, err = d.GetInsertDateRangeNextCursor()
			if err != nil {
				return err
			}
		} else {
			insertData.InsertNum = insertNum
			dateRange, err = d.GetInsertDateRangeCurrentCursor()
			if err != nil {
				return err
			}
		}
	}
	return d.displayDateSelectView(insertData, dateRange)
}

func (d *DailyListView) displayDateSelectView(insertData dto.InsertData, dateRange DateRange) error {
	insertData.TargetDailyData = d.dailyList
	dateSelectView := NewDateSelectView(
		d.gui,
		d.memoDirPath,
		insertData,
		dateRange,
		[]View{d},
		d.dailyDataRepository,
		d.noteRepository,
	)
	err := dateSelectView.Create()
	if err != nil {
		return err
	}
	err = dateSelectView.Focus()
	if err != nil {
		return err
	}
	return nil
}

func (d *DailyListView) Focus() error {
	err := d.Reload()
	if err != nil {
		return err
	}
	err = d.ViewBase.Focus()
	if err != nil {
		return err
	}
	return nil
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
