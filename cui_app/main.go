package main

func main() {
    rshinMemo := NewRshinMemo()
    defer rshinMemo.Close()

    rshinMemo.Run()
}
