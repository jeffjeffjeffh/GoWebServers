package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

type Chirp struct{
	AuthorID int `json:"author_id"`
	Body string `json:"body"`
	ID int `json:"id"`
}

type returnErrorVals struct{
	Error string `json:"error"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	params, err := decodeChirpParams(r)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	err = validateChirpLength(*params.Body)
	if err != nil {
		respBody := returnErrorVals{
			Error: "Chirp is too long",
		}

		data, err := json.Marshal(respBody)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
			
		writeJSON(w, data, http.StatusBadRequest)
		return
	}

	cleanedChirp := cleanChirp(*params.Body)

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

	createdChirp, err := cfg.db.CreateChirp(cleanedChirp, id)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(createdChirp)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
	}

	writeJSON(w, data, http.StatusCreated)
}

func validateChirpLength(text string) error {
		maxChirpLength := 140
		if len(text) > maxChirpLength {
			return errors.New("chirp length exceeds 140 characters")
		}
		return nil
}

func cleanChirp(chirp string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	chirpWords := strings.Split(chirp, " ")
	cleanedWords := []string{}
	
	for i := range chirpWords {
		if slices.Contains(profaneWords, chirpWords[i]) {
			cleanedWords = append(cleanedWords, "****")
			} else {
			cleanedWords = append(cleanedWords, chirpWords[i])
		}
	}

	return strings.Join(cleanedWords, " ")
}