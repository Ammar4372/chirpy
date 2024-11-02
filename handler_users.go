package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Ammar4372/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	params := database.CreateUserParams{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, fmt.Sprintf("%s\n", err))
		return
	}
	params.CreatedAt = time.Now().UTC()
	params.UpdatedAt = time.Now().UTC()
	user, err := cfg.db.CreateUser(context.Background(), params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, fmt.Sprintf("%s\n", err))
		return
	}
	responseWithJson(w, http.StatusCreated, user)

}
