package main

import (
	"net/http"
	"strconv"
	"github.com/MPRaiden/chirpy/internal/auth"
	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerChirpsDelete (w http.ResponseWriter, r *http.Request) {
	chirpIDString := chi.URLParam(r, "chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not find specified chirp.")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to extract bearer token.")
		return
	}

	userIDString, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "You are not authorized to delete this chirp.")
		return
	}
	
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error processinng user ID")
		return
	}
	
	if dbChirp.AuthorID != userID {
		respondWithError(w, http.StatusForbidden, "You can't delete someone else's chirp.")
		return
	}

	err = cfg.DB.DeleteChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting chirp.")
		return
	}

	w.WriteHeader(http.StatusOK)
	// possibly add a success response message or chirpHasBeenDeleted message here


}
