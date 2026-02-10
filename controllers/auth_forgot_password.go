package controllers

import (
	"backend/dto"
	"backend/services"
	"backend/utils"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

func ForgotPassword(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
            return
        }

        var req dto.ForgotPasswordRequest
        if err := json.NewDecoder(r.Body).Decode(&req); 
		err != nil {
            utils.Error(w, http.StatusBadRequest, "Format JSON tidak valid")
            return
        }

        if req.Email == "" {
			utils.Error(w, http.StatusBadRequest, "Email wajib diisi")
			return
		}

        userID, err := services.GetUserByEmail(db, req.Email)
        if err != nil {
            if err == sql.ErrNoRows {
                utils.Error(w, http.StatusNotFound, "Email tidak terdaftar")
				return
            }
            utils.Error(w, http.StatusInternalServerError, "Terjadi kesalahan pada server")
			return
        }

        token := utils.GenerateRandomToken(32)
        expiresAt := time.Now().Add(1 * time.Hour)

        if err := services.GeneratePasswordResetToken(db, userID, token, expiresAt); 
		err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal membuat token reset")
			return
        }

        if err := utils.SendResetPasswordEmail(req.Email, token); 
		err != nil {
            panic(err)
        }

        utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
            "success": true,
            "status":  http.StatusOK,
            "message": "Link reset password sudah dikirim ke email.",
        })
    }
}
