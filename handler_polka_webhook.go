package main

import (
	"net/http"
	"encoding/json"
	"github.com/MPRaiden/chirpy/internal/database"
	"errors"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request){
	type parameters struct {
		Event string `json:"event"`
		Data map[string]int `json:"data"`
	}
	// Check for API key
	apiKey := r.Header.Get("Authorization")
	if apiKey != cfg.polkaAPIKey {
		respondWithError(w, http.StatusUnauthorized, "You are not authorized to make this request")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Event == "user.upgraded" {
		_, err := cfg.DB.UpgradeUser(params.Data["user_id"])
		if err != nil {
			if errors.Is(err, database.ErrNotExist) {
				respondWithError(w, http.StatusNotFound, "User does not exist")
			} else {
				respondWithError(w, http.StatusInternalServerError, "Could not upgrade user to Chirpy Red status")
			}
			return
		}
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status":"success"})
}
