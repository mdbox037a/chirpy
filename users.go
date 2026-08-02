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
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

type userReqParams struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type polkaUpgradeParams struct {
	Event string          `json:"event"`
	Data  polkaDataParams `json:"data"`
}

type polkaDataParams struct {
	UserID uuid.UUID `json:"user_id"`
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

func (cfg *apiConfig) handlerUsersModify(wr http.ResponseWriter, req *http.Request) {
	var reqParams userReqParams
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&reqParams)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(wr, http.StatusInternalServerError, "Something went wrong - failed to decode modify user body")
		return
	}

	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(wr, http.StatusUnauthorized, "Bad request - malformed or missing access token in request header")
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(wr, http.StatusUnauthorized, "Unauthorized - could not validate JWT")
		return
	}

	hashedNewPassword, err := auth.HashPassword(reqParams.Password)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(wr, http.StatusInternalServerError, "Something went wrong - failed to hash new password")
		return
	}

	upUsPar := database.UpdateUserParams{
		ID:             userID,
		Email:          reqParams.Email,
		HashedPassword: hashedNewPassword,
	}

	dbUser, err := cfg.dbQueries.UpdateUser(req.Context(), upUsPar)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(wr, http.StatusInternalServerError, "Something went wrong - failed to update user info in database")
		return
	}

	resUser := mapDBUserToResUser(dbUser)
	respondWithJSON(wr, http.StatusOK, resUser)
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

	resUser := mapDBUserToResUser(dbUser)

	token, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(wr, http.StatusInternalServerError, "Something went wrong - failed to generate JWT")
		return
	}
	resUser.Token = token

	refreshToken := auth.MakeRefreshToken()
	refreshTokenParams := database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: dbUser.ID,
	}
	err = cfg.dbQueries.CreateRefreshToken(req.Context(), refreshTokenParams)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(wr, http.StatusInternalServerError, "Something went wrong - failed to add refresh token to db")
		return
	}
	resUser.RefreshToken = refreshToken

	respondWithJSON(wr, http.StatusOK, resUser)
}

func (cfg *apiConfig) handlerUserUpgrade(wr http.ResponseWriter, req *http.Request) {
	var reqUpgradeParams polkaUpgradeParams
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&reqUpgradeParams)
	if err != nil {
		log.Printf("Error: %v", err)
		respondWithError(wr, http.StatusInternalServerError, "Something went wrong - failed to decode request json data")
		return
	}

	if reqUpgradeParams.Event != "user.upgraded" {
		wr.WriteHeader(http.StatusNoContent)
		return
	}

	dbUserData, err := cfg.dbQueries.UpgradeUser(req.Context(), reqUpgradeParams.Data.UserID)
	if err != nil || !dbUserData.IsChirpyRed {
		log.Printf("Error: %v", err)
		respondWithError(wr, http.StatusNotFound, "Not found - user not upgraded")
		return
	}
	wr.WriteHeader(http.StatusNoContent)
}

func mapDBUserToResUser(dbUser database.User) User {
	return User{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
	}
}
