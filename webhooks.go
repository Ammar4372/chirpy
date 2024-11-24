package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Ammar4372/chirpy/internal/auth"
	"github.com/google/uuid"
)

type Req struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	key, _ := auth.GetApiKey(r.Header)
	if key != cfg.polkaApiKey {
		responseWithJson(w, http.StatusUnauthorized, "")
		return
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	body := Req{}
	err := decoder.Decode(&body)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid json", err)
		return
	}
	if body.Event != "user.upgraded" {
		responseWithJson(w, http.StatusNoContent, "")
		return
	}
	err = cfg.db.UpgradeUser(context.Background(), body.Data.UserID)
	if err != nil {
		responseWithError(w, http.StatusNotFound, "User not found", nil)
	}
	responseWithJson(w, http.StatusNoContent, "")
}
