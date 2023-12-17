package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
	}

	chirp, err := cfg.db.GetChirp(id)
	if err != nil {
		writeError(w, err, http.StatusNotFound)
	}

	resBody, err := json.Marshal(chirp)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
	}

	writeJSON(w, resBody, http.StatusOK)
}