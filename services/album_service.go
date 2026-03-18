package services

import (
	"backend/models"
	"database/sql"
)

func GetAlbumsByKategori(db *sql.DB, kategori string) ([]models.Album, error) {
	query := `SELECT id, judul, gambar, kategori, created_at FROM album WHERE kategori = ? ORDER BY created_at DESC`
	rows, err := db.Query(query, kategori)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var albums []models.Album
	for rows.Next() {
		var album models.Album
		if err := rows.Scan(&album.ID, &album.Judul, &album.Gambar, &album.Kategori, &album.CreatedAt); err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if albums == nil {
		albums = []models.Album{}
	}

	return albums, nil
}

func GetAlbumByID(db *sql.DB, id int) (*models.Album, error) {
	query := `SELECT id, judul, gambar, kategori, created_at FROM album WHERE id = ?`
	var album models.Album
	err := db.QueryRow(query, id).Scan(&album.ID, &album.Judul, &album.Gambar, &album.Kategori, &album.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &album, nil
}

func CreateAlbum(db *sql.DB, album models.Album) (int, error) {
	query := `INSERT INTO album (judul, gambar, kategori) VALUES (?, ?, ?)`
	result, err := db.Exec(query, album.Judul, album.Gambar, album.Kategori)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func DeleteAlbum(db *sql.DB, id int) error {
	query := `DELETE FROM album WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}
