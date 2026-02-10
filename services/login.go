package services

import (
	"backend/dto"
	"backend/utils"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func Login(db *sql.DB, req dto.LoginRequest) (string, error) {
	var id int
	var username, hashedPassword string

	// 1. Cari user di tabel admin
	query := "SELECT id, username, password FROM admin WHERE username = ? LIMIT 1"
	err := db.QueryRow(query, req.Username).Scan(&id, &username, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("username atau password salah")
		}
		return "", err
	}

	// 2. Bandingkan password (Bcrypt)
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		return "", errors.New("username atau password salah")
	}

	// 3. Generate JWT Token
	// Catatan: Karena tabel admin simpel, kita kirim string kosong untuk FullName/ProfilePicture
	token := utils.GenerateToken(id, username, username, "admin", "")

	return token, nil
}