package main

import (
	"errors"
	"internal/auth"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerTokenRevoke(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetAuthString(r)
	if err != nil {
		writeError(w, err, http.StatusUnauthorized)
		return
	}
	
	token, err := auth.ParseToken(tokenString, cfg.jwtSecret)
	if err != nil {
		writeError(w, err, http.StatusUnauthorized)
		return
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		log.Println(err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if issuer != "chirpy-refresh" {
		err := errors.New("wrong token type for revocation")
		log.Println(err)
		writeError(w, err, http.StatusBadRequest)
		return
	}

	err = cfg.db.RevokeToken(tokenString)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	writeJSON(w, nil, http.StatusOK)
}