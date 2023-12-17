package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type userLoginResponse struct{
	Email string `json:"email"`
	ID int `json:"id"`
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	params, err := decodeUserParams(r)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
	}
	log.Println("login params decoded")

	user, err := cfg.db.Login(*params.Email, *params.Password)
	if err != nil {
		if err.Error() == "user not found" {
			writeError(w, err, http.StatusNotFound)
			return
		}
		if err.Error() == "wrong password" {
			writeError(w, err, http.StatusUnauthorized)
			return
		}
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	expiration := validateExpiration(params.Expiration)

	newToken := cfg.generateJwt(user.ID, expiration)
	signedString, err := newToken.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("token generated")

	userResp := userLoginResponse{
		Email: user.Email,
		ID: user.ID,
		Token: signedString,
	}

	data, err := json.Marshal(userResp)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
	}

	writeJSON(w, data, http.StatusOK)
}

func validateExpiration(expInSeconds *int) time.Duration {
	var expiration time.Duration
	if expInSeconds == nil || time.Duration(*expInSeconds) > time.Duration(time.Hour * 24) {
		expiration = time.Hour * time.Duration(24)
	} else {
		expiration = time.Duration(*expInSeconds) * time.Second
	}

	return expiration
}

func (cfg *apiConfig) generateJwt(id int, expiration time.Duration) *jwt.Token {
	now := time.Now().UTC()

	claims := jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
		Subject: fmt.Sprint(id),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token
}