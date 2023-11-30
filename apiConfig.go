package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
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

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	return
}

func (cfg *apiConfig) validateChirp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Chirp string `json:"body"`
		}
		params := parameters{}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&params)
		if err != nil {
			writeError(w, err, http.StatusBadRequest)
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
				writeError(w, err, http.StatusBadRequest)
				return
			}
	
			writeJSON(w, data, http.StatusBadRequest)
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
			writeError(w, err, http.StatusBadRequest)
			return
		}
	
		writeJSON(w, data, http.StatusOK)
		next.ServeHTTP(w, r)
	})
}

