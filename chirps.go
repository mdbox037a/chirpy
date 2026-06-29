package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mdbox037a/chirpy/internal/database"
)

type reqChirp struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type resChirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerNewChirp(wr http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	reqChirp := reqChirp{}
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

	createChirpParams := database.CreateChirpParams{
		Body:   reqChirp.Body,
		UserID: reqChirp.UserID,
	}
	dbChirp, err := cfg.dbQueries.CreateChirp(req.Context(), createChirpParams)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(wr, http.StatusInternalServerError, "Failed to add chirp to database")
		return
	}
	resChirp := mapDbChirpToResChirp(dbChirp)
	respondWithJSON(wr, http.StatusCreated, resChirp)
}

func validateChirp(reqChirp *reqChirp) string {
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

func mapDbChirpToResChirp(dbChirp database.Chirp) resChirp {
	return resChirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
}
