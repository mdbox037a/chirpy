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
		respondWithError(wr, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// validation blocks
	if len(post.Body) > 140 {
		respondWithError(wr, http.StatusBadRequest, "Chirp is too long")
		return
	}
	cleanMsg := replaceProfanity(post.Body)

	type successResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	sr := successResponse{
		CleanedBody: cleanMsg,
	}
	respondWithJSON(wr, http.StatusOK, sr)
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
