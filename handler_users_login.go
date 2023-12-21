package main

import (
	"encoding/json"
	"internal/auth"
	"log"
	"net/http"
)

type userLoginResponse struct{
	Email string `json:"email"`
	ID int `json:"id"`
	ChirpyRed bool `json:"is_chirpy_red"`
	Token string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("logging user in...")

	params, err := decodeUserParams(r)
	if err != nil {
		log.Println(err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	user, err := cfg.db.Login(*params.Email, *params.Password)
	if err != nil {
		var code int
		if err.Error() == "user not found" {
			code = http.StatusNotFound
		} else if err.Error() == "wrong password" {
			code = http.StatusUnauthorized
		} else {
			code = http.StatusInternalServerError
		}
		log.Println(err)
		writeError(w, err, code)
		return
	}

	newRefreshToken, err := auth.GenerateJwt(user.ID, nil, "chirpy-refresh")
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	signedRefreshToken, err := newRefreshToken.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		log.Println(err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	newToken, err := auth.GenerateJwt(user.ID, params.Expiration, "chirpy-access")
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	
	signedToken, err := newToken.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		log.Println(err)
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	userResp := userLoginResponse{
		Email: user.Email,
		ID: user.ID,
		ChirpyRed: user.ChirpyRed,
		Token: signedToken,
		RefreshToken: signedRefreshToken,
	}

	data, err := json.Marshal(userResp)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
	}

	log.Println("user logged in")
	writeJSON(w, data, http.StatusOK)
}