package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")

	id, err := uuid.Parse(chirpIDString)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse chirp ID")
		return
	}

	chirp, err := cfg.dbQueries.GetChirpById(r.Context(), id)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Unable to find the chirp with this id")
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt.String(),
		UpdatedAt: chirp.UpdatedAt.String(),
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
	})
}
