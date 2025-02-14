package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/Ammar4372/chirpy/internal/auth"
	"github.com/Ammar4372/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	UserID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	params := database.CreateChripParams{}
	defer req.Body.Close()
	docoder := json.NewDecoder(req.Body)
	err = docoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, fmt.Sprintf("%s\n", err), err)
		return
	}
	if len(params.Body) > 140 {

		responseWithError(w, http.StatusBadRequest, "Chirp is too long", err)
		return
	}
	words := strings.Split(params.Body, " ")
	for i, v := range words {
		lower := strings.ToLower(v)
		if lower == "kerfuffle" || lower == "sharbert" || lower == "fornax" {
			words[i] = "****"
		}
	}
	params.Body = strings.Join(words, " ")
	params.UserID = UserID
	params.CreatedAt = time.Now().UTC()
	params.UpdatedAt = time.Now().UTC()
	chirp, err := cfg.db.CreateChrip(context.Background(), params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, fmt.Sprintf("%s", err), err)
		return
	}
	responseWithJson(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	authorId := req.URL.Query().Get("author_id")

	if authorId != "" {
		authorID, err := uuid.Parse(authorId)
		if err != nil {
			responseWithError(w, http.StatusBadRequest, "Invalid author id", err)
			return
		}
		chirps, err := cfg.db.GetChirpByAuthorId(context.Background(), authorID)
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, fmt.Sprintf("%s", err), err)
			return
		}
		responseWithJson(w, http.StatusOK, chirps)
		return
	}
	chirps, err := cfg.db.GetAllChirps(context.Background())
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, fmt.Sprintf("%s", err), err)
		return
	}
	sortby := req.URL.Query().Get("sort")
	switch sortby {
	case "desc":
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
		fmt.Print(sortby)
		break
	default:
		break
	}
	responseWithJson(w, http.StatusOK, chirps)
}
func (cfg *apiConfig) handlerChirpById(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	Id, err := uuid.Parse(id)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid Id", err)
		return
	}
	chirp, err := cfg.db.GetChirpById(context.Background(), Id)
	if err.Error() == "sql: no rows in result set" {
		responseWithError(w, http.StatusNotFound, "Resource Not Found", err)
		return
	}
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, fmt.Sprintf("%s\n", err), err)
		return
	}
	responseWithJson(w, http.StatusOK, chirp)
}
