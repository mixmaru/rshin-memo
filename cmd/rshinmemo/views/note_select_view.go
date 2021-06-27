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
)

const NOTE_SELECT_VIEW = "note_select"
const NOTE_SELECT_SEARCH_INPUT_VIEW = "note_select_search_input_view"

type NoteSelectView struct {
	gui             *gocui.Gui
	noteSelectView  *gocui.View
	searchInputView *gocui.View
	memoDirPath     string
	notes           []string

	insertData dto.InsertData

	dailYDataRepository repositories.DailyDataRepositoryInterface
	noteRepository      repositories.NoteRepositoryInterface

	*ViewBase
}

func (n *NoteSelectView) Delete() error {
	err := n.deleteSearchView()
	if err != nil {
		return err
	}
	return deleteView(n.gui, n.viewName)
}

func (n *NoteSelectView) Focus() error {
	return focus(n.gui, n.viewName)
}

func (n *NoteSelectView) AllDelete() error {
	err := n.deleteSearchView()
	if err != nil {
		return err
	}
	return allDelete(n, n.parentView)
}

func (n *NoteSelectView) deleteSearchView() error {
	n.gui.DeleteKeybindings(NOTE_SELECT_SEARCH_INPUT_VIEW)
	err := n.gui.DeleteView(NOTE_SELECT_SEARCH_INPUT_VIEW)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (n *NoteSelectView) deleteThisView(g *gocui.Gui, v *gocui.View) error {
	return deleteThisView(n, n.parentView)
}

func (n *NoteSelectView) Resize() error {
	width, height := n.gui.Size()
	return resize(n.gui, n.viewName, width/2-25, 0, width/2+25, height-2, n.childView)
}

func NewNoteSelectView(
	gui *gocui.Gui,
	memoDirPath string,
	insertData dto.InsertData,
	parentView View,
	dailYDataRepository repositories.DailyDataRepositoryInterface,
	noteRepository repositories.NoteRepositoryInterface,
) *NoteSelectView {
	retObj := &NoteSelectView{
		gui:                 gui,
		insertData:          insertData,
		memoDirPath:         memoDirPath,
		dailYDataRepository: dailYDataRepository,
		noteRepository:      noteRepository,
	}
	retObj.ViewBase = NewViewBase(NOTE_SELECT_VIEW, gui, parentView)
	return retObj
}

// 新規作成
func (n *NoteSelectView) Create(notes []string) error {
	n.notes = notes
	width, height := n.gui.Size()
	v, err := createOrResizeView(n.gui, NOTE_SELECT_VIEW, width/2-25, 0, width/2+25, height-2)
	if err != nil {
		return err
	}
	n.noteSelectView = v

	n.noteSelectView.Highlight = true
	n.noteSelectView.SelBgColor = gocui.ColorGreen
	n.noteSelectView.SelFgColor = gocui.ColorBlack

	n.searchInputView, err = createOrResizeView(n.gui, NOTE_SELECT_SEARCH_INPUT_VIEW, 0, height-2, width, height)
	n.searchInputView.Frame = false
	n.searchInputView.Editable = true
	n.searchInputView.Editor = &Editor{}

	n.setContents()

	err = n.setEvents()
	if err != nil {
		return err
	}
	return nil
}

func (n *NoteSelectView) setEvents() error {
	if err := n.gui.SetKeybinding(NOTE_SELECT_VIEW, gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return errors.Wrap(err, "キーバインド失敗")
	}
	if err := n.gui.SetKeybinding(NOTE_SELECT_VIEW, 'j', gocui.ModNone, cursorDown); err != nil {
		return errors.Wrap(err, "キーバインド失敗")
	}
	if err := n.gui.SetKeybinding(NOTE_SELECT_VIEW, gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return errors.Wrap(err, "キーバインド失敗")
	}
	if err := n.gui.SetKeybinding(NOTE_SELECT_VIEW, 'k', gocui.ModNone, cursorUp); err != nil {
		return errors.Wrap(err, "キーバイーンド失敗")
	}
	if err := n.gui.SetKeybinding(NOTE_SELECT_VIEW, gocui.KeyEnter, gocui.ModNone, n.insertNoteToDailyList); err != nil {
		return errors.Wrap(err, "キーバイーンド失敗")
	}
	if err := n.gui.SetKeybinding(NOTE_SELECT_VIEW, gocui.KeyEsc, gocui.ModNone, n.deleteThisView); err != nil {
		return errors.Wrap(err, "キーバイーンド失敗")
	}
	if err := n.gui.SetKeybinding(NOTE_SELECT_VIEW, gocui.KeyCtrlF, gocui.ModNone, n.focusSearchView); err != nil {
		return errors.Wrap(err, "キーバイーンド失敗")
	}

	if err := n.gui.SetKeybinding(NOTE_SELECT_SEARCH_INPUT_VIEW, gocui.KeyEsc, gocui.ModNone, n.focusNoteSelectView); err != nil {
		return errors.Wrap(err, "キーバイーンド失敗")
	}
	return nil
}

func (n *NoteSelectView) insertNoteToDailyList(g *gocui.Gui, v *gocui.View) error {
	if n.isSelectedNewNote() {
		return n.addNote()
	} else {
		return n.insertExistedNoteToDailyList()
	}
}

func (n *NoteSelectView) focusSearchView(g *gocui.Gui, v *gocui.View) error {
	_, err := n.gui.SetCurrentView(NOTE_SELECT_SEARCH_INPUT_VIEW)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (n *NoteSelectView) focusNoteSelectView(g *gocui.Gui, v *gocui.View) error {
	n.searchInputView.Clear()
	err := n.searchInputView.SetCursor(0, 0)
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = n.gui.SetCurrentView(NOTE_SELECT_VIEW)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (n *NoteSelectView) addNote() error {
	noteNameInputView := NewNoteNameInputView(
		n.gui,
		n.memoDirPath,
		n.insertData,
		n,
		n.dailYDataRepository,
		n.noteRepository,
	)
	err := noteNameInputView.Create()
	if err != nil {
		return err
	}
	// フォーカスの移動
	err = noteNameInputView.Focus()
	if err != nil {
		return errors.Wrap(err, "フォーカス移動失敗")
	}
	n.childView = noteNameInputView
	return nil
}

func (n *NoteSelectView) insertExistedNoteToDailyList() error {
	// noteNameを取得
	noteName, err := n.getNoteNameOnCursor()
	if err != nil {
		return err
	}
	//r.openViews = append(r.openViews, r.noteSelectView)
	n.insertData.NoteName = noteName

	if err := n.createNewDailyList(); err != nil {
		return err
	}

	err = utils.OpenVim(filepath.Join(n.memoDirPath, noteName+".txt"))
	if err != nil {
		return err
	}

	// 不要なviewを閉じる
	err = n.AllDelete()
	if err != nil {
		return err
	}
	return nil
}

func (n *NoteSelectView) createNewDailyList() error {
	// Note作成を依頼
	dailyData, err := n.insertData.GenerateNewDailyData()
	if err != nil {
		return err
	}
	useCase := usecases.NewSaveDailyDataUseCase(n.noteRepository, n.dailYDataRepository)
	err = useCase.Handle(dailyData)
	if err != nil {
		// todo: エラーメッセージビューへメッセージを表示する
		return err
	}
	return nil
}

func (n *NoteSelectView) setContents() {
	fmt.Fprintln(n.noteSelectView, utils.ConvertStringForView("新規追加"))
	for _, note := range n.notes {
		fmt.Fprintln(n.noteSelectView, utils.ConvertStringForView(note))
	}
}

func (n *NoteSelectView) getNoteNameOnCursor() (string, error) {
	_, y := n.noteSelectView.Cursor()
	noteName, err := n.noteSelectView.Line(y)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return utils.ConvertStringForLogic(noteName), nil
}

func (n *NoteSelectView) isSelectedNewNote() bool {
	_, y := n.noteSelectView.Cursor()
	if y == 0 {
		return true
	} else {
		return false
	}
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	v.MoveCursor(0, 1, false)
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	v.MoveCursor(0, -1, false)
	return nil
}
