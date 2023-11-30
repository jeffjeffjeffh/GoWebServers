package main

import (
	"encoding/json"
	// "fmt"
	"net/http"
	"slices"
	"strings"
	// "internal/database"
)

type apiConfig struct {
	hits int
}

type parameters struct {
	Chirp string `json:"body"`
}



func (cfg *apiConfig) cleanBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, err := decodeParams(r)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		
		profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
		chirpWords := strings.Split(params.Chirp, " ")
		cleanedWords := []string{}
		
		for i := range chirpWords {
			if slices.Contains(profaneWords, chirpWords[i]) {
				cleanedWords = append(cleanedWords, "****")
				} else {
				cleanedWords = append(cleanedWords, chirpWords[i])
			}
		}

		type returnVals struct {
				CleanedBody string `json:"cleaned_body"`
		}
		respBody := returnVals{
			CleanedBody: strings.Join(cleanedWords, " "),
		}
	
		data, err := json.Marshal(respBody)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
	

		writeJSON(w, data, http.StatusOK)
		next.ServeHTTP(w, r)
	})
}