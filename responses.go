package main

import (
	"log"
	"net/http"
)

func writeError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	log.Printf("Error marshalling JSON: %s", err)
}

func writeJSON(w http.ResponseWriter, data []byte, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}