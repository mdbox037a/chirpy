package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mdbox037a/chirpy/internal/auth"
	"github.com/mdbox037a/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type userReqParams struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds *int   `json:"expires_in_seconds"`
}

func (cfg *apiConfig) handlerUsersCreate(wr http.ResponseWriter, req *http.Request) {

	var params userReqParams
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(wr, http.StatusInternalServerError, "Something went wrong - failed to decode request body")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(wr, http.StatusInternalServerError, "Something went wrong - failed to hash user password")
	}

	crUsPar := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}

	dbUser, err := cfg.dbQueries.CreateUser(req.Context(), crUsPar)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(wr, http.StatusInternalServerError, "Something went wrong - failed to add user to database")
		return
	}

	resUser := mapDBUserToResUser(dbUser)
	respondWithJSON(wr, http.StatusCreated, resUser)
}

func (cfg *apiConfig) handlerUsersLogin(wr http.ResponseWriter, req *http.Request) {
	var params userReqParams
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(wr, http.StatusInternalServerError, "Something went wrong - failed to decode login request body")
		return
	}

	dbUser, err := cfg.dbQueries.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		log.Printf("Error: %v", err)
		// TODO: keeping it simple for now - may want to revisit to handle for server-side issue on the lookup
		respondWithError(wr, http.StatusUnauthorized, "incorrect email or password")
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(wr, http.StatusInternalServerError, "Something went wrong - password hashing comparison failed")
		return
	}

	if !match {
		respondWithError(wr, http.StatusUnauthorized, "incorrect email or password")
		return
	}

	expiry := 3600
	if params.ExpiresInSeconds != nil {
		requestedSeconds := *params.ExpiresInSeconds
		if requestedSeconds < 3600 && requestedSeconds > 0 {
			expiry = requestedSeconds
		}
	}

	resUser := mapDBUserToResUser(dbUser)

	token, err := auth.MakeJWT(resUser.ID, cfg.jwtSecret, time.Duration(expiry)*time.Second)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(wr, http.StatusInternalServerError, "Something went wrong - failed to generate JWT")
		return
	}

	respondWithJSON(wr, http.StatusOK, resUser)
}

func mapDBUserToResUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
}
