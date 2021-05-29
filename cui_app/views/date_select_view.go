package views

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/mixmaru/rshin-memo/cui_app/utils"
	"github.com/pkg/errors"
	"sort"
	"time"
)

const DATE_SELECT_VIEW = "date_select"

type DateSelectView struct {
	gui   *gocui.Gui
	view  *gocui.View
	dates []time.Time
}

func NewDateSelectView(gui *gocui.Gui) *DateSelectView {
	retObj := &DateSelectView{
		gui: gui,
	}
	return retObj
}

// 新規作成
func (n *DateSelectView) Create(dates []time.Time) error {
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].After(dates[j])
	})
	n.dates = dates
	width, height := n.gui.Size()
	v, err := createOrResizeView(n.gui, DATE_SELECT_VIEW, width/2-25, 0, width/2+25, height-1)
	if err != nil {
		return err
	}
	n.view = v

	n.view.Highlight = true
	n.view.SelBgColor = gocui.ColorGreen
	n.view.SelFgColor = gocui.ColorBlack

	n.setContents()

	return nil
}

func (n *DateSelectView) setContents() {
	fmt.Fprintln(n.view, utils.ConvertStringForView("手入力する"))
	for _, date := range n.dates {
		fmt.Fprintln(n.view, utils.ConvertStringForView(date.Format("2006-01-02")))
	}
}

func (n *DateSelectView) Focus() error {
	_, err := n.gui.SetCurrentView(DATE_SELECT_VIEW)
	if err != nil {
		return errors.Wrap(err, "フォーカス移動失敗")
	}
	return nil
}

func (n *DateSelectView) GetDateOnCursor() string {
	_, y := n.view.Cursor()
	return n.dates[y-1].Format("2006-01-02")
}

func (n *DateSelectView) Delete() error {
	err := n.gui.DeleteView(DATE_SELECT_VIEW)
	if err != nil {
		return errors.Wrapf(err, "Viewの削除に失敗。%v", NOTE_NAME_INPUT_VIEW)
	}
	return nil
}

func (n *DateSelectView) IsSelectedOtherDate() bool {
	_, y := n.view.Cursor()
	if y == 0 {
		return true
	} else {
		return false
	}
}
