package main

func main() {
    rshinMemo := NewRshinMemo()
    defer rshinMemo.Close()

    rshinMemo.Run()

	// g, err := gocui.NewGui(gocui.OutputNormal)
	// if err != nil {
	// 	log.Panicln(err)
	// }
	// defer g.Close()

    // g.Cursor = true

	// g.SetManagerFunc(dailyListViewLayout)

	// if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
	// 	log.Panicln(err)
	// }

	// if err := g.SetKeybinding("hello", gocui.KeyArrowDown, gocui.ModNone, down); err != nil {
	// 	log.Panicln(err)
	// }

	// if err := g.SetKeybinding("hello", gocui.KeyArrowUp, gocui.ModNone, up); err != nil {
	// 	log.Panicln(err)
	// }

	// if err := g.SetKeybinding("hello", gocui.KeyEnter, gocui.ModNone, enter); err != nil {
	// 	log.Panicln(err)
	// }

	// if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
	// 	log.Panicln(err)
	// }
}

// func quit(g *gocui.Gui, v *gocui.View) error {
// 	return gocui.ErrQuit
// }
// 
// // カーソルを一つ下に動かす
// func down(g *gocui.Gui, v *gocui.View) error {
//     // 現在の位置を取得
//     _, y := v.Cursor()
//     v.SetCursor(0, y + 1)
//     return nil
// }
// 
// // カーソルを一つうえに動かす
// func up(g *gocui.Gui, v *gocui.View) error {
//     // 現在の位置を取得
//     _, y := v.Cursor()
//     v.SetCursor(0, y - 1)
//     return nil
// }
// 
// func enter(g *gocui.Gui, v *gocui.View) error {
//     // 現在の位置を取得
//     _, y := v.Cursor()
//     str, _ := v.Line(y)
// 	fmt.Fprintln(v, str)
//     return nil
// }
