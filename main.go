package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	PORT := "8080"
	apiCfg := apiConfig{
		hits: 0,
	}

	router := chi.NewRouter()
	router.Handle("/app/*", apiCfg.incrementCount(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	router.Handle("/app", apiCfg.incrementCount(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	apiRouter := chi.NewRouter()
	apiRouter.Get("/metrics", apiCfg.getCount)
	apiRouter.Get("/reset", apiCfg.resetCount)
	apiRouter.Post("/validate_chirp", apiCfg.validateChirp)
	apiRouter.Get("/healthz", healthzHandler)
	router.Mount("/api", apiRouter)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.adminGetCount)
	router.Mount("/admin", adminRouter)

	corsRouter := corsWrapper(router)
	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: corsRouter,
	}

	log.Printf("Server listening on port: " + PORT)
	log.Fatal(server.ListenAndServe())
}