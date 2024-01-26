package main

import (
	"net/http"
	"log"
)

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Printf("Error writing response: %v", err)
		}
}

func main() {
    mux := http.NewServeMux()

    fsPublic := http.FileServer(http.Dir("public/"))
    // Serves index.html at the route '/app' and everything else in the 'public' folder under '/app/*'
    mux.Handle("/app/", http.StripPrefix("/app/", fsPublic))

    fsAssets := http.FileServer(http.Dir("assets/"))
    // Serves assets under the route '/app/assets/*'
    mux.Handle("/app/assets/", http.StripPrefix("/app/assets/", fsAssets))

    mux.HandleFunc("/healthz", healthzHandler)
    corsMux := middlewareCors(mux)

    err := http.ListenAndServe("localhost:8080", corsMux)
    if err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}

