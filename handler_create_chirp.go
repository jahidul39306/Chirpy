package main

import (
	"encoding/json"
	"net/http"

	"github.com/jahidul39306/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	// type parameters struct {
	// 	Body   string `json:"body"`
	// 	UserID string `json:"user_id"`
	// }

	decoder := json.NewDecoder(r.Body)
	params := database.CreateChirpsParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode params")
		return
	}

	// type chirp struct {
	// 	ID        string `json:"id"`
	// 	CreatedAt string `json:"created_at"`
	// 	UpdatedAt string `json:"updated_at"`
	// 	Body      string `json:"body"`
	// 	UserID    string `json:"user_id"`
	// }

	c, err := cfg.dbQueries.CreateChirps(r.Context(), params)

	respondWithJSON(w, http.StatusCreated, c)

}
