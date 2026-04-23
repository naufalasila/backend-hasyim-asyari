package services

import (
	"backend/models"
	"database/sql"
	"log"
)

func GetAllPembina(db *sql.DB) ([]models.Pembina, error) {
	query := `
		SELECT id, nama, COALESCE(jabatan, ''), COALESCE(pendidikan, ''),
		       COALESCE(foto, ''), urutan, created_at
		FROM pembina
		ORDER BY urutan ASC, created_at ASC
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("ERROR DATABASE (GetAllPembina): %v", err)
		return nil, err
	}
	defer rows.Close()

	var list []models.Pembina
	for rows.Next() {
		var p models.Pembina
		err := rows.Scan(&p.ID, &p.Nama, &p.Jabatan, &p.Pendidikan,
			&p.Foto, &p.Urutan, &p.CreatedAt)
		if err != nil {
			log.Printf("ERROR DATABASE (Scan pembina row): %v", err)
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func CreatePembina(db *sql.DB, p models.Pembina) (int64, error) {
	query := `
		INSERT INTO pembina (nama, jabatan, pendidikan, foto, urutan)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := db.Exec(query, p.Nama, p.Jabatan, p.Pendidikan, p.Foto, p.Urutan)
	if err != nil {
		log.Printf("ERROR DATABASE (CreatePembina): %v", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("ERROR DATABASE (LastInsertId pembina): %v", err)
		return 0, err
	}

	log.Printf("Pembina created successfully: id=%d, nama=%s", id, p.Nama)
	return id, nil
}

func DeletePembina(db *sql.DB, id int) error {
	query := `DELETE FROM pembina WHERE id = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		log.Printf("ERROR DATABASE (DeletePembina): %v", err)
	}
	return err
}

func UpdatePembina(db *sql.DB, p models.Pembina) error {
	query := `
		UPDATE pembina SET nama = ?, jabatan = ?, pendidikan = ?,
		       foto = ?, urutan = ?
		WHERE id = ?
	`
	_, err := db.Exec(query, p.Nama, p.Jabatan, p.Pendidikan,
		p.Foto, p.Urutan, p.ID)
	if err != nil {
		log.Printf("ERROR DATABASE (UpdatePembina): %v", err)
		return err
	}
	log.Printf("Pembina updated successfully: id=%d, nama=%s", p.ID, p.Nama)
	return nil
}

func GetPembinaByID(db *sql.DB, id int) (*models.Pembina, error) {
	query := `
		SELECT id, nama, COALESCE(jabatan, ''), COALESCE(pendidikan, ''),
		       COALESCE(foto, ''), urutan, created_at
		FROM pembina
		WHERE id = ?
	`
	var p models.Pembina
	err := db.QueryRow(query, id).Scan(&p.ID, &p.Nama, &p.Jabatan,
		&p.Pendidikan, &p.Foto, &p.Urutan, &p.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("ERROR DATABASE (GetPembinaByID): %v", err)
		return nil, err
	}
	return &p, nil
}
