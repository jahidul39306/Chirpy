package main

import (
	"encoding/json"
	"net/http"

	"github.com/jahidul39306/Chirpy/internal/auth"
	"github.com/jahidul39306/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode user parameter")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to hash password")
		return
	}

	updatedUser, err := cfg.dbQueries.UpdateUserByID(r.Context(), database.UpdateUserByIDParams{
		ID:             userID,
		HashedPassword: hashedPassword,
		Email:          params.Email,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type UserUpdate struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Email     string `json:"email"`
	}
	respondWithJSON(w, http.StatusOK, UserUpdate{
		ID:        updatedUser.ID.String(),
		Email:     updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt.Time.String(),
		UpdatedAt: updatedUser.UpdatedAt.Time.String(),
	})
}
