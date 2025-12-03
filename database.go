package main

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./bookings.db")
	if err != nil {
		return err
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS bookings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		slot_time DATETIME NOT NULL UNIQUE,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_slot_time ON bookings(slot_time);

	CREATE TABLE IF NOT EXISTS blocked_slots (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		slot_time DATETIME NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_blocked_slot_time ON blocked_slots(slot_time);
	`

	_, err = db.Exec(createTableSQL)
	return err
}

func generateAvailableSlots() []AvailableSlot {
	slots := []AvailableSlot{}
	now := time.Now()

	// Generate slots for the next 90 days (approximately 3 months)
	for day := 0; day < 90; day++ {
		date := now.AddDate(0, 0, day)

		// Generate slots from 9 AM to 5 PM (9:00, 10:00, 11:00, 12:00, 13:00, 14:00, 15:00, 16:00)
		for hour := 9; hour <= 16; hour++ {
			slotTime := time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, time.UTC)

			// Only include future slots
			if slotTime.After(now) {
				slots = append(slots, AvailableSlot{
					SlotTime:  slotTime.Format(time.RFC3339),
					Available: true,
				})
			}
		}
	}

	return slots
}
