package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ammar4372/chirpy/internal/auth"
	"github.com/Ammar4372/chirpy/internal/database"
	"github.com/google/uuid"
)

type reqBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type respBody struct {
	Id           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAtAt  time.Time `json:"updated_at"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	body := reqBody{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&body)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid Json", err)
		return
	}
	hash, err := auth.HashPassword(body.Password)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid Password", err)
		return
	}
	params := database.CreateUserParams{}
	params.Email = body.Email
	params.HashedPassword = hash
	params.CreatedAt = time.Now().UTC()
	params.UpdatedAt = time.Now().UTC()
	user, err := cfg.db.CreateUser(context.Background(), params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "internal server error", err)
		return
	}
	responseWithJson(w, http.StatusCreated, user)
}

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, req *http.Request) {
	body := reqBody{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&body)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid Json", err)
		return
	}
	user, err := cfg.db.GetUserByEmail(context.Background(), body.Email)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "internal server error", err)
		return
	}
	err = auth.CheckPasswordHash(body.Password, user.HashedPassword)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Wrong Password", err)
		return
	}

	expiresIn := time.Duration(time.Hour)
	token, err := auth.MakeJWT(user.ID, cfg.secret, expiresIn)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "internal server error", err)
		return
	}
	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "internal server error", err)
		return
	}
	params := database.CreateTokenParams{
		Token:     refresh_token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * 60 * time.Hour),
		RevokedAt: sql.NullTime{
			Valid: false,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = cfg.db.CreateToken(context.Background(), params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "internal server error", err)
		return
	}
	resp := respBody{
		Id:           user.ID,
		Email:        user.Email,
		CreatedAt:    user.CreatedAt,
		UpdatedAtAt:  user.UpdatedAt,
		Token:        token,
		RefreshToken: refresh_token,
	}
	responseWithJson(w, http.StatusOK, resp)
}
