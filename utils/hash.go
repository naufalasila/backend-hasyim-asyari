package utils

import (
    "crypto/rand"
    "encoding/hex"
    "golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
    hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashed), nil
}

func CheckPassword(password, hashed string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}

func GenerateRandomToken(length int) string {
    bytes := make([]byte, length)
    _, err := rand.Read(bytes)
    if err != nil {
        return ""
    }
    return hex.EncodeToString(bytes)
}