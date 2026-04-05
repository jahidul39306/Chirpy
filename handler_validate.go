package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"unicode/utf8"
)

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnMsg struct {
		Cleaned_Body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	size := utf8.RuneCountInString(params.Body)

	if size > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	bad_words := []string{"kerfuffle", "sharbert", "fornax"}
	target_words := []string{}
	words := strings.Split(params.Body, " ")

	for _, w := range words {
		found := slices.Contains(bad_words, strings.ToLower(w))
		if found {
			target_words = append(target_words, "****")
		} else {
			target_words = append(target_words, w)
		}
	}
	new_text := strings.Join(target_words, " ")

	respondWithJSON(w, http.StatusOK, returnMsg{
		Cleaned_Body: new_text,
	})
}
