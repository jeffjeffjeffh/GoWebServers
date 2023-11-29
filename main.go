package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	PORT := "8080"
	router := chi.NewRouter()
	apiCfg := apiConfig{
		hits: 0,
	}

	router.Handle("/app/*", apiCfg.incrementCount(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	router.Handle("/app", apiCfg.incrementCount(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	router.Get("/metrics", apiCfg.getCount)
	router.Get("/reset", apiCfg.resetCount)
	router.Get("/healthz", healthzHandler)

	corsRouter := corsWrapper(router)
	server := &http.Server{
		Addr: ":" + PORT,
		Handler: corsRouter,
	}

	log.Printf("Server listening on port: " + PORT)
	log.Fatal(server.ListenAndServe())
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}