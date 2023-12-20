package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type userUpdateResponse struct{
	Email string `json:"email"`
	ID int `json:"id"`
}

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	token, _, err := authenticateUser(r, "chirpy-access", cfg.jwtSecret)
	if err != nil {
		writeError(w, err, http.StatusUnauthorized)
		return
	}

	userId, err := token.Claims.GetSubject()
	if err != nil {
		log.Println("could not get subject from token")
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	newUserInfo, err := decodeUserParams(r)
	if err != nil {
		log.Println("error decoding request body")
		writeError(w, err, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(userId)
	if err != nil {
		log.Println("error converting id in request to int")
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	updatedUser, err := cfg.db.UpdateUser(*newUserInfo.Email, *newUserInfo.Password, id)
	if err != nil {
		if err.Error() == "id not found" {
			writeError(w, err, http.StatusBadRequest)
		} else {
			writeError(w, err, http.StatusInternalServerError)
		}
		return
	}

	userResp := userUpdateResponse{
		Email: updatedUser.Email,
		ID: updatedUser.ID,
	}

	data, err := json.Marshal(userResp)
	if err != nil {
		log.Println("error marshalling user response")
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	log.Println("user data updated")
	writeJSON(w, data, http.StatusOK)
}

