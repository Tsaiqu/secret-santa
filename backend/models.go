package main

// Participant reprezentuje osobę biorącą udział w losowaniu.
type Participant struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Preferences string `json:"preferences"`
	Token       string `json:"token"`     // Klucz dostępu
	TargetID    *int   `json:"target_id"` // ID osoby, którą wylosowaliśmy (może być null przed losowaniem)
}

// Struktura do odbierania danych z formularza zapisu
type SignupRequest struct {
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
