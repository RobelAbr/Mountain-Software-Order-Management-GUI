package main

import (
	"database/sql"
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	// Verbindung zur Datenbank herstellen
	var err error
	db, err = sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/mountain_software")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Fyne App initialisieren
	a := app.New()
	w := a.NewWindow("Mountain Software Auftragsübersicht")

	// Login GUI erstellen
	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("Email")
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Passwort")
	loginButton := widget.NewButton("Login", func() {
		email := emailEntry.Text
		password := passwordEntry.Text
		if authenticateUser(email, password) {
			showOrderWindow(a)
		} else {
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Fehler",
				Content: "Ungültige Email oder Passwort",
			})
		}
	})

	// Login Layout
	w.SetContent(container.NewVBox(
		emailEntry,
		passwordEntry,
		loginButton,
	))

	w.ShowAndRun()
}

// Benutzer authentifizieren
func authenticateUser(email, password string) bool {
	var dbPassword string
	err := db.QueryRow("SELECT password FROM customers WHERE email = ?", email).Scan(&dbPassword)
	if err != nil {
		return false
	}
	return password == dbPassword
}

// Fenster zum Anzeigen der Aufträge
func showOrderWindow(a fyne.App) {
	w := a.NewWindow("Auftragsstatus einsehen")
	customerIdEntry := widget.NewEntry()
	customerIdEntry.SetPlaceHolder("Kunden ID")
	statusLabel := widget.NewLabel("")

	checkButton := widget.NewButton("Aufträge anzeigen", func() {
		customerId := customerIdEntry.Text
		orders, err := getOrders(customerId)
		if err != nil {
			statusLabel.SetText(fmt.Sprintf("Fehler: %v", err))
			return
		}
		statusLabel.SetText(orders)
	})

	w.SetContent(container.NewVBox(
		customerIdEntry,
		checkButton,
		statusLabel,
	))

	w.Show()
}

// Aufträge aus der Datenbank holen
func getOrders(customerId string) (string, error) {
	rows, err := db.Query("SELECT auftragsnummer FROM orders WHERE customer_id = ?", customerId)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var orders string
	for rows.Next() {
		var auftragsnummer string
		if err := rows.Scan(&auftragsnummer); err != nil {
			return "", err
		}
		orders += fmt.Sprintf("Auftragsnummer: %s\n", auftragsnummer)
	}
	return orders, nil
}
