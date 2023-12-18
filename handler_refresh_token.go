package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type tokenResponse struct{
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	authHeader := strings.Split(r.Header.Get("Authorization"), "")
	if len(authHeader) < 2 {
		err := errors.New("malformed auth header")
		log.Println(err)
		writeError(w, err, http.StatusUnauthorized)
		return
	}

	tokenStr := authHeader[1]
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) { return []byte(cfg.jwtSecret), nil})
	if err != nil {
		log.Println(err)
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

	tokenIsRevoked, err := cfg.db.CheckTokenStatus(tokenStr)
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

	expiration := time.Hour * time.Duration(24 * 60)
	newToken := cfg.generateJwt(id, expiration, "chirpy-refresh")

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