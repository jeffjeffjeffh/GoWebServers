package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerChirpsList(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.ListChirps()
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
	}

	data, err := json.Marshal(chirps)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
	}

	writeJSON(w, data, http.StatusOK)
}