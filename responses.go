package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type resError struct{
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)

	resBody := resError{
		Error: err.Error(),
	}

	data, err := json.Marshal(resBody)
	if err != nil {
		fmt.Println("errception")
		return
	}

	w.Write(data)
}

func writeJSON(w http.ResponseWriter, data []byte, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}