package main

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/Ammar4372/chirpy/internal/auth"
	"github.com/Ammar4372/chirpy/internal/database"
)

type resp struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerTokenRefresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "No Refresh Token", err)
		return
	}
	data, err := cfg.db.GetUserFromRefreshToken(context.Background(), token)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Invalid Token", err)
		return
	}
	if data.RevokedAt.Valid || data.ExpiresAt.Before(time.Now()) {
		responseWithError(w, http.StatusUnauthorized, "Token Expired", err)
		return
	}
	access_token, err := auth.MakeJWT(data.UserID, cfg.secret, time.Hour)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "internal server error", err)
		return
	}
	res := resp{
		Token: access_token,
	}
	responseWithJson(w, http.StatusOK, res)
}

func (cfg *apiConfig) handlerTokenRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "No Refresh Token", err)
		return
	}
	params := database.RevokeTokenParams{
		Token: token,
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: time.Now(),
	}
	err = cfg.db.RevokeToken(context.Background(), params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Internal Server Error", err)
		return
	}

	responseWithJson(w, http.StatusNoContent, nil)
}
