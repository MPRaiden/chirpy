package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/MPRaiden/chirpy/database"
	"encoding/json"
	"strconv"
	"os"
	"flag"
)

type apiConfig struct {
	fileserverHits int
}

var db *database.DB

func main() {
	// Deletes database.json on server start if debug mode
	var debug bool
	flag.BoolVar(&debug, "debug", false, "Enable debug mode")
	flag.Parse()
	if debug {
	deleteDatabase()
	}

	var err error
	db, err = database.NewDB("database.json")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %s", err)
	}
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	// Create routers using chi library
	r := chi.NewRouter()
	api := chi.NewRouter()
	admin := chi.NewRouter()

	// Serve files from root project
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	// Handle different endpoints
	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)
	api.Get("/healthz", handlerReadiness)
	api.Post("/chirps", postChirp)
	api.Get("/chirps", getChirps)
	api.Get("/chirps/{chirp_id}", getChirp)
	api.Post("/users", createUser)
	api.HandleFunc("/reset", apiCfg.handlerReset)
	admin.Get("/metrics", apiCfg.handlerMetrics)

	// Mount different routers
	r.Mount("/api", api)
	r.Mount("/admin", admin)

	corsMux := middlewareCors(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func deleteDatabase() {
    err := os.Remove("database.json")
    if err != nil {
        log.Fatalf("Failed to delete the database: %s", err)
    } else {
        log.Println("Successfully deleted the database.")
    }
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "<html>\n<body>\n<h1>Welcome, Chirpy Admin</h1>\n<p>Chirpy has been visited %d times!</p>\n</body>\n</html>", cfg.fileserverHits)
}


func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func validateChirp(body string) (string, error) {
	if len(body) > 140 {
		return "", fmt.Errorf("Chirp is too long")
	}

	cleanedChirp := cleanChirp(body)
	return cleanedChirp, nil
}


func postChirp(w http.ResponseWriter, r *http.Request) {
    var newChirp database.Chirp
    err := json.NewDecoder(r.Body).Decode(&newChirp)
    if err != nil {
        log.Printf("Error decoding parameters: %s", err)
        respondWithError(w, 500, "Unable to parse request body.")
        return
    }

    validatedChirp, err := validateChirp(newChirp.Body)
    if err != nil {
        respondWithError(w, 400, err.Error())
        return
    }

    createdChirp, err := db.CreateChirp(validatedChirp)
    if err != nil {
        log.Printf("Error creating chirp: %s", err)
        respondWithError(w, 500, "Unable to create chirp.")
        return
    }

	respondWithJSON(w, http.StatusCreated, createdChirp)
}

func getChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := db.GetChirps()
	if err != nil {
		log.Printf("Error getting chirps: %s", err)
		respondWithError(w, 500, "Unable to get chirps.")
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func getChirp(w http.ResponseWriter, r *http.Request) {
	// Get chirp id
	chirpID := chi.URLParam(r, "chirp_id")
	id, err := strconv.Atoi(chirpID)
	if err != nil {
		log.Printf("Error getting chirps: %s", err)
		respondWithError(w, 500, "Unable to convert string to int.")
		return
	}

	// Load chirps DB
	chirps, err := db.GetChirps()
	if err != nil {
		log.Printf("Error getting chirps: %s", err)
		respondWithError(w, 500, "Unable to get chirps.")
		return
	}

	// First specific chirp
	for _, chirp := range chirps {
		if chirp.ID == id {
			respondWithJSON(w, http.StatusOK, chirp)
			return
		}
	}
	respondWithError(w, 404, "Not found")
}

func createUser(w http.ResponseWriter, r *http.Request) {
    var newUser database.User
    err := json.NewDecoder(r.Body).Decode(&newUser)
    if err != nil {
        log.Printf("Error decoding parameters: %s", err)
        respondWithError(w, 500, "Unable to parse request body.")
        return
    }

    createdUser, err := db.CreateUser(newUser.Email)
    if err != nil {
        log.Printf("Error creating user: %s", err)
        respondWithError(w, 500, "Unable to create user.")
        return
    }

    respondWithJSON(w, http.StatusCreated, createdUser)
}

