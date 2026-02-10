package services

import (
    "database/sql"
    "time"
)

func GenerateVerificationToken(db *sql.DB, userID int, token string, expiresAt time.Time) error {
    _, err := db.Exec(
        "INSERT INTO email_verification_tokens (user_id, token, expires_at) VALUES (?, ?, ?)",
        userID, token, expiresAt,
    )
    return err
}

func ValidateVerificationToken(db *sql.DB, token string) (int, time.Time, error) {
    var userID int
    var expiresAt time.Time

    err := db.QueryRow(
        "SELECT user_id, expires_at FROM email_verification_tokens WHERE token = ?",
        token,
    ).Scan(&userID, &expiresAt)

    return userID, expiresAt, err
}

func VerifyEmail(db *sql.DB, userID int) error {
    _, err := db.Exec(
        "UPDATE users SET is_verified = 1, updated_at = CURRENT_TIMESTAMP WHERE id_user = ?",
        userID,
    )
    return err
}

func DeleteVerificationToken(db *sql.DB, token string) error {
    _, err := db.Exec("DELETE FROM email_verification_tokens WHERE token = ?", token)
    return err
}
