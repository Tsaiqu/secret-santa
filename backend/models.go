package main

// Participant reprezentuje osobę biorącą udział w losowaniu.
type Participant struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Preferences string `json:"preferences"`
}

// Struktura do odbierania danych z formularza zapisu
type SignupRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Preferences string `json:"preferences"`
}
