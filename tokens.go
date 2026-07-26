package main

import (
	"log"
	"net/http"
	"time"

	"github.com/mdbox037a/chirpy/internal/auth"
)

type TokenResponse struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerTokenRefresh(wr http.ResponseWriter, req *http.Request) {
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

	newToken, err := auth.MakeJWT(tokenInfo.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(wr, http.StatusInternalServerError, "Something went wrong - failed to generate JWT")
		return
	}
	resToken := TokenResponse{
		Token: newToken,
	}

	respondWithJSON(wr, http.StatusOK, resToken)
}

func (cfg *apiConfig) handlerTokenRevoke(wr http.ResponseWriter, req *http.Request) {
	reqRefreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(wr, http.StatusBadRequest, "Bad request - failed to get refresh token from request headers")
		return
	}
	err = cfg.dbQueries.RevokeRefreshToken(req.Context(), reqRefreshToken)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(wr, http.StatusInternalServerError, "Something went wrong - failed to revoke refresh token")
	}
	wr.WriteHeader(http.StatusCreated)
}
