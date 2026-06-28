package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type reqChirp struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func handlerNewChirp(wr http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	chirp := reqChirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(wr, http.StatusInternalServerError, "Failed to decode client request into chirp struct")
		return
	}

	msg := validateChirp(&chirp)
	if msg != "" {
		respondWithError(wr, http.StatusBadRequest, msg)
	}
	// TODO: bookmark June 28, 2026
	// next call CreateChirp to add to db
}

func validateChirp(chirp *reqChirp) string {
	if len(chirp.Body) > 140 {
		return fmt.Sprint("Chirp body is too long (max 140 characters)")
	}
	cleanMsg := replaceProfanity(chirp.Body)
	chirp.Body = cleanMsg
	return ""
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
