package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jahidul39306/Chirpy/internal/auth"
	"github.com/jahidul39306/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode parameters")
		return
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

	token, err := auth.MakeJWT(user.ID, cfg.secretKey, time.Duration(60)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	refresh_token := auth.MakeRefreshToken()
	_, err = cfg.dbQueries.CreateRefreshTokens(r.Context(), database.CreateRefreshTokensParams{
		Token:  refresh_token,
		UserID: user.ID,
	})

	type userInfo struct {
		ID           string `json:"id"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		IsChirpyRed  bool   `json:"is_chirpy_red"`
	}

	respondWithJSON(w, http.StatusOK, userInfo{
		ID:           user.ID.String(),
		CreatedAt:    user.CreatedAt.Time.String(),
		UpdatedAt:    user.UpdatedAt.Time.String(),
		Email:        user.Email,
		Token:        token,
		RefreshToken: refresh_token,
		IsChirpyRed:  user.IsChirpyRed,
	})
}
