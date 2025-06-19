package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {

	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type userResponse struct {
		ID          string `json:"id"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
		Email       string `json:"email"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}

	decoder := json.NewDecoder(r.Body)
	var req requestBody
	err := decoder.Decode(&req)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Email is required")
		return
	}
	if req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Password is required")
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

	// Get user from DB
	user, err := cfg.db.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User not found")
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}
	updatedUser, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:           user.ID,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			respondWithError(w, http.StatusBadRequest, "Email already exists")
			return
		}
		log.Printf("Error updating user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	response := userResponse{
		ID:          updatedUser.ID.String(),
		CreatedAt:   updatedUser.CreatedAt.String(),
		UpdatedAt:   updatedUser.UpdatedAt.String(),
		Email:       updatedUser.Email,
		IsChirpyRed: updatedUser.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusOK, response)
	return
}
