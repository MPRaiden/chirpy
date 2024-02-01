package main

import (
	"net/http"
	"strings"
	"encoding/json"
	"log"
)

func cleanChirp(chirp string) string {
	badWords := []string{"kerfuffle","sharbert","fornax"}

	for _, bw := range badWords {
		if strings.Contains(strings.ToLower(chirp), bw) {
			chirp = strings.Replace(chirp, bw, "****", -1)
			chirp = strings.Replace(chirp, strings.Title(bw), "****", -1)
		}
	}

	return chirp
}


func respondWithError(w http.ResponseWriter, code int, msg string) {
    response := map[string]string{"error": msg}
    jsonResp, err := json.Marshal(response)
    if err != nil {
        log.Printf("Error marshalling JSON: %s", err)
        w.WriteHeader(500)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(jsonResp)
}


func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    jsonResp, err := json.Marshal(payload)
    if err != nil {
        log.Printf("Error marshalling JSON: %s", err)
        w.WriteHeader(500)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(jsonResp)
}


func validateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	
	cleanedChirp := cleanChirp(params.Body)
	response := map[string]string{"cleaned_body": cleanedChirp}
	respondWithJSON(w, 200, response)
}

