// controllers/profile.go
package controllers

import (
	"backend/dto"
	"backend/middleware"
	"backend/services"
	"backend/utils"
	"database/sql"
	"net/http"
	"path/filepath"
	"strings"
)

func GetProfile(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            utils.Error(w, http.StatusMethodNotAllowed, "Hanya GET yang diizinkan")
            return
        }
        claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
        if !ok {
            utils.Error(w, http.StatusUnauthorized, "Akses ditolak")
            return
        }

        userData, err := services.GetProfileData(db, claims.IDUser)
        if err != nil {
            if err == sql.ErrNoRows {
                utils.Error(w, http.StatusNotFound, "User tidak ditemukan")
            } else {
                utils.Error(w, http.StatusInternalServerError, "Gagal ambil data")
            }
            return
        }

        status, err := services.GetUserPendaftarStatus(db, claims.IDUser)
        if err != nil && err != sql.ErrNoRows {
            utils.Error(w, http.StatusInternalServerError, "Gagal cek status")
            return
        }

        statusKeanggotaan := "Bukan Calon Anggota"
        if status == "diterima" || status == "pending" {
            if status == "pending" {
                statusKeanggotaan = "Calon Anggota (Pending)"
            } else {
                statusKeanggotaan = "Calon Anggota"
            }
        }

        response := dto.UserProfileResponse{
            FullName:          userData.FullName,
            Email:             userData.Email,
            ProfilePicture:    userData.ProfilePicture,
            TanggalBergabung:  userData.CreatedAt.Format("2006-01-02"),
            StatusKeanggotaan: statusKeanggotaan,
        }

        utils.JSONResponse(w, http.StatusOK, response)
    }
}

func UpdateProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut && r.Method != http.MethodPost {
			utils.Error(w, http.StatusMethodNotAllowed, "Hanya PUT/POST yang diizinkan")
			return
		}

		claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		if !ok {
			utils.Error(w, http.StatusUnauthorized, "Akses ditolak")
			return
		}

		err := r.ParseMultipartForm(5 << 20)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Gagal parsing form")
			return
		}

		fullName := strings.TrimSpace(r.FormValue("full_name"))
		if fullName == "" {
			utils.Error(w, http.StatusBadRequest, "Nama lengkap wajib diisi")
			return
		}

		var newProfilePicture string

		file, header, err := r.FormFile("profile_picture")
		if err == nil {
			defer file.Close()

			ext := strings.ToLower(filepath.Ext(header.Filename))
			if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
				utils.Error(w, http.StatusBadRequest, "Format foto tidak didukung")
				return
			}

			uploadedFile, err := utils.UploadFoto(file, header, utils.ProfilePhotoPath)
			if err != nil {
				utils.Error(w, http.StatusInternalServerError, "Gagal upload foto")
				return
			}
			newProfilePicture = uploadedFile

			if claims.ProfilePicture != "" && claims.ProfilePicture != "default.jpg" {
				utils.HapusFoto(utils.ProfilePhotoPath, claims.ProfilePicture)
			}
		}

		err = services.UpdateUserProfile(db, claims.IDUser, fullName, newProfilePicture)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal update profil")
			return
		}

		claims.FullName = fullName
		if newProfilePicture != "" {
			claims.ProfilePicture = newProfilePicture
		}

		utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Profil berhasil diperbarui",
			"data": map[string]string{
				"full_name":       fullName,
				"profile_picture": newProfilePicture,
			},
		})
	}
}
