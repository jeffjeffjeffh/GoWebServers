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
	RefreshToken string `json:"refresh_token"`
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
	newToken := cfg.generateJwt(user.ID, expiration, "chirpy-access")
	signedToken, err := newToken.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		log.Println(err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	REFRESH_EXPIRATION := time.Hour * time.Duration(24 * 60)
	newRefreshToken := cfg.generateJwt(user.ID, REFRESH_EXPIRATION, "chirpy-refresh")
	signedRefreshToken, err := newRefreshToken.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		log.Println(err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	log.Println("tokens generated")

	userResp := userLoginResponse{
		Email: user.Email,
		ID: user.ID,
		Token: signedToken,
		RefreshToken: signedRefreshToken,
	}

	data, err := json.Marshal(userResp)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
	}

	writeJSON(w, data, http.StatusOK)
}

func validateExpiration(expInSeconds *int) time.Duration {
	var expiration time.Duration
	if expInSeconds == nil || time.Duration(*expInSeconds) > time.Duration(time.Hour) {
		expiration = time.Hour
	} else {
		expiration = time.Duration(*expInSeconds) * time.Second
	}

	return expiration
}

func (cfg *apiConfig) generateJwt(id int, expiration time.Duration, tokenType string) *jwt.Token {
	now := time.Now().UTC()

	claims := jwt.RegisteredClaims{
		Issuer: tokenType,
		IssuedAt: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
		Subject: fmt.Sprint(id),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token
}