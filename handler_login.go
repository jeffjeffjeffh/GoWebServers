package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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


	userResp := userResponse{
		Email: user.Email,
		ID: user.ID,
	}

	data, err := json.Marshal(userResp)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
	}

	writeJSON(w, data, http.StatusOK)
}