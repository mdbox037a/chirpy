package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(wr http.ResponseWriter, code int, msg string) {
	type resError struct {
		Error string `json:"error"`
	}

	resErr := resError{
		Error: msg,
	}
	respondWithJSON(wr, code, resErr)
}

func respondWithJSON(wr http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		wr.WriteHeader(http.StatusInternalServerError)
		return
	}

	wr.Header().Set("Content-Type", "application/json")
	wr.WriteHeader(code)
	wr.Write(data)
}
