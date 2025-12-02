package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// --- HANDLERY ---

func handleSignup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest

	// 1. Dekodowanie JSON z ciała żądania
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Nieprawidłowy format danych JSON")
		return
	}
	defer r.Body.Close()

	// 2. Prosta walidacja
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Email) == "" {
		respondError(w, http.StatusBadRequest, "Imię i email są wymagane")
		return
	}

	// 3. Zapis do bazy
	// `prepare` dla bezpieczeństwa
	stmt, err := db.Prepare("INSERT INTO participants(name, email, preferences) VALUES(?, ?, ?)")
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Błąd bazy danych")
		log.Println("Prepare error: ", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(req.Name, req.Email, req.Preferences)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constrint failed: participants.email") {
			respondError(w, http.StatusConflict, "Ten adres email jest już zapisany")
			return
		}
		respondError(w, http.StatusInternalServerError, "Nie udało się zapisać uczestnika")
		return
	}

	log.Printf("Nowy uczestnik zapisany: %s (%s)\n", req.Name, req.Email)
	respondJSON(w, http.StatusCreated, map[string]string{"message": "Zapisano pomyślnie"})
}

func handleListParticipants(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, email. preferences FROM participants ORDER BY id DESC")
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Błąd podczas pobierania danych")
		log.Println("Query error: ", err)
		return
	}
	defer rows.Close()

	participants := []Participant{}
	for rows.Next() {
		var p Participant
		if err := rows.Scan(&p.ID, &p.Name, &p.Email, &p.Preferences); err != nil {
			log.Println("Row scan error: ", err)
			continue
		}
		participants = append(participants, p)
	}

	respondJSON(w, http.StatusOK, participants)
}

// TODO: dodać weryfikację hasła admina
func handleDrawAndSend(w http.ResponseWriter, r *http.Request) {
	// 1. Pobierz wszystkich uczestników z bazy
	// (Używamy tej samej logiki co w handleListParticipants, ale wewnętrznie)
	rows, err := db.Query("SELECT name, email, preferences FROM participants")
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Błąd bazy danych przed losowaniem")
		return
	}
	defer rows.Close()

	// Używamy struktury tymaczasowej, bo ID nie jest nam tu potrzebne
	type Giver struct {
		Name        string
		Email       string
		Preferences string
	}
	var givers []Giver
	for rows.Next() {
		var g Giver
		rows.Scan(&g.Name, &g.Email, &g.Preferences)
		givers = append(givers, g)
	}

	count := len(givers)
	if count < 2 {
		respondError(w, http.StatusBadRequest, "Za mało uczestników, aby przeprowadzić losowanie")
		return
	}

	log.Println("Rozpoczynam losowanie dla ", count, " osób...")

	// 2. Algorytm losowania

	rand.Seed(time.Now().UnixNano())

	rand.Shuffle(count, func(i, j int) {
		givers[i], givers[j] = givers[j], givers[i]
	})

	emailErrors := 0

	for i := 0; i < count; i++ {
		giver := givers[i]
		receiver := givers[(i+1)%count]

		fmt.Printf("[LOSOWANIE] %s kupuje prezent dla -> %s\n", giver.Name, receiver.Name)

		// 3. Wysyłka maila (placeholder)
		// Tutaj wywołamy funkcję z pliku mailer.go, na razie symulujemy
		err := sendSecretSantaEmail(giver.Email, giver.Name, receiver.Name, receiver.Preferences)
		if err != nil {
			log.Printf("Błąd wysyłki do %s: %v \n", giver.Email, err)
			emailErrors++
		}
	}

	responseMsg := fmt.Sprintf("Losowanie zakończone. Wysłano %d maili.", count-emailErrors)
	if emailErrors > 0 {
		responseMsg += fmt.Sprintf(" Uwaga: %d maili nie dotarło (sprawdź logi).", emailErrors)
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": responseMsg})
}

// --- FUNCKJĘ POMOCNICZE DO ODPOWIEDZI HTTP ---

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Błąd serwera podczas kodowania JSON"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
