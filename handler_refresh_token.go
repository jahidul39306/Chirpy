package main

import (
	"net/http"
	"time"

	"github.com/jahidul39306/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.secretKey, time.Duration(60)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	type responseBody struct {
		AccessToken string `json:"token"`
	}
	respondWithJSON(w, http.StatusOK, responseBody{
		AccessToken: accessToken,
	})
}
