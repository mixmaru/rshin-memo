package main

func main() {
	rshinMemo := NewRshinMemo()
	err := rshinMemo.Run()
	if err != nil {
		panic(err)
	}
}
