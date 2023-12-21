package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsList(w http.ResponseWriter, r *http.Request) {
	var authorId *int
	authorIdStr := r.URL.Query().Get("author_id")
	
	if authorIdStr != "" {
		id, err := strconv.Atoi(authorIdStr)
		if err != nil {
			log.Println(err)
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		authorId = &id
	} else {
		authorId = nil
	}

	chirps, err := cfg.db.ListChirps(authorId)
	if err != nil {
		log.Println(err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(chirps)
	if err != nil {
		log.Println(err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, data, http.StatusOK)
}