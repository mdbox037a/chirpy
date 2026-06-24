package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerValidateChirp(wr http.ResponseWriter, req *http.Request) {
	type reqBody struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	post := reqBody{}
	err := decoder.Decode(&post)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(wr, 500, "Something went wrong")
		return
	}

	if len(post.Body) > 140 {
		respondWithError(wr, 400, "Chirp is too long")
		return
	}

	// TODO: call return json func
}

func respondWithError(wr http.ResponseWriter, code int, msg string) {
	type resError struct {
		Error string `json:"error"`
	}

	resErr := resError{
		Error: msg,
	}
	data, err := json.Marshal(resErr)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		wr.WriteHeader(500)
		return
	}

	log.Printf("%s", msg)
	wr.WriteHeader(code)
	wr.Header().Set("Content-Type", "application/json")
	wr.Write(data)
}
