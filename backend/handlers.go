package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
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
	token := uuid.New().String()

	// `prepare` dla bezpieczeństwa
	stmt, err := db.Prepare("INSERT INTO participants(name, email, preferences, token) VALUES(?, ?, ?, ?)")
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Błąd bazy danych")
		log.Println("Prepare error: ", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(req.Name, req.Email, req.Preferences, token)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constrint failed: participants.email") {
			respondError(w, http.StatusConflict, "Ten adres email jest już zapisany")
			return
		}
		respondError(w, http.StatusInternalServerError, "Nie udało się zapisać uczestnika")
		return
	}

	log.Printf("Nowy uczestnik zapisany: %s (%s)\nToken: %s", req.Name, req.Email, token)
	respondJSON(w, http.StatusCreated, map[string]string{
		"message": "Zapisano! Sprawdź maila po link",
		"token":   token,
	})
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
// TODO: update początku funkcji (wczytyawania z bazy) - nowy model danych w DB
func handleDrawAndSend(w http.ResponseWriter, r *http.Request) {
	// 1. Pobierz wszystkich uczestników z bazy
	// (Używamy tej samej logiki co w handleListParticipants, ale wewnętrznie)
	rows, err := db.Query("SELECT id, name, email, preferences FROM participants")
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
		ID          int
	}
	var givers []Giver
	for rows.Next() {
		var g Giver
		if err := rows.Scan(&g.ID, &g.Name, &g.Email, &g.Preferences); err != nil {
			log.Println("Row scan error: ", err)
			continue
		}
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

	// --- ZAPISYWANIE DO BAZY ---

	// Rozpoczynamy transakcję SQL (żeby zapisać wszystko albo nic xd)
	tx, err := db.Begin()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Błąd transakcji")
		return
	}

	for i := 0; i < count; i++ {
		giver := givers[i] // Uwaga: musisz pobierać z bazy pełne obiekty Participant (z ID!), nie uproszczone
		receiver := givers[(i+1)%count]

		// 1. Zapisz w bazie, że Giver wylosował Receivera
		_, err := tx.Exec("UPDATE participants SET target_id = ? WHERE id = ?", receiver.ID, giver.ID)
		if err != nil {
			tx.Rollback()
			log.Println("Błąd zapisu pary: ", err)
			respondError(w, http.StatusInternalServerError, "Błąd zapisu losowania")
			return
		}
	}

	if err := tx.Commit(); err != nil {
		respondError(w, http.StatusInternalServerError, "Błąd zatwierdzania transakcji")
		return
	}

	// 3. Wysyłka maila
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

// Sprawdzanie statusu (Logowanie magicznym linkiem)
// GET /api/my-status?token=xyz
func handleMyStatus(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		respondError(w, http.StatusUnauthorized, "Brak tokenu")
		return
	}

	var me Participant
	var targetID sql.NullInt64 // Używamy NullInt64 bo może być NULL przed losowaniem

	// 1. Znajdź mnie po tokenie
	err := db.QueryRow("SELECT id, name, email, preferences, target_id FROM participants WHERE token = ?", token).Scan(&me.ID, &me.Name, &me.Email, &me.Preferences, &targetID)

	if err == sql.ErrNoRows {
		respondError(w, http.StatusUnauthorized, "Nieprawidłowy token")
		return
	} else if err != nil {
		respondError(w, http.StatusInternalServerError, "Błąd bazy")
		return
	}

	response := MyStatusResponse{
		Me:         me,
		IsDrawDone: false,
	}

	// 2. Jeśli mam przypisaną osobę (target_id nie jest null), pobierz jej dane
	if targetID.Valid {
		response.IsDrawDone = true
		var targetName, targetPrefs string
		err := db.QueryRow("SELECT name, preferences FROM participants WHERE id = ?", targetID.Int64).Scan(&targetName, &targetPrefs)
		if err == nil {
			response.TargetName = targetName
			response.TargetPrefs = targetPrefs
		}
	}

	respondJSON(w, http.StatusOK, response)
}

func handleUpdatePreferences(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		respondError(w, http.StatusUnauthorized, "Brak tokenu")
		return
	}

	var req struct {
		Preferences string `json:"preferences"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Nieprawidłowy JSON")
		return
	}
	defer r.Body.Close()

	if strings.TrimSpace(req.Preferences) == "" {
		respondError(w, http.StatusBadRequest, "Preferencje nie mogą być puste")
		return
	}

	// Sprawdź czy token istnieje i zaktualizuj
	res, err := db.Exec("UPDATE participants SET preferences = ? WHERE token = ?", req.Preferences, token)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Błąd bazy danych")
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		respondError(w, http.StatusUnauthorized, "Nieprawidłowy token")
		return
	}

	// Sprawdź czy ktoś mnie wylosował i wyślij powiadomeinie
	var myName string
	var myID int
	err = db.QueryRow("SELECT id, name FROM participants WHERE token = ?", token).Scan(&myID, &myName)
	if err == nil {
		// Szukamy Mikołaja (osoby, która ma target_id = myID)
		var giverEmail, giverName string
		err = db.QueryRow("SELECT email, name FROM participants WHERE target_id = ?", myID).Scan(&giverEmail, &giverName)
		if err == nil {
			// Mamy Mikołaja, wysyłamy powiadomienie
			go func() {
				if err := sendPreferencesUpdateEmail(giverEmail, giverName, myName, req.Preferences); err != nil {
					log.Printf("Błąd wysyłki powiadomienia o aktualizacji preferencji do %s: %v\n", giverEmail, err)
				}
			}()
		}
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Preferencje zaktualizowane!"})
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
