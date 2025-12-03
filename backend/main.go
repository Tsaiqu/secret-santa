package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Ostrzeżenie: Nie znaleziono pliku .env, szukam w zmiennych systemowych.")
	}

	// Inicjalizacja bazy
	initDB()
	defer db.Close() // Zamknij bazę przy wyłączeniu aplikacji

	mux := http.NewServeMux()

	// --- ENDPOINTY API ---

	mux.HandleFunc("POST /api/signup", handleSignup)

	mux.HandleFunc("GET /api/admin/participants", handleListParticipants)

	mux.HandleFunc("POST /api/admin/draw", handleDrawAndSend)

	mux.HandleFunc("GET /api/my-status", handleMyStatus)

	mux.HandleFunc("GET /api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Dzień dobry, tu serwerek :)"))
	})

	corsHandler := enableCORS(mux)

	log.Println("Serwer Mikołaja startuje na porcie :8080...")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
