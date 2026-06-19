package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	filedserverHits atomic.Int32
}

//middleware that records hit on requests
//wraps each user request path
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler{

	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		//increments fileserverhits
		cfg.filedserverHits.Add(1)
		next.ServeHTTP(w,r)
	})
}

//handler method of type apiconfig that accesses its hits
func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=itf-8")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.filedserverHits.Load())))
}

//handler to reset apiconfig hits
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request){
	cfg.filedserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("Hits reset to 0"))
}


func main() {
	//create apiconfig instance
	apiCfg := &apiConfig{}

	//create app instance
	mux := http.NewServeMux()

	//strips app path, wraps middleware and handles request
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir("."))),
	))
	mux.HandleFunc("/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("/reset", apiCfg.handlerReset)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()
}