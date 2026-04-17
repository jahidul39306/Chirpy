package main

import (
	"net/http"
	"slices"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleGetAllChirps(w http.ResponseWriter, r *http.Request) {
	authorIDString := r.URL.Query().Get("author_id")
	sortDirection := r.URL.Query().Get("sort")
	if authorIDString != "" {
		authorID, err := uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		chirps, err := cfg.dbQueries.GetChirpByAuthorID(r.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if sortDirection == "desc" {
			slices.Reverse(chirps)
			respondWithJSON(w, http.StatusOK, chirps)
			return
		}
		respondWithJSON(w, http.StatusOK, chirps)
		return
	}

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
	if sortDirection == "desc" {
		slices.Reverse(allChirps)
		respondWithJSON(w, http.StatusOK, allChirps)
		return
	}
	respondWithJSON(w, http.StatusOK, allChirps)
}
