package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	token, _, err := authenticateUser(r, "chirpy-access", cfg.jwtSecret)
	if err != nil {
		writeError(w, err, http.StatusUnauthorized)
		return
	}

	idString, err := token.Claims.GetSubject()
	if err != nil {
		log.Println(err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println(err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	chirpIdString := chi.URLParam(r, "id")
	chirpId, err := strconv.Atoi(chirpIdString)
	if err != nil {
		log.Println(err)
		writeError(w, err, http.StatusBadRequest)
		return
	}

	err = cfg.db.DeleteChirp(id, chirpId)
	if err != nil {
		if err.Error() == "chirp not found" {
			writeError(w, err, http.StatusBadRequest)
		} else {
			writeError(w, err, http.StatusForbidden)
		}
		return
	}
	
	w.WriteHeader(http.StatusOK)
}