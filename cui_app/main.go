package main

import (
	"log"
)

func main() {
	rshinMemo := NewRshinMemo()
	defer rshinMemo.Close()

	err := rshinMemo.Run()
	if err != nil {
		log.Panicf("%+v", err)
	}
}
