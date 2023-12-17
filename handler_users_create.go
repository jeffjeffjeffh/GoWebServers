package main

import (
	"encoding/json"
	"net/http"
)

type userResponse struct{
	Email string `json:"email"`
	ID int `json:"id"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	params, err := decodeUserParams(r)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
	}

	createdUser, err := cfg.db.CreateUser(*params.Email, *params.Password)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
	}

	userResp := userResponse{
		createdUser.Email,
		createdUser.ID,
	}

	data, err := json.Marshal(userResp)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
	}

	writeJSON(w, data, http.StatusCreated)
}