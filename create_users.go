package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handerCreateUser(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type userResponse struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Email     string `json:"email"`
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
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
		// log.Printf("Error hashing password: %v", err)
		// respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		// return
	}
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:        req.Email,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			respondWithError(w, http.StatusBadRequest, "Email already exists")
			return
		}
		log.Printf("Error creating user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	resp := userResponse{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
		Email:     user.Email,
	}
	respondWithJSON(w, http.StatusCreated, resp)
	return

}
