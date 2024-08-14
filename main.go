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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) getHits() int {
	return cfg.fileServerHits
}
func (cfg *apiConfig) resetHits() {
	cfg.fileServerHits = 0
}

func blockHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func main() {
	mux := http.NewServeMux()
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(FILEPATHROOT)))
	apiCfg := apiConfig{
		fileServerHits: 0,
	}
	//[APP] main handler
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(handler))

	//[API] server health handler
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("POST /api/healthz", blockHandler)
	mux.HandleFunc("DELETE /api/healthz", blockHandler)

	//[API] req number handler
	mux.HandleFunc("GET /api/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Hits: %d", apiCfg.getHits())))
	})

	mux.HandleFunc("POST /api/metrics", blockHandler)
	mux.HandleFunc("DELETE /api/metrics", blockHandler)

	//[API] req number reset handler
	mux.HandleFunc("GET /api/reset", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("charset", "utf-8")
		w.WriteHeader(http.StatusOK)
		apiCfg.resetHits()
		w.Write([]byte("Reset succesful"))
	})

	mux.HandleFunc("POST /api/reset", blockHandler)
	mux.HandleFunc("DELETE /api/reset", blockHandler)

	//server setup
	server := &http.Server{
		Handler: mux,
		Addr:    ":" + PORT,
	}

	//server launch
	log.Printf("Serving files from %s on port: %s\n", FILEPATHROOT, PORT)
	log.Fatal(server.ListenAndServe())

}
