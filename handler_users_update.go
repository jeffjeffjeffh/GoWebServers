package main

import (
	"encoding/json"
	"errors"
	"internal/auth"
	"log"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

type userUpdateResponse struct{
	Email string `json:"email"`
	ID int `json:"id"`
}

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	token, err := cfg.authenticateUser(w, r)
	if err != nil {
		// error is already logged and written to the response
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

func (cfg *apiConfig) authenticateUser(w http.ResponseWriter, r *http.Request) (*jwt.Token, error) {
	authString, err := auth.GetAuthString(r)
	if err != nil {
		log.Println(err)
		writeError(w, err, http.StatusUnauthorized)
		return nil, err
	}

	parsedToken, err := auth.ParseToken(authString, cfg.jwtSecret)
	if err != nil {
		log.Println(err)
		writeError(w, err, http.StatusUnauthorized)
		return nil, err
	}

	tokenType, err := parsedToken.Claims.GetIssuer()
	if err != nil {
		log.Println(err)
		writeError(w, err, http.StatusInternalServerError)
		return nil, err
	}
	
	if tokenType == "chirpy-refresh" {
		err := errors.New("received refresh token, expected other type")
		writeError(w, err, http.StatusUnauthorized)
		return nil, err
	}

	tokenIsValid := parsedToken.Valid
	if !tokenIsValid {
		err := errors.New("invalid authorization token")
		log.Println("invalid authorization token")
		writeError(w, err, http.StatusUnauthorized)
		return nil, err
	}
	
	log.Println("user authenticated")
	return parsedToken, nil
}