package main

import (
	"net/http"
	"sort"
	"strconv"
	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
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

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:   dbChirp.ID,
		Body: dbChirp.Body,
		AuthorID: dbChirp.AuthorID,
	})
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	// Get sorting query param
	sortParam := r.URL.Query().Get("sort")
	
	// If chirp id is provided only return that chirp
	author_id := r.URL.Query().Get("author_id")
	if author_id != "" { 
		author_id_int, err := strconv.Atoi(author_id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Could not convert author_id to integer")
			return
		}

		author_chirps := []Chirp{}
		for _, dbChirp := range dbChirps {
			if dbChirp.AuthorID == author_id_int {
				author_chirps = append(author_chirps, Chirp{
					ID: dbChirp.ID,
					Body: dbChirp.Body,
					AuthorID: dbChirp.AuthorID,
				})
			}
		}

		if sortParam == "desc" {
			sort.Slice(author_chirps, func(i, j int) bool {
					return author_chirps[i].ID > author_chirps[j].ID
				})
			} else {
				sort.Slice(author_chirps, func(i, j int) bool {
					return author_chirps[i].ID < author_chirps[j].ID
				})
		}

		respondWithJSON(w, http.StatusOK, author_chirps)
		return
	}
	
	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
			AuthorID: dbChirp.AuthorID,
		})
	}

	if sortParam == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID > chirps[j].ID
		})
	} else {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID < chirps[j].ID
		})
	}
	
	respondWithJSON(w, http.StatusOK, chirps)
}
