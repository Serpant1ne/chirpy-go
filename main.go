package main

import (
	"log"
	"net/http"
)

const (
	FILEPATHROOT = "."
	PORT         = "8080"
)

func main() {
	handler := http.NewServeMux()
	handler.Handle("/", http.FileServer(http.Dir(FILEPATHROOT)))
	server := &http.Server{
		Handler: handler,
		Addr:    ":" + PORT,
	}

	log.Printf("Serving files from %s on port: %s\n", FILEPATHROOT, PORT)
	log.Fatal(server.ListenAndServe())
}
