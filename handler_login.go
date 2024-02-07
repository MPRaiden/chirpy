package main

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/MPRaiden/chirpy/internal/auth"

)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
		Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	// Parse the request body
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	// Get the user by their email
	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	// Check their password
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	// Create JWT token
	jwtToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour) 
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create a JWT token")
		return
	}

	// Create refresh token
	refreshToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour*24*60, "chirpy-refresh")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create a refresh token")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
		Token: jwtToken,
		RefreshToken: refreshToken,
	})
}
