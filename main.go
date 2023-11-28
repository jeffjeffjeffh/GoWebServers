package main

import (
	"log"
	"net/http"
)

func main() {
	PORT := "8080"

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(".")))

	corsMux := corsWrapper(mux)

	server := &http.Server{
		Addr: ":" + PORT,
		Handler: corsMux,
	}

	log.Printf("Server listening on port: " + PORT)
	log.Fatal(server.ListenAndServe())
}

func corsWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}