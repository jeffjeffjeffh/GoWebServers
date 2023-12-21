package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
)

func (cfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	authStrings := strings.Split(authHeader, " ")

	if authHeader == "" || len(authStrings) != 2 {
		err := errors.New("missing or malformed authorization")
		log.Println(err)
		writeError(w, err, http.StatusUnauthorized)
		return
	}

	authString := authStrings[1]

	if authString != os.Getenv("POLKA_API_KEY") {
		err := errors.New("polka api key mismatch")
		log.Println(err)
		writeError(w, err, http.StatusUnauthorized)
		return
	}

	hook, err := decodeWebhooksParams(r)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	if *hook.Event != "user.upgraded" {
		log.Println("ignoring webhook")
		w.WriteHeader(http.StatusOK)
		return
	}

	err = cfg.db.UpdateUserMembership(*hook.Data.UserID)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	log.Println("user's chirpy red status updated")
	w.WriteHeader(http.StatusOK)
}