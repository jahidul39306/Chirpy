package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jahidul39306/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	const maxExpiry = 3600
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode parameters")
		return
	}

	if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds > maxExpiry {
		params.ExpiresInSeconds = maxExpiry
	}

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Wrong credentials")
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Wrong credentials")
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secretKey, time.Duration(params.ExpiresInSeconds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)

	type userInfo struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Email     string `json:"email"`
		Token     string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, userInfo{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt.Time.String(),
		UpdatedAt: user.UpdatedAt.Time.String(),
		Email:     user.Email,
		Token:     token,
	})
}
