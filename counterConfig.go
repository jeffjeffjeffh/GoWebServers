package main

import (
	"fmt"
	"net/http"
)

type apiConfig struct {
	hits int
}

func (config *apiConfig) incrementCount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config.hits++
		next.ServeHTTP(w, r)
	})
}

func (config *apiConfig) getCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %v", config.hits)))
}

func (config *apiConfig) resetCount(w http.ResponseWriter, r *http.Request)  {
	config.hits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits reset to 0: %v", config.hits)))
}