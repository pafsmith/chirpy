package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
}

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
		Body string `json:"body"`
	}

	// --- JWT Authentication ---
	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
		return
	}

	userID, err := auth.ValidateJWT(tokenStr, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}
	// --- End JWT Authentication ---

	decoder := json.NewDecoder(r.Body)
	var req requestBody
	err = decoder.Decode(&req)
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
		UserID: userID,
	})

	if err != nil {
		log.Printf("Error creating chirp: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create chirp")
		return
	}
	resp := Chirp{

		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      cleaned,
		UserID:    chirp.UserID,
	}
	respondWithJSON(w, http.StatusCreated, resp)
	return
}
