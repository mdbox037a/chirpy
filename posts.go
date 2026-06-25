package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

	// validation blocks
	if len(post.Body) > 140 {
		respondWithError(wr, 400, "Chirp is too long")
		return
	}
	cleanMsg := replaceProfanity(post.Body)

	type successResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	sr := successResponse{
		CleanedBody: cleanMsg,
	}
	respondWithJSON(wr, 200, sr)
}

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
		wr.WriteHeader(500)
		return
	}

	wr.Header().Set("Content-Type", "application/json")
	wr.WriteHeader(code)
	wr.Write(data)
}

func replaceProfanity(msg string) string {
	profanity := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Split(msg, " ")
	for i, word := range words {
		if _, exists := profanity[strings.ToLower(word)]; exists {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
