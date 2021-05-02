package main

import (
	"log"

	"github.com/jroimartin/gocui"
)

func main() {
	rshinMemo := NewRshinMemo()
	defer rshinMemo.Close()

	err := rshinMemo.Run()
	if err != gocui.ErrQuit {
		log.Panicf("%+v", err)
	}
}
