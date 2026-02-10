package services

import (
	"backend/models"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// GenerateRegistrasiID creates a unique registration ID like "REG-2026020001"
func GenerateRegistrasiID(db *sql.DB) (string, error) {
	// Get current year and month
	now := time.Now()
	prefix := fmt.Sprintf("REG-%d%02d", now.Year(), now.Month())

	// Count existing registrations this month to get the next sequence number
	var count int
	query := `SELECT COUNT(*) FROM pendaftar WHERE registrasi_id LIKE ?`
	err := db.QueryRow(query, prefix+"%").Scan(&count)
	if err != nil {
		log.Printf("ERROR DATABASE (GenerateRegistrasiID): %v", err)
		return "", err
	}

	// Generate the new ID
	newID := fmt.Sprintf("%s%04d", prefix, count+1)
	return newID, nil
}

// CreatePendaftar inserts a new pendaftar record into the database
func CreatePendaftar(db *sql.DB, p models.Pendaftar) (int64, string, error) {
	// Generate registrasi_id
	registrasiID, err := GenerateRegistrasiID(db)
	if err != nil {
		return 0, "", err
	}

	// Handle empty date - set to NULL or use current date
	var tanggalLahir interface{}
	if p.TanggalLahir == "" {
		tanggalLahir = nil
	} else {
		tanggalLahir = p.TanggalLahir // Frontend sends YYYY-MM-DD format
	}

	// Using correct column names: file_skl, file_kk, file_akte, file_pip
	query := `
		INSERT INTO pendaftar 
		(registrasi_id, nama_lengkap, nisn, tempat_lahir, tanggal_lahir, alamat, asal_sekolah, 
		nama_ayah, pekerjaan_ayah, nama_ibu, pekerjaan_ibu, no_hp_ortu,
		foto_profil, file_skl, file_kk, file_akte, file_pip, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'Menunggu Verifikasi')
	`

	log.Printf("Attempting to insert pendaftar: registrasi_id=%s, nama=%s, nisn=%s", registrasiID, p.NamaLengkap, p.NISN)

	result, err := db.Exec(query,
		registrasiID,
		p.NamaLengkap, p.NISN, p.TempatLahir, tanggalLahir, p.Alamat, p.AsalSekolah,
		p.NamaAyah, p.PekerjaanAyah, p.NamaIbu, p.PekerjaanIbu, p.NoHpOrtu,
		p.FotoProfil, p.FileSKL, p.FileKK, p.FileAkte, p.FilePIP,
	)
	if err != nil {
		log.Printf("ERROR DATABASE (CreatePendaftar): %v", err)
		return 0, "", err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("ERROR DATABASE (LastInsertId): %v", err)
		return 0, "", err
	}

	log.Printf("Pendaftar created successfully: id=%d, registrasi_id=%s", id, registrasiID)
	return id, registrasiID, nil
}

// GetAllPendaftar retrieves all pendaftar records from the database
func GetAllPendaftar(db *sql.DB) ([]models.Pendaftar, error) {
	// Using correct column names: file_skl, file_kk, file_akte, file_pip
	query := `
		SELECT id, registrasi_id, nama_lengkap, nisn, tempat_lahir, 
		COALESCE(tanggal_lahir, ''), alamat, asal_sekolah,
		nama_ayah, pekerjaan_ayah, nama_ibu, pekerjaan_ibu, COALESCE(no_hp_ortu, ''),
		COALESCE(foto_profil, ''), COALESCE(file_skl, ''), COALESCE(file_kk, ''), 
		COALESCE(file_akte, ''), COALESCE(file_pip, ''), status, created_at
		FROM pendaftar
		ORDER BY created_at DESC
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("ERROR DATABASE (GetAllPendaftar): %v", err)
		return nil, err
	}
	defer rows.Close()

	var pendaftarList []models.Pendaftar
	for rows.Next() {
		var p models.Pendaftar
		err := rows.Scan(
			&p.ID, &p.RegistrasiID, &p.NamaLengkap, &p.NISN, &p.TempatLahir, &p.TanggalLahir, &p.Alamat, &p.AsalSekolah,
			&p.NamaAyah, &p.PekerjaanAyah, &p.NamaIbu, &p.PekerjaanIbu, &p.NoHpOrtu,
			&p.FotoProfil, &p.FileSKL, &p.FileKK, &p.FileAkte, &p.FilePIP, &p.Status, &p.CreatedAt,
		)
		if err != nil {
			log.Printf("ERROR DATABASE (Scan row): %v", err)
			return nil, err
		}
		pendaftarList = append(pendaftarList, p)
	}
	return pendaftarList, nil
}

// UpdatePendaftarStatus updates the status of a pendaftar record
func UpdatePendaftarStatus(db *sql.DB, id int, status string) error {
	query := `UPDATE pendaftar SET status = ? WHERE id = ?`
	_, err := db.Exec(query, status, id)
	if err != nil {
		log.Printf("ERROR DATABASE (UpdatePendaftarStatus): %v", err)
	}
	return err
}
