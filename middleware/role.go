// middleware/role.go
package middleware

import (
    "net/http"
    "backend/utils"
)

func Role(requiredRole string) func(http.HandlerFunc) http.HandlerFunc {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            claims, ok := r.Context().Value(UserContextKey).(*utils.Claims)
            if !ok {
                utils.Error(w, http.StatusUnauthorized, "Akses ditolak: pengguna tidak terautentikasi")
				return
            }

            if claims.Role != requiredRole {
                utils.Error(w, http.StatusForbidden, "Akses ditolak: peran tidak sesuai")
				return
            }

            next(w, r)
        }
    }
}