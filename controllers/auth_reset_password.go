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

func ResetPassword(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPut {
            utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
        }

        var req dto.ResetPasswordRequest
        if err := json.NewDecoder(r.Body).Decode(&req); 
		err != nil {
            utils.Error(w, http.StatusBadRequest, "Format JSON tidak valid")
			return
        }

        if req.Token == "" || req.NewPassword == "" || req.ConfirmNewPassword == "" {
            utils.Error(w, http.StatusBadRequest, "Token, password baru, dan konfirmasi password wajib diisi")
			return
        }

        if req.NewPassword != req.ConfirmNewPassword {
            utils.Error(w, http.StatusBadRequest, "Password baru dan konfirmasi password tidak sama")
			return
        }

        if !utils.IsValidPassword(req.NewPassword) {
            utils.Error(w, http.StatusBadRequest, "Password baru harus minimal 8 karakter, mengandung huruf besar, angka, dan simbol")
			return
        }

        userID, expiresAt, err := services.ValidatePasswordResetToken(db, req.Token)
        if err != nil {
            if err == sql.ErrNoRows {
                utils.Error(w, http.StatusBadRequest, "Token reset password tidak valid")
				return
            }
            panic(err)
        }

        if expiresAt.Before(time.Now()) {
            utils.Error(w, http.StatusBadRequest, "Token reset password telah kedaluwarsa")
			return
        }

        hashedPassword, err := utils.HashPassword(req.NewPassword)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal memproses password baru")
			return
        }

        err = services.ResetPassword(db, userID, hashedPassword)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal mereset password")
			return
        }

        err = services.DeletePasswordResetToken(db, req.Token)
        if err != nil {
            panic(err)
        }

        utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
            "success": true,
            "status":  http.StatusOK,
            "message": "Password berhasil direset.",
        })
    }
}
