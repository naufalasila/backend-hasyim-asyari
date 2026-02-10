// services/profile.go
package services

import (
    "database/sql"
    "backend/models"
)

func GetProfileData(db *sql.DB, userID int) (*models.User, error) {
    var user models.User
    query := `
        SELECT id_user, full_name, email, profile_picture, created_at 
        FROM users 
        WHERE id_user = ?
    `
    err := db.QueryRow(query, userID).Scan(
        &user.IDUser,
        &user.FullName,
        &user.Email,
        &user.ProfilePicture,
        &user.CreatedAt,
    )
    return &user, err
}

func GetUserPendaftarStatus(db *sql.DB, userID int) (string, error) {
    var status string
    query := `SELECT status FROM pendaftar WHERE user_id = ? LIMIT 1`
    err := db.QueryRow(query, userID).Scan(&status)
    if err == sql.ErrNoRows {
        return "", nil
    }
    return status, err
}

func UpdateUserProfile(db *sql.DB, userID int, fullName, profilePicture string) error {
    query := `
        UPDATE users 
        SET full_name = ?, profile_picture = ?, updated_at = NOW() 
        WHERE id_user = ?
    `
    _, err := db.Exec(query, fullName, profilePicture, userID)
    return err
}