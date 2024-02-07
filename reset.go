package main

import (
	"net/http"
	"github.com/MPRaiden/chirpy/internal/auth"
	"time"
)
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func (cfg *apiConfig) handlerPostRefresh(w http.ResponseWriter, r *http.Request) {
	// Get token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	// Checks if token is of refresh type
	isRefresh, err := auth.IsRefreshToken(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not check token type")
		return
	}
	if !isRefresh {
		respondWithError(w, http.StatusUnauthorized, "Provided token is not a refresh token")
		return
	}

	// Check if the token is revoked
	isRevoked, err := cfg.DB.IsTokenRevoked(token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not check if token is revoked")
		return
	}
	if isRevoked {
		respondWithError(w, http.StatusUnauthorized, "Token has been revoked")
		return
	}

	// Extract the userID from the token's claims
	claims, err := auth.ExtractClaims(token, []byte(cfg.jwtSecret))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not extract claims from token")
		return
	}

	userID, ok := claims["userid"].(float64) // assuming "userid" claim is used to store userID
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "User ID not found in token claims")
		return
	}

	// Generate new JWT access token using the same method you use for login
	newToken, err := auth.MakeJWT(int(userID), cfg.jwtSecret, 1*time.Hour, "chirpy-access")	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate new access token.")
		return
	}

	// Respond with new token
	response := map[string]string{
		"token": newToken,
	}
	respondWithJSON(w, http.StatusOK, response)
	
}
