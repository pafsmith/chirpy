package main

import (
	"chirpy/internal/auth"
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	chirp_id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
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

	// Get and validate access token (JWT)
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or malformed access token")
		return
	}
	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid access token")
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Not authorized to delete chirp")
		return
	}
	err = cfg.db.DeleteChirp(r.Context(), chirp_id)
	if err != nil {
		log.Printf("Error deleting chirp: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to delete chirp")
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)

}
