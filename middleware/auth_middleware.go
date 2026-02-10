// middleware/auth.go
package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v5"
    "backend/utils"
)

type contextKey string
var UserContextKey contextKey = "user_claims"

func Auth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")

        if authHeader == "" {
            utils.Error(w, http.StatusUnauthorized, "Header Authorization diperlukan")
            return
        }

        if !strings.HasPrefix(authHeader, "Bearer ") {
            utils.Error(w, http.StatusUnauthorized, "Format token harus Bearer <token>")
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        claims := &utils.Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }
            return utils.Secret, nil
        })

        if err != nil {
            utils.Error(w, http.StatusUnauthorized, "Token tidak valid")
            return
        }

        if !token.Valid {
            utils.Error(w, http.StatusUnauthorized, "Token tidak valid")
            return
        }

        ctx := context.WithValue(r.Context(), UserContextKey, claims)
        r = r.WithContext(ctx)

        next(w, r)
    }
}
