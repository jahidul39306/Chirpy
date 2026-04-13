package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"

	"slices"

	"github.com/jahidul39306/Chirpy/internal/auth"
	"github.com/jahidul39306/Chirpy/internal/database"
)

type Chirp struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Body      string `json:"body"`
	UserID    string `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	UserID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode params")
		return
	}

	new_text, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := cfg.dbQueries.CreateChirps(r.Context(), database.CreateChirpsParams{
		Body:   new_text,
		UserID: UserID,
	})
	if err != nil {
		log.Printf("Error creating chirp: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Unable to create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt.String(),
		UpdatedAt: chirp.UpdatedAt.String(),
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
	})

}

func validateChirp(text string) (string, error) {
	size := utf8.RuneCountInString(text)

	if size > 140 {
		return "", fmt.Errorf("Chirp is too long")
	}

	bad_words := []string{"kerfuffle", "sharbert", "fornax"}
	target_words := []string{}
	words := strings.Split(text, " ")

	for _, w := range words {
		found := slices.Contains(bad_words, strings.ToLower(w))
		if found {
			target_words = append(target_words, "****")
		} else {
			target_words = append(target_words, w)
		}
	}
	new_text := strings.Join(target_words, " ")
	return new_text, nil
}
