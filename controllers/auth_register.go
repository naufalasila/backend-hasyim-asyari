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

func Register(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
            return
        }

        var req dto.RegisterRequest
        if err := json.NewDecoder(r.Body).Decode(&req); 
		err != nil {
            utils.Error(w, http.StatusBadRequest, "Format JSON tidak valid")
            return
        }

        if req.Username == "" || req.Email == "" || req.Password == "" || req.ConfirmPassword == "" {
            utils.Error(w, http.StatusBadRequest, "Semua field wajib diisi")
            return
        }

        _, _, _, _, _, err := services.GetUserByUsername(db, req.Username)
        if err == nil {
            utils.Error(w, http.StatusConflict, "Username sudah digunakan")
            return
        } else if err != sql.ErrNoRows {
            utils.Error(w, http.StatusInternalServerError, "Gagal memeriksa username")
            return
        }

        if !utils.IsValidUsername(req.Username) {
            utils.Error(w, http.StatusBadRequest, "Username tidak valid (minimal 3 karakter, huruf/angka/_)")
            return
        }

        _, err = services.GetUserByEmail(db, req.Email)
        if err == nil {
            utils.Error(w, http.StatusConflict, "Email sudah digunakan")
            return
        } else if err != sql.ErrNoRows {
            utils.Error(w, http.StatusInternalServerError, "Gagal memeriksa email")
            return
        }

        if !utils.IsValidEmail(req.Email) {
            utils.Error(w, http.StatusBadRequest, "Format email tidak valid")
            return
        }

        if req.Password != req.ConfirmPassword {
            utils.Error(w, http.StatusBadRequest, "Password dan konfirmasi tidak cocok")
            return
        }

        if !utils.IsValidPassword(req.Password) {
            utils.Error(w, http.StatusBadRequest, "Password harus minimal 8 karakter, mengandung huruf besar, angka, dan simbol")
            return
        }

        hashedPassword, err := utils.HashPassword(req.Password)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal memproses password")
			return
        }

        userID, err := services.Register(db, req.Username, req.Email, hashedPassword, "user")
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal mendaftarkan akun")
			return
        }

        verificationToken := utils.GenerateRandomToken(32)
        expiresAt := time.Now().Add(24 * time.Hour)

        err = services.GenerateVerificationToken(db, userID, verificationToken, expiresAt)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal membuat token verifikasi")
			return
        }

        err = utils.SendVerificationEmail(req.Email, verificationToken)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal mengirim email verifikasi. Pastikan email valid.")
            return
        }

        utils.JSONResponse(w, http.StatusCreated, map[string]interface{}{
            "success": true,
            "status":  http.StatusCreated,
            "message": "Akun berhasil dibuat. Silakan cek email untuk verifikasi.",
            "user": map[string]string{
                "username": req.Username,
                "email":    req.Email,
            },
        })
    }
}
