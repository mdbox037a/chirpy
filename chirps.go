package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/mdbox037a/chirpy/internal/database"
)

func (cfg *apiConfig) handlerNewChirp(wr http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	reqChirp := database.CreateChirpParams{}
	err := decoder.Decode(&reqChirp)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(wr, http.StatusInternalServerError, "Failed to decode client request into chirp struct")
		return
	}

	msg := validateChirp(&reqChirp)
	if msg != "" {
		respondWithError(wr, http.StatusBadRequest, msg)
	}

	resChirp, err := cfg.dbQueries.CreateChirp(req.Context(), reqChirp)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(wr, http.StatusInternalServerError, "Failed to add chirp to database")
		return
	}
	respondWithJSON(wr, http.StatusCreated, resChirp)
}

func validateChirp(reqChirp *database.CreateChirpParams) string {
	// maybe a bit clunky now, but leaving room for more validation steps later
	if len(reqChirp.Body) > 140 {
		return fmt.Sprint("Chirp body is too long (max 140 characters)")
	}
	cleanMsg := replaceProfanity(reqChirp.Body)
	reqChirp.Body = cleanMsg
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
