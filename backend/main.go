package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

// Here's all the MAGIC
// Ta linijka mówi kompilatorowi Go:
// "Weź cały folder 'build' i wpakuj go do środka pliku '.exe' jako zmienną"
//
//go:embed build/*
var frontendFiles embed.FS

func main() {
	// 1. Ładowanie .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Ostrzeżenie: Nie znaleziono pliku .env, szukam w zmiennych systemowych.")
	}

	// Inicjalizacja bazy
	initDB()
	defer db.Close() // Zamknij bazę przy wyłączeniu aplikacji

	mux := http.NewServeMux()

	// --- ENDPOINTY API ---

	// 2. API
	mux.HandleFunc("POST /api/signup", handleSignup)
	mux.HandleFunc("GET /api/admin/participants", handleListParticipants)
	mux.HandleFunc("POST /api/admin/draw", handleDrawAndSend)
	mux.HandleFunc("GET /api/my-status", handleMyStatus)
	mux.HandleFunc("POST /api/update-preferences", handleUpdatePreferences)
	mux.HandleFunc("GET /api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Dzień dobry, tu serwerek :)"))
	})

	// --- 3. Serwowanie frontendu ---
	frontendFS, err := fs.Sub(frontendFiles, "build")
	if err != nil {
		log.Fatal(err)
	}

	// Tworzymy handler plików
	fileServer := http.FileServer(http.FS(frontendFS))

	// Obsługa wszystkich innych ścieżek ("/")
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Sprawdzić czy API nie zostało wywołane przez pomyłke tutaj
		// (opcjonalnie, ale dobra praktyka)

		// Trich dla SPA (Single Page Application):
		// Jeśli użytkownik wejdzie na "/admin" albo "/status",
		// fizycznie takiego pliku nie ma na serwerze.
		// Musimy wtedy zaserwoawć "index.html", a Svelte w przeglądarce
		// zzobaczy URL i wyświetli odpowiedni widok.

		path := r.URL.Path
		// Sprawdź czy plik istnieje w naszym wirtualnym systemie plików
		_, err := frontendFS.Open(path[1:]) // usuwamy pierwszy slash

		if err != nil {
			// Jeśli plik nie istnieje (np. /admin), serwujemy index.html
			// To pozwala działać routingowi Svelte
			r.URL.Path = "/"
		}

		fileServer.ServeHTTP(w, r)
	})

	// Konfiguracja CORS
	// Na razie zostawiamy
	corsHandler := enableCORS(mux)

	log.Println("🎅 Serwer Mikołaja (full stack) startuje na porcie :8080...")
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
