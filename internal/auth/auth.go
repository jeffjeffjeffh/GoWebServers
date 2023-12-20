package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwt(id int, reqExpiration *int, tokenType string) (*jwt.Token, error) {
	REFRESH_TOKEN_EXPIRATION := time.Hour * time.Duration(24 * 60)
	var expiration time.Duration

	if tokenType == "chirpy-refresh" {
		expiration = REFRESH_TOKEN_EXPIRATION
	} else if tokenType == "chirpy-access" {
		expiration = validateExpiration(reqExpiration)
	} else {
		err := errors.New("invalid token type")
		log.Println(err)
		return nil, err
	}

	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Issuer:    tokenType,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
		Subject:   fmt.Sprint(id),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token, nil
}

func validateExpiration(expInSeconds *int) time.Duration {
	var expiration time.Duration
	if expInSeconds == nil || time.Duration(*expInSeconds) > time.Duration(time.Hour) {
		expiration = time.Hour
	} else {
		expiration = time.Duration(*expInSeconds) * time.Second
	}

	return expiration
}

func GetTokenFromRequest(r *http.Request, secret string) (*jwt.Token, error) {
	authString, err := GetAuthString(r)
	if err != nil {
		return nil, err
	}

	token, err := ParseToken(authString, secret)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func GetAuthString(r *http.Request) (string, error) {
	authHeader := strings.Split(r.Header.Get("Authorization"), " ")
	if len(authHeader) != 2 {
		err := errors.New("malformed auth header")
		log.Println(err)
		return "", err
	}

	return authHeader[1], nil
}

func ParseToken(tokenString, secret string) (*jwt.Token, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(token *jwt.Token) (interface{}, error) { return []byte(secret), nil},
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return token, nil
}