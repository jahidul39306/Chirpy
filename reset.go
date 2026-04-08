package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}
	// w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte(http.StatusText(http.StatusOK)))
	err := cfg.dbQueries.DeleteUsers(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete users")
		return
	}
	cfg.fileserverHits.Store(0)
	type Response struct {
		Msg string `json:"msg"`
	}
	response := Response{
		Msg: "OK",
	}
	respondWithJSON(w, http.StatusOK, response)
}
