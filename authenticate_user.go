package main

import (
	"auth"
	"errors"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func authenticateUser(r *http.Request, expectedTokenType, secret string) (*jwt.Token, *string, error) {
	authString, err := auth.GetAuthString(r)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	parsedToken, err := auth.ParseToken(authString, secret)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	tokenType, err := parsedToken.Claims.GetIssuer()
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	if tokenType != expectedTokenType {
		err := errors.New("wrong token type")
		return nil, nil, err
	}

	tokenIsValid := parsedToken.Valid
	if !tokenIsValid {
		err := errors.New("invalid auth token")
		log.Println(err)
		return nil, nil, err
	}

	log.Println("authentication successful")
	return parsedToken, &authString, nil
}