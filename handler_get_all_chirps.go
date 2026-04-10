package main

import (
	"net/http"
)

func (cfg *apiConfig) handleGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.dbQueries.GetAllChirps(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to load all chirps")
		return
	}

	allChirps := []Chirp{}

	for _, chirp := range chirps {
		allChirps = append(allChirps, Chirp{
			ID:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt.String(),
			UpdatedAt: chirp.UpdatedAt.String(),
			Body:      chirp.Body,
			UserID:    chirp.UserID.String(),
		})
	}

	respondWithJSON(w, http.StatusOK, allChirps)
}
