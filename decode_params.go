package main

import (
	"encoding/json"
	"errors"
	"log"
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
	ChirpyRed *bool `json:"chirpy_red"`
}

func decodeUserParams(r *http.Request) (userParams, error) {
	params := userParams{}
	decoder := json.NewDecoder(r.Body)
	
	err := decoder.Decode(&params)
	if err != nil {
		log.Println(err)
		return userParams{}, err
	}
	if params.Email == nil {
		err := errors.New("invalid POST request; no user email found")
		log.Println(err)
		return userParams{}, err
	}
	if params.Password == nil {
		err := errors.New("invalid POST request; no user password found")
		log.Println(err)
		return userParams{}, err
	}

	return params, nil
}

type webhookParams struct{
	Event *string `json:"event"`
	Data *struct{
		UserID *int `json:"user_id"`
	} `json:"data"`
}

func decodeWebhooksParams(r *http.Request) (webhookParams, error) {
	log.Println("decoding webhook params...")

	hook := webhookParams{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&hook)
	if err != nil {
		log.Println(err)
		return webhookParams{}, err
	}
	if hook.Event == nil || hook.Data == nil || hook.Data.UserID == nil {
		err := errors.New("malformed hook")
		log.Println(err)
		return webhookParams{}, err
	}

	return hook, nil
}