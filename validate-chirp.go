package main

import (
	"net/http"
	"encoding/json"
	"log"
)

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
		response := map[string]string{"error":"Chirp is too long"}
		dat, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			w.Write(dat)
		}
		return
	}
	response := map[string]bool{"valid":true}
	dat, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(dat)
	}
}

