package main

import "time"

type Event struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatorID   int       `json:"creator_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type Participant struct {
	ID          int    `json:"id"`
	EventID     int    `json:"event_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Preferences string `json:"preferences"`
	Token       string `json:"token"`     // Klucz dostępu
	TargetID    *int   `json:"target_id"` // ID osoby, którą wylosowaliśmy (może być null przed losowaniem)
}

type EventItem struct {
	ID            int       `json:"id"`
	EventID       int       `json:"event_id"`
	ParticipantID int       `json:"participant_id"`
	ItemName      int       `json:"item_name"`
	Quantity      int       `json:"quantity"`
	CreatedAt     time.Time `json:"created_at"`
}

// Struktura do odbierania danych z formularza zapisu
type SignupRequest struct {
	EventID     int    `json:"event_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Preferences string `json:"preferences"`
}

type MyStatusResponse struct {
	Me          Participant `json:"me"`
	TargetName  string      `json:"target_name,omitempty"`
	TargetPrefs string      `json:"target_prefs,omitempty"`
	IsDrawDone  bool        `json:"is_draw_done"`
}
