package main

import (
	"internal/testDatabase"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	hits int
	db *testDatabase.DB
}

func main() {
	PORT := "8080"
	DB_FILE := "database.json"

	db, err := testDatabase.NewDB(DB_FILE)
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		hits: 0,
		db: db,
	}

	router := chi.NewRouter()
	fsHandler := apiCfg.handlerIncrementMetrics(http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	router.Handle("/app", fsHandler)
	router.Handle("/app/*", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/metrics", apiCfg.handlerGetMetrics)
	apiRouter.Get("/reset", apiCfg.handlerResetMetrics)
	apiRouter.Get("/chirps", apiCfg.handlerChirpsList)
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