package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

    g.Cursor = true

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("hello", gocui.KeyArrowDown, gocui.ModNone, down); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("hello", gocui.KeyArrowUp, gocui.ModNone, up); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("hello", gocui.KeyEnter, gocui.ModNone, enter); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	// maxX, maxY := g.Size()
	// if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
	if v, err := g.SetView("hello", 0, 0, 10, 10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "1111")
		fmt.Fprintln(v, "2222")
		fmt.Fprintln(v, "3333")
		fmt.Fprintln(v, "4444")
		fmt.Fprintln(v, "5555")
		fmt.Fprintln(v, "6666")
		// str, _ := v.Line(1)
		// fmt.Fprintln(v, str)
        // v.SetCursor(0, 2)
        // _, y := v.Cursor()
        // selectedStr, _ := v.Line(y)
		// fmt.Fprintln(v, selectedStr)
	}
    g.SetCurrentView("hello")
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// カーソルを一つ下に動かす
func down(g *gocui.Gui, v *gocui.View) error {
    // 現在の位置を取得
    _, y := v.Cursor()
    v.SetCursor(0, y + 1)
    return nil
}

// カーソルを一つうえに動かす
func up(g *gocui.Gui, v *gocui.View) error {
    // 現在の位置を取得
    _, y := v.Cursor()
    v.SetCursor(0, y - 1)
    return nil
}

func enter(g *gocui.Gui, v *gocui.View) error {
    // 現在の位置を取得
    _, y := v.Cursor()
    str, _ := v.Line(y)
	fmt.Fprintln(v, str)
    return nil
}
