package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type chirpParams struct{
	Body *string `json:"body"`
}

func decodeChirpParams(r *http.Request) (chirpParams, error) {
	params := chirpParams{}
	decoder := json.NewDecoder(r.Body)
	
	err := decoder.Decode(&params)
	if params.Body == nil {
		return chirpParams{}, errors.New("invalid POST request; no body found")
	}
	if err != nil {
		return chirpParams{}, err
	}

	return params, err
}

type userParams struct{
	Email *string `json:"email"`
	Password *string `json:"password"`
	Expiration *int `json:"expires_in_seconds"`
}

func decodeUserParams(r *http.Request) (userParams, error) {
	params := userParams{}
	decoder := json.NewDecoder(r.Body)
	
	err := decoder.Decode(&params)
	if params.Email == nil {
		return userParams{}, errors.New("invalid POST request; no user email found")
	}
	if params.Password == nil {
		return userParams{}, errors.New("invalid POST request; no user password found")
	}
	if err != nil {
		return userParams{}, err
	}

	return params, err
}