package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Reset is only allowed in development mode")
	}
	cfg.fileserverHits.Store(0)
	log.Println("Hits counter reset")
	err := cfg.db.ResetUsers(r.Context())
	if err != nil {
		log.Printf("Error resetting users: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to reset users")
		return
	}
	log.Println("Users table reset")

}
