package main

import (
	"net/http"
	"encoding/json"
)

type Chirp struct {
	ID int `json:"id"`
	Body string `json:"body"`
}

func decodeParams(r *http.Request) (parameters, error) {
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	return params, err
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	params, err := decodeParams(r)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	validatedChirp, err := validateChirpLength(params.Chirp)
	if err != nil {
		type returnVals struct {
				Error string `json:"error"`
			}
			respBody := returnVals{
				Error: "Chirp is too long",
			}
			
		writeJSON(w, respBody, http.StatusBadRequest)
	}
}

func validateChirpLength(text string) (Chirp, error) {
		maxChirpLength := 140
		if len(text) > maxChirpLength {
			
	
			data, err := json.Marshal(respBody)
			if err != nil {
				writeError(w, err, http.StatusInternalServerError)
				return
			}
	
			writeJSON(w, data, http.StatusBadRequest)
			return
		}
}