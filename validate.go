package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidate(w http.ResponseWriter, req *http.Request) {
	type data struct {
		Body string `json:"body"`
	}

	type resp struct {
		Body string `json:"cleaned_body"`
	}
	docoder := json.NewDecoder(req.Body)
	reqBody := data{}
	err := docoder.Decode(&reqBody)
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Something Went Wrong", err)
		return
	}
	if len(reqBody.Body) > 140 {

		responseWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	words := strings.Split(reqBody.Body, " ")
	for i, v := range words {
		lower := strings.ToLower(v)
		if lower == "kerfuffle" || lower == "sharbert" || lower == "fornax" {
			words[i] = "****"
		}
	}
	cleaned_body := strings.Join(words, " ")
	respBody := resp{
		Body: cleaned_body,
	}
	responseWithJson(w, http.StatusOK, respBody)
}
