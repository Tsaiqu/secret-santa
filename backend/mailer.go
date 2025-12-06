package main

import (
	"fmt"
	"net/smtp"
	"os"
)

const (
	smtpHost = "smtp.gmail.com"
	smtpPort = "587"
)

func sendSecretSantaEmail(recipientEmail, recipientName, targetName, targetPrefs string) error {
	// 1. Pobieramy dane nadawcy ze zmiennych środowiskowych
	senderEmail := os.Getenv("SANTA_EMAIL")
	senderPassword := os.Getenv("SANTA_PASSWORD")

	if senderEmail == "" || senderPassword == "" {
		return fmt.Errorf("Brak konfiguracji emaila (ustaw zmienne SANTA_EMAIL i SANTA_PASSWORD)")
	}

	subject := "🎅 Twój wynik Świątecznego Losowania!"
	body := fmt.Sprintf(`Ho Ho Ho %s!
		Świąteczne losowanie zostało zakończone.

		Twoja wylosowana osoba to: %s

		Oto co Mikołaj wie o jej preferencjach:
		---------------------------------------
		%s
		---------------------------------------

		Powodzenia w szukaniu prezentu!
		Twój Świąteczny Bot 🎄
		`, recipientName, targetName, targetPrefs)

	// Składanie wsyzstkiego w całość (Nagłówki + Treść)
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s", recipientEmail, subject, body))

	// 3. Uwierzytelnianie (Logowanie do gmaila)
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	// 4. Fizyczne wysłanie maila
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{recipientEmail}, msg)
	if err != nil {
		return err
	}

	return nil
}

func sendPreferencesUpdatedEmail(recipientEmail, recipientName, targetName, newPreferences string) error {
	senderEmail := os.Getenv("SANTA_EMAIL")
	senderPassword := os.Getenv("SANTA_PASSWORD")

	if senderEmail == "" || senderPassword == "" {
		return fmt.Errorf("Brak konfiguracji emaila (ustaw zmienne SANTA_EMAIL i SANTA_PASSWORD)")
	}

	subject := fmt.Sprintf("Aktualizacja: %s zmienił(a) swoje preferencje!", targetName)
	body := fmt.Sprintf(`Ho Ho Ho %s!

Twoja wylosowana osoba (%s) zaktualizowała swoje preferencje.
Oto co teraz pisze w liście do Mikołaja:

---------------------------------------
%s
---------------------------------------

Powodzenia!
Twój Świąteczny Bot 🎄
`, recipientName, targetName, newPreferences)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s", recipientEmail, subject, body))

	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{recipientEmail}, msg)
	if err != nil {
		return err
	}

	return nil
}
