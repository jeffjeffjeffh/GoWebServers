package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type apiConfig struct {
	hits int
}

func (cfg *apiConfig) incrementCount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.hits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) getCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %v", cfg.hits)))
}

func (cfg *apiConfig) adminGetCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>`, cfg.hits)))
}

func (cfg *apiConfig) resetCount(w http.ResponseWriter, r *http.Request) {
	cfg.hits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits reset to 0: %v", cfg.hits)))
}

func (cfg *apiConfig) validateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Chirp string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, err, http.StatusBadRequest)
		return
	}

	maxChirpLength := 140
	if len(params.Chirp) > maxChirpLength {
		type returnVals struct {
			Error string `json:"error"`
		}
		respBody := returnVals{
			Error: "Chirp is too long",
		}

		data, err := json.Marshal(respBody)
		if err != nil {
			respondWithError(w, err, http.StatusBadRequest)
			return
		}

		respondWithJSON(w, data, http.StatusBadRequest)
		return
	}

	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	chirpWords := strings.Split(params.Chirp, " ")
	for i := range chirpWords {
		for j := range profaneWords {
			if profaneWords[j] == strings.ToLower(chirpWords[i]) {
				chirpWords[i] = "****"
			}
		}
	}
	cleanedWords := strings.Join(chirpWords, " ")

	type returnVals struct {
			CleanedBody string `json:"cleaned_body"`
	}
	respBody := returnVals{
		CleanedBody: cleanedWords,
	}

	data, err := json.Marshal(respBody)
	if err != nil {
		respondWithError(w, err, http.StatusBadRequest)
		return
	}

	respondWithJSON(w, data, http.StatusOK)
	return
}

func respondWithError(w http.ResponseWriter, e error, code int) {
	w.WriteHeader(code)
	log.Printf("Error marshalling JSON: %s", e)
}

func respondWithJSON(w http.ResponseWriter, data []byte, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}