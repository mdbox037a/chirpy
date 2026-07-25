package main

import (
	"log"
	"net/http"
	"time"

	"github.com/mdbox037a/chirpy/internal/auth"
)

func handlerRefreshToken(cfg *apiConfig) (wr http.ResponseWriter, req http.Request) {
	reqRefreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(wr, http.StatusBadRequest, "Bad request - failed to get refresh token from request headers")
		return
	}

	tokenInfo, err := cfg.dbQueries.GetUserFromRefreshToken(req.Context(), reqRefreshToken)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(wr, http.StatusUnauthorized, "Unauthorized - token does not exist in db")
		return
	}
	if tokenInfo.ExpiresAt.Before(time.Now()) || tokenInfo.RevokedAt.Valid {
		respondWithError(wr, http.StatusUnauthorized, "Unauthorized - token has been revoked or is expired")
		return
	}
}
