package main

import (
	"fmt"
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
	fmt.Fprintf(writer, "hello world")
}
