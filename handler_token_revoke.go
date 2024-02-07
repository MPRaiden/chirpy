package main

import (
	"net/http"
	"github.com/MPRaiden/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevokeTokenn (w http.ResponseWriter, r *http.Request) {
	// Get token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	// Check if token is refresh
	isRefresh, err := auth.IsRefreshToken(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not check token type")
		return
	}

	if isRefresh {
		// Revoke the refresh token in the database.
		err = cfg.DB.RevokeRefreshToken(token)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Could not revoke token")
			return
		 }
	}

	// Respond with a 200 status code and a success message
	w.WriteHeader(http.StatusOK)
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Successfully revoked token"})
}
