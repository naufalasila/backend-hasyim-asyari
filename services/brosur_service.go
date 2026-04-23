package services

import (
	"backend/models"
	"database/sql"
	"log"
)

func GetLatestBrosur(db *sql.DB) (*models.Brosur, error) {
	query := `SELECT id, filename, created_at FROM brosur ORDER BY id DESC LIMIT 1`
	var b models.Brosur
	err := db.QueryRow(query).Scan(&b.ID, &b.Filename, &b.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("ERROR DATABASE (GetLatestBrosur): %v", err)
		return nil, err
	}
	return &b, nil
}

func SaveBrosur(db *sql.DB, filename string) (int64, error) {
	// Delete all old entries, keep only the new one
	_, _ = db.Exec(`DELETE FROM brosur`)

	query := `INSERT INTO brosur (filename) VALUES (?)`
	result, err := db.Exec(query, filename)
	if err != nil {
		log.Printf("ERROR DATABASE (SaveBrosur): %v", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("ERROR DATABASE (LastInsertId brosur): %v", err)
		return 0, err
	}

	log.Printf("Brosur saved successfully: id=%d, filename=%s", id, filename)
	return id, nil
}
