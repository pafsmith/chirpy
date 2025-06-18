package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	type chirpResponse struct {
		ID        string    `json:"id"`
		CreatedAt string    `json:"created_at"`
		UpdatedAt string    `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		log.Printf("Error fetching chirps: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch chirps")
		return
	}
	var resp []chirpResponse
	for _, chirp := range chirps {
		resp = append(resp, chirpResponse{
			ID:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt.String(),
			UpdatedAt: chirp.UpdatedAt.String(),
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, resp)
	return
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	type chirpResponse struct {
		ID        string    `json:"id"`
		CreatedAt string    `json:"created_at"`
		UpdatedAt string    `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	idStr := r.PathValue("id")
	chirp_id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID format")
		return
	}
	chirp, err := cfg.db.GetChirpByID(r.Context(), chirp_id)

	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Chirp not found")
			return
		}
		log.Printf("Error fetching chirp: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch chirp")
		return
	}
	resp := chirpResponse{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt.String(),
		UpdatedAt: chirp.UpdatedAt.String(),
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	respondWithJSON(w, http.StatusOK, resp)
	return
}
