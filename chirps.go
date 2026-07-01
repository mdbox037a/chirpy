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

func (cfg *apiConfig) handlerChirpsGet(wr http.ResponseWriter, req *http.Request) {
	dbChirps, err := cfg.dbQueries.GetChirps(req.Context())
	if err != nil {
		log.Printf("Erro: %v", err)
		respondWithError(wr, http.StatusInternalServerError, "Failed to retrieve chirps from database")
		return
	}

	resChirps := make([]resChirp, len(dbChirps))
	for i, dbChirp := range dbChirps {
		resChirp := mapDbChirpToResChirp(dbChirp)
		resChirps[i] = resChirp
	}

	respondWithJSON(wr, http.StatusOK, resChirps)
}

func (cfg *apiConfig) handlerChirpNew(wr http.ResponseWriter, req *http.Request) {
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
	replaceProfanity(reqChirp)
	return ""
}

func replaceProfanity(reqChirp *reqChirp) {
	profanity := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Split(reqChirp.Body, " ")
	for i, word := range words {
		if _, exists := profanity[strings.ToLower(word)]; exists {
			words[i] = "****"
		}
	}

	reqChirp.Body = strings.Join(words, " ")
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
