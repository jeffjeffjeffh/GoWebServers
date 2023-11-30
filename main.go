package main

import (
	// "fmt"
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
	fsHandler := apiCfg.handlerIncrementMetrics(http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	router.Handle("/app", fsHandler)
	router.Handle("/app/*", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/metrics", apiCfg.handlerGetMetrics)
	apiRouter.Get("/reset", apiCfg.handlerResetMetrics)
	apiRouter.Post("/chirps", apiCfg.handlerChirpsCreate)
	apiRouter.Get("/healthz", handlerHealthz)
	router.Mount("/api", apiRouter)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.handlerAdminGetMetrics)
	router.Mount("/admin", adminRouter)

	corsRouter := corsWrapper(router)
	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: corsRouter,
	}

	log.Printf("Server listening on port: " + PORT)
	log.Fatal(server.ListenAndServe())
}