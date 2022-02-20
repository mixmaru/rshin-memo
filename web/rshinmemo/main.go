package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	port := "8080"
	http.HandleFunc("/", hello)
	log.Printf("Server listening on port %s", port)
	log.Print(http.ListenAndServe(":"+port, nil))
}

func hello(writer http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("template/index.html")
	if err != nil {
		log.Fatalf("template error: %v", err)
	}
	if err := t.Execute(writer, struct {
		Title   string
		Content string
	}{
		Title:   "メモ一覧",
		Content: "内容",
	}); err != nil {
		log.Printf("failed to execute template: %v", err)
	}
}
