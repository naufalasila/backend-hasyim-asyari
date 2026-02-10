package services

import (
	"backend/models"
	"database/sql"
	"log"
)

// GetAllBerita retrieves all berita records ordered by date (newest first)
func GetAllBerita(db *sql.DB) ([]models.Berita, error) {
	query := `
		SELECT id, tema, judul, COALESCE(gambar, ''), isi, tanggal, created_at
		FROM berita
		ORDER BY tanggal DESC, created_at DESC
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("ERROR DATABASE (GetAllBerita): %v", err)
		return nil, err
	}
	defer rows.Close()

	var beritaList []models.Berita
	for rows.Next() {
		var b models.Berita
		err := rows.Scan(&b.ID, &b.Tema, &b.Judul, &b.Gambar, &b.Isi, &b.Tanggal, &b.CreatedAt)
		if err != nil {
			log.Printf("ERROR DATABASE (Scan row): %v", err)
			return nil, err
		}
		beritaList = append(beritaList, b)
	}
	return beritaList, nil
}

// CreateBerita inserts a new berita record into the database
func CreateBerita(db *sql.DB, b models.Berita) (int64, error) {
	query := `
		INSERT INTO berita (tema, judul, gambar, isi, tanggal)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := db.Exec(query, b.Tema, b.Judul, b.Gambar, b.Isi, b.Tanggal)
	if err != nil {
		log.Printf("ERROR DATABASE (CreateBerita): %v", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("ERROR DATABASE (LastInsertId): %v", err)
		return 0, err
	}

	log.Printf("Berita created successfully: id=%d, judul=%s", id, b.Judul)
	return id, nil
}

// DeleteBerita removes a berita record from the database
func DeleteBerita(db *sql.DB, id int) error {
	query := `DELETE FROM berita WHERE id = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		log.Printf("ERROR DATABASE (DeleteBerita): %v", err)
	}
	return err
}

// GetBeritaByID retrieves a single berita by ID
func GetBeritaByID(db *sql.DB, id int) (*models.Berita, error) {
	query := `
		SELECT id, tema, judul, COALESCE(gambar, ''), isi, tanggal, created_at
		FROM berita
		WHERE id = ?
	`
	var b models.Berita
	err := db.QueryRow(query, id).Scan(&b.ID, &b.Tema, &b.Judul, &b.Gambar, &b.Isi, &b.Tanggal, &b.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("ERROR DATABASE (GetBeritaByID): %v", err)
		return nil, err
	}
	return &b, nil
}
