package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"time"

	"github.com/google/uuid"
)

func handlerUsersCreate(wr http.ResponseWriter, req *http.Request) {
	type user struct {
		ID        uuid.UUID    `json:"id"`
		CreatedAt time.Time    `json:"created_at"`
		UpdatedAt time.Time    `json:"updated_at"`
		Email     mail.Address `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	post := reqBody{}
	err := decoder.Decode(&post)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(wr, http.StatusInternalServerError, "Something went wrong")
		return
	}
}
