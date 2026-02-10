package controllers

import (
    "backend/services"
    "backend/utils"
    "net/http"
    "time"
    "database/sql"
)

func VerifyEmail(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.URL.Query().Get("token")
        if token == "" {
            utils.Error(w, http.StatusBadRequest, "Token verifikasi wajib disertakan")
			return
        }

        userID, expiresAt, err := services.ValidateVerificationToken(db, token)
        if err != nil {
            if err == sql.ErrNoRows {
                utils.Error(w, http.StatusBadRequest, "Token verifikasi tidak valid")
				return
            }
            panic(err)
        }

        if expiresAt.Before(time.Now()) {
            utils.Error(w, http.StatusBadRequest, "Token verifikasi sudah kedaluwarsa")
			return
        }

        if err := services.VerifyEmail(db, userID); 
		err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal memverifikasi email")
            return
        }

        if err := services.DeleteVerificationToken(db, token); 
		err != nil {
            panic(err)
        }

        utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
            "success": true,
            "status":  http.StatusOK,
            "message": "Email berhasil diverifikasi.",
        })
    }
}
