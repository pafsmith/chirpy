package main

import (
	"database/sql"
	"log"
	"net/http"
	"sort" // Added for sorting

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	authorID := uuid.Nil
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString != "" {
		authorID, err = uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID")
			return
		}
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		if authorID != uuid.Nil && dbChirp.UserID != authorID {
			continue
		}

		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		})
	}

	// Sorting logic
	sortParam := r.URL.Query().Get("sort")
	if sortParam != "desc" {
		// Default to ascending
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})
	} else {
		// Descending
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[j].CreatedAt.Before(chirps[i].CreatedAt)
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
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
