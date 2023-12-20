package main

import (
	"flag"
	"internal/database"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	hits int
	db *database.DB
	jwtSecret string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	PORT := "8080"
	DB_FILE := "database.json"
	apiCfg := apiConfig{
		hits: 0,
		jwtSecret: os.Getenv("JWT_SECRET"),
	}

	debugDB := flag.Bool("debug", false, "when set to true, will create a new database.json on every restart")
	flag.Parse()

	if *debugDB {
		db, err := database.CreateDB(DB_FILE)
		if err != nil {
			log.Fatal(err)
		}
		apiCfg.db = db
	} else {
		db, err := database.LoadDB(DB_FILE)
		if err != nil {
			log.Fatal(err)
		}
		apiCfg.db = db
	}

	router := chi.NewRouter()
	fsHandler := apiCfg.handlerIncrementMetrics(http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	router.Handle("/app", fsHandler)
	router.Handle("/app/*", fsHandler)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.handlerAdminGetMetrics)
	router.Mount("/admin", adminRouter)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/chirps", apiCfg.handlerChirpsList)
	apiRouter.Get("/chirps/{id}", apiCfg.handlerChirpsGet)
	apiRouter.Get("/healthz", handlerHealthz)
	apiRouter.Get("/metrics", apiCfg.handlerGetMetrics)
	apiRouter.Get("/reset", apiCfg.handlerResetMetrics)
	apiRouter.Post("/chirps", apiCfg.handlerChirpsCreate)
	apiRouter.Post("/login", apiCfg.handlerLogin)
	apiRouter.Post("/users", apiCfg.handlerUsersCreate)
	apiRouter.Post("/refresh", apiCfg.handlerTokenRefresh)
	apiRouter.Post("/revoke", apiCfg.handlerTokenRevoke)
	apiRouter.Put("/users", apiCfg.handlerUsersUpdate)
	apiRouter.Delete("/chirps/{id}", apiCfg.handlerChirpsDelete)
	router.Mount("/api", apiRouter)


	corsRouter := corsWrapper(router)
	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: corsRouter,
	}

	log.Printf("Server listening on port: " + PORT)
	log.Fatal(server.ListenAndServe())
}