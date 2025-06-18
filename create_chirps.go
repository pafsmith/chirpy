package main

import (
	"chirpy/internal/database"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func cleanProfanity(body string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(body, " ")
	for i, word := range words {
		for _, bad := range badWords {
			if strings.ToLower(word) == bad {
				words[i] = "****"
			}
		}
	}
	return strings.Join(words, " ")
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type responseBody struct {
		ID        string    `json:"id"`
		CreatedAt string    `json:"created_at"`
		UpdatedAt string    `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	var req requestBody
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	if len(req.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleaned := cleanProfanity(req.Body)
	if cleaned == "" {
		respondWithError(w, http.StatusBadRequest, "Chirp cannot be empty")
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: req.UserID,
	})

	if err != nil {
		log.Printf("Error creating chirp: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create chirp")
		return
	}
	resp := responseBody{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt.String(),
		UpdatedAt: chirp.UpdatedAt.String(),
		Body:      cleaned,
		UserID:    chirp.UserID,
	}
	respondWithJSON(w, http.StatusCreated, resp)
	return
}
