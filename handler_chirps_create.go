package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"github.com/MPRaiden/chirpy/internal/auth"
	"strconv"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
	AuthorID int `json:"author_id"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	var params parameters

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	authorID, err := cfg.GetTokenUserID(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not logged in")
		return
	}

	chirp, err := cfg.DB.CreateChirp(cleaned, authorID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:   chirp.ID,
		Body: chirp.Body,
		AuthorID: chirp.AuthorID,
	})
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(body, badWords)
	return cleaned, nil
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}

func (cfg *apiConfig) GetTokenUserID (r *http.Request) (int, error) {
	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		return 0, err
	}

	userIDStr, err := auth.ValidateJWT(tokenStr, cfg.jwtSecret)
	if err != nil {
		return 0, err
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
