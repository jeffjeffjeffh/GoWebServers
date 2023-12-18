package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type userUpdateResponse struct{
	Email string `json:"email"`
	ID int `json:"id"`
}

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	token, err := cfg.validateUser(w, r)
	if err != nil {
		return
	}
	if token == nil {
		return
	}

	userId, err := token.Claims.GetSubject()
	if err != nil {
		log.Println("could not get Subject from token")
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

	writeJSON(w, data, http.StatusOK)
}

func (cfg *apiConfig) validateUser(w http.ResponseWriter, r *http.Request) (*jwt.Token, error) {
	reqToken := r.Header.Get("Authorization")
	if len(reqToken) == 0 {
		err := errors.New("no auth header included in request")
		writeError(w, err, http.StatusUnauthorized)
		return nil, err
	}

	strippedToken := reqToken[strings.Index(reqToken, " ")+1:]
	if strippedToken == "" {
		err := errors.New("malformed authorization header")
		writeError(w, err, http.StatusUnauthorized)
		return nil, err
	}

	claims := jwt.RegisteredClaims{}
	
	parsedToken, err := jwt.ParseWithClaims(strippedToken, &claims, func(token *jwt.Token) (interface{}, error) {return []byte(cfg.jwtSecret), nil})
	if err != nil {
		log.Printf("%s: %s", err.Error(), strippedToken)
		writeError(w, err, http.StatusUnauthorized)
		return nil, err
	}

	tokenType, err := parsedToken.Claims.GetIssuer()
	if err != nil {
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