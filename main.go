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
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(FILEPATHROOT)))
	apiCfg := apiConfig{
		fileServerHits: 0,
	}

	//[APP] main handler
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(handler))

	//[API] server health handler
	mux.HandleFunc("GET /api/healthz", apiCfg.healthHandler)

	// [ADMIN] req number handler
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)

	//[API] req number reset handler
	mux.HandleFunc("GET /api/reset", apiCfg.resetHandler)

	//server setup
	server := &http.Server{
		Handler: mux,
		Addr:    ":" + PORT,
	}

	//server launch
	log.Printf("Serving files from %s on port: %s\n", FILEPATHROOT, PORT)
	log.Fatal(server.ListenAndServe())

}
