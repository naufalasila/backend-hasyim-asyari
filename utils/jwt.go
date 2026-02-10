// utils/jwt.go
package utils

import (
    "os"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

var Secret []byte

func InitJWTSecret() {
    secret := os.Getenv("JWT_SECRET")

    if secret == "" {
        panic("JWT_SECRET tidak ditemukan di .env")
    }
    
    // Minimal panjang secret untuk keamanan HS256 adalah 32 karakter
    if len(secret) < 32 {
        panic("JWT_SECRET terlalu pendek! Gunakan minimal 32 karakter.")
    }
    
    Secret = []byte(secret)
}

type Claims struct {
    IDUser          int    `json:"id_user"`
    Username        string `json:"username"`
    FullName        string `json:"full_name"`
    Role            string `json:"role"`
    ProfilePicture  string `json:"profile_picture"`
    jwt.RegisteredClaims
}

func GenerateToken(userID int, username, fullName, role, profilePicture string) string {
    claims := &Claims{
        IDUser:         userID,
        Username:       username,
        FullName:       fullName,
        Role:           role,
        ProfilePicture: profilePicture,
        RegisteredClaims: jwt.RegisteredClaims{
            // Token berlaku selama 24 jam
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString(Secret)
    if err != nil {
        // Gunakan panic hanya untuk error fatal saat setup, 
        // namun untuk flow produksi, sebaiknya fungsi ini mengembalikan error.
        panic("Gagal membuat token JWT: " + err.Error())
    }

    return signedToken
}