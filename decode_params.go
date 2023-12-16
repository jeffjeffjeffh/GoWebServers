package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type parameters struct{
	Chirp *string `json:"body"`
}

func decodeParams(r *http.Request) (parameters, error) {
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	
	err := decoder.Decode(&params)
	if params.Chirp == nil {
		return parameters{}, errors.New("invalid POST request; no body found")
	}
	if err != nil {
		return parameters{}, err
	}

	return params, err
}