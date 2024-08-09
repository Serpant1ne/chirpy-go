package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	FILEPATHROOT = "."
	PORT         = "8080"
)

type apiConfig struct {
	fileServerHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	cfg.fileServerHits++
	return next
}

func (cfg *apiConfig) getHits() int {
	return cfg.fileServerHits
}
func (cfg *apiConfig) resetHits() {
	cfg.fileServerHits = 0
}

func main() {
	mux := http.NewServeMux()
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(FILEPATHROOT)))
	apiCfg := apiConfig{
		fileServerHits: 0,
	}
	//main handler
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(handler))

	//server health handler
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	//req number handler
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Hits: %d", apiCfg.getHits())))
	})

	//req number reset handler
	mux.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("charset", "utf-8")
		w.WriteHeader(http.StatusOK)
		apiCfg.resetHits()
		w.Write([]byte("Reset succesful"))
	})

	//server setup
	server := &http.Server{
		Handler: mux,
		Addr:    ":" + PORT,
	}

	//server launch
	log.Printf("Serving files from %s on port: %s\n", FILEPATHROOT, PORT)
	log.Fatal(server.ListenAndServe())
}
