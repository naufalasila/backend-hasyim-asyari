package services

import (
    "database/sql"
    "time"
)

func GeneratePasswordResetToken(db *sql.DB, userID int, token string, expiresAt time.Time) error {
    _, err := db.Exec(
        "INSERT INTO password_resets (user_id, token, expires_at) VALUES (?, ?, ?)",
        userID, token, expiresAt,
    )
    return err
}

func ValidatePasswordResetToken(db *sql.DB, token string) (int, time.Time, error) {
    var userID int
    var expiresAt time.Time

    err := db.QueryRow(
        "SELECT user_id, expires_at FROM password_resets WHERE token = ?",
        token,
    ).Scan(&userID, &expiresAt)

    return userID, expiresAt, err
}

func ResetPassword(db *sql.DB, userID int, newPasswordHash string) error {
    _, err := db.Exec(
        "UPDATE users SET password = ?, updated_at = CURRENT_TIMESTAMP WHERE id_user = ?",
        newPasswordHash, userID,
    )
    return err
}

func DeletePasswordResetToken(db *sql.DB, token string) error {
    _, err := db.Exec("DELETE FROM password_resets WHERE token = ?", token)
    return err
}
