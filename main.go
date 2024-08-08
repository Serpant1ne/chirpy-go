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
	mux := http.NewServeMux()
	mux.Handle("/app/*", http.StripPrefix("/app", http.FileServer(http.Dir(FILEPATHROOT))))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("charset", "utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ok"))
	})
	server := &http.Server{
		Handler: mux,
		Addr:    ":" + PORT,
	}

	log.Printf("Serving files from %s on port: %s\n", FILEPATHROOT, PORT)
	log.Fatal(server.ListenAndServe())
}
