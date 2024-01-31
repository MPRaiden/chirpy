package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
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
	api.HandleFunc("/reset", apiCfg.handlerReset)
	api.Post("/validate_chirp", validateChirp)
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
