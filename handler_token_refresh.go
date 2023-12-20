package main

import (
	"encoding/json"
	"errors"
	"internal/auth"
	"log"
	"net/http"
	"strconv"
)

type tokenResponse struct{
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerTokenRefresh(w http.ResponseWriter, r *http.Request) {
	authString, err := auth.GetAuthString(r) 
	if err != nil {		
		log.Println(err)
		writeError(w, err, http.StatusUnauthorized)
		return
	}

	token, err := auth.ParseToken(authString, cfg.jwtSecret)
	if err != nil {
		log.Println(err)
		writeError(w, err, http.StatusUnauthorized)
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		log.Println(err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if issuer != "chirpy-refresh" {
		err := errors.New("expected refresh token, received other type")
		log.Println(err)
		writeError(w, err, http.StatusUnauthorized)
		return
	}

	if !token.Valid {
		err := errors.New("invalid token")
		log.Println(err)
		writeError(w, err, http.StatusUnauthorized)
		return
	}

	tokenIsRevoked, err := cfg.db.CheckTokenStatus(authString)
	if err != nil {
		log.Println(err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if tokenIsRevoked {
		err := errors.New("token is revoked")
		log.Println(err)
		writeError(w, err, http.StatusUnauthorized)
		return
	}

	idStr, err := token.Claims.GetSubject()
	if err != nil {
		log.Println(err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	newToken, err := auth.GenerateJwt(id, nil, "chirpy-access")
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	newTokenStr, err := newToken.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		log.Println(err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	tokenRes := tokenResponse{
		Token: newTokenStr,
	}

	data, err := json.Marshal(tokenRes)
	if err != nil {
		log.Println(err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, data, http.StatusOK)
}