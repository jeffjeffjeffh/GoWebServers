package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	params, err := decodeUserParams(r)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
	}

	createdUser, err := cfg.db.CreateUser(*params.Email)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
	}

	data, err := json.Marshal(createdUser)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
	}

	writeJSON(w, data, http.StatusCreated)
}