package controllers

import (
	"backend/services"
	"backend/utils"
	"database/sql"
	"log"
	"net/http"
)

// GetBrosur returns the latest brosur info (public)
func GetBrosur(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
		}

		brosur, err := services.GetLatestBrosur(db)
		if err != nil {
			log.Printf("ERROR: Gagal mengambil data brosur: %v", err)
			utils.Error(w, http.StatusInternalServerError, "Gagal mengambil data brosur")
			return
		}

		if brosur == nil {
			utils.Success(w, http.StatusOK, "Belum ada brosur", nil)
			return
		}

		utils.Success(w, http.StatusOK, "Data brosur berhasil diambil", brosur)
	}
}

// UploadBrosur handles POST /api/admin/brosur (admin only)
func UploadBrosur(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
		}

		err := r.ParseMultipartForm(50 << 20) // 50MB max
		if err != nil {
			log.Printf("ERROR: Gagal parse multipart form: %v", err)
			utils.Error(w, http.StatusBadRequest, "Gagal memproses form data")
			return
		}

		file, header, err := r.FormFile("brosur")
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "File brosur wajib diupload")
			return
		}

		// Delete old brosur file if exists
		oldBrosur, _ := services.GetLatestBrosur(db)
		if oldBrosur != nil && oldBrosur.Filename != "" {
			utils.HapusFoto(utils.BrosurPath, oldBrosur.Filename)
		}

		filename, uploadErr := utils.UploadFoto(file, header, utils.BrosurPath)
		if uploadErr != nil {
			log.Printf("ERROR: Gagal upload brosur: %v", uploadErr)
			utils.Error(w, http.StatusInternalServerError, "Gagal upload brosur")
			return
		}

		id, err := services.SaveBrosur(db, filename)
		if err != nil {
			log.Printf("ERROR: Gagal menyimpan brosur ke database: %v", err)
			utils.Error(w, http.StatusInternalServerError, "Gagal menyimpan data brosur")
			return
		}

		log.Printf("Brosur berhasil diupload: id=%d, filename=%s", id, filename)
		utils.Success(w, http.StatusCreated, "Brosur berhasil diupload", map[string]interface{}{
			"id":       id,
			"filename": filename,
		})
	}
}
