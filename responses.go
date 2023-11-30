package main

import (
	"log"
	"net/http"
)

func writeError(w http.ResponseWriter, e error, code int) {
	w.WriteHeader(code)
	log.Printf("Error marshalling JSON: %s", e)
}

func writeJSON(w http.ResponseWriter, data []byte, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}