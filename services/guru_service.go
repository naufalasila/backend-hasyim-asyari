package services

import (
	"backend/models"
	"database/sql"
	"log"
)

// GetGuruByJenjang retrieves all guru for a specific jenjang, ordered by urutan
func GetGuruByJenjang(db *sql.DB, jenjang string) ([]models.Guru, error) {
	query := `
		SELECT id, nama, COALESCE(jabatan, ''), COALESCE(mata_pelajaran, ''),
		       COALESCE(pendidikan, ''), COALESCE(foto, ''), jenjang, urutan, created_at
		FROM guru
		WHERE jenjang = ?
		ORDER BY urutan ASC, created_at ASC
	`
	rows, err := db.Query(query, jenjang)
	if err != nil {
		log.Printf("ERROR DATABASE (GetGuruByJenjang): %v", err)
		return nil, err
	}
	defer rows.Close()

	var guruList []models.Guru
	for rows.Next() {
		var g models.Guru
		err := rows.Scan(&g.ID, &g.Nama, &g.Jabatan, &g.MataPelajaran,
			&g.Pendidikan, &g.Foto, &g.Jenjang, &g.Urutan, &g.CreatedAt)
		if err != nil {
			log.Printf("ERROR DATABASE (Scan guru row): %v", err)
			return nil, err
		}
		guruList = append(guruList, g)
	}
	return guruList, nil
}

// CreateGuru inserts a new guru record
func CreateGuru(db *sql.DB, g models.Guru) (int64, error) {
	query := `
		INSERT INTO guru (nama, jabatan, mata_pelajaran, pendidikan, foto, jenjang, urutan)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	result, err := db.Exec(query, g.Nama, g.Jabatan, g.MataPelajaran,
		g.Pendidikan, g.Foto, g.Jenjang, g.Urutan)
	if err != nil {
		log.Printf("ERROR DATABASE (CreateGuru): %v", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("ERROR DATABASE (LastInsertId guru): %v", err)
		return 0, err
	}

	log.Printf("Guru created successfully: id=%d, nama=%s", id, g.Nama)
	return id, nil
}

// DeleteGuru removes a guru record
func DeleteGuru(db *sql.DB, id int) error {
	query := `DELETE FROM guru WHERE id = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		log.Printf("ERROR DATABASE (DeleteGuru): %v", err)
	}
	return err
}

// UpdateGuru updates an existing guru record
func UpdateGuru(db *sql.DB, g models.Guru) error {
	query := `
		UPDATE guru SET nama = ?, jabatan = ?, mata_pelajaran = ?, pendidikan = ?,
		       foto = ?, jenjang = ?, urutan = ?
		WHERE id = ?
	`
	_, err := db.Exec(query, g.Nama, g.Jabatan, g.MataPelajaran,
		g.Pendidikan, g.Foto, g.Jenjang, g.Urutan, g.ID)
	if err != nil {
		log.Printf("ERROR DATABASE (UpdateGuru): %v", err)
		return err
	}
	log.Printf("Guru updated successfully: id=%d, nama=%s", g.ID, g.Nama)
	return nil
}

// GetGuruByID retrieves a single guru by ID
func GetGuruByID(db *sql.DB, id int) (*models.Guru, error) {
	query := `
		SELECT id, nama, COALESCE(jabatan, ''), COALESCE(mata_pelajaran, ''),
		       COALESCE(pendidikan, ''), COALESCE(foto, ''), jenjang, urutan, created_at
		FROM guru
		WHERE id = ?
	`
	var g models.Guru
	err := db.QueryRow(query, id).Scan(&g.ID, &g.Nama, &g.Jabatan, &g.MataPelajaran,
		&g.Pendidikan, &g.Foto, &g.Jenjang, &g.Urutan, &g.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("ERROR DATABASE (GetGuruByID): %v", err)
		return nil, err
	}
	return &g, nil
}
