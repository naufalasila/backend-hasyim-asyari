package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// GetAllBerita handles GET /api/berita - public endpoint
func GetAllBerita(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
		}

		beritaList, err := services.GetAllBerita(db)
		if err != nil {
			log.Printf("ERROR: Gagal mengambil data berita: %v", err)
			utils.Error(w, http.StatusInternalServerError, "Gagal mengambil data berita")
			return
		}

		utils.Success(w, http.StatusOK, "Data berita berhasil diambil", beritaList)
	}
}

// CreateBerita handles POST /api/admin/berita with multipart form data
func CreateBerita(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
		}

		// Parse multipart form (max 10MB for images)
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			log.Printf("ERROR: Gagal parse multipart form: %v", err)
			utils.Error(w, http.StatusBadRequest, "Gagal memproses form data")
			return
		}

		// Extract text fields
		berita := models.Berita{
			Tema:    r.FormValue("tema"),
			Judul:   r.FormValue("judul"),
			Isi:     r.FormValue("isi"),
			Tanggal: r.FormValue("tanggal"),
		}

		// Validate required fields
		if berita.Judul == "" || berita.Isi == "" || berita.Tanggal == "" {
			utils.Error(w, http.StatusBadRequest, "Judul, isi, dan tanggal wajib diisi")
			return
		}

		// Handle image upload
		if file, header, err := r.FormFile("gambar"); err == nil {
			filename, uploadErr := utils.UploadFoto(file, header, utils.BeritaImagePath)
			if uploadErr != nil {
				log.Printf("ERROR: Gagal upload gambar berita: %v", uploadErr)
				utils.Error(w, http.StatusInternalServerError, "Gagal upload gambar")
				return
			}
			berita.Gambar = filename
			log.Printf("Gambar berita uploaded: %s", filename)
		}

		// Save to database
		id, err := services.CreateBerita(db, berita)
		if err != nil {
			log.Printf("ERROR: Gagal menyimpan berita ke database: %v", err)
			utils.Error(w, http.StatusInternalServerError, "Gagal menyimpan berita")
			return
		}

		log.Printf("Berita berhasil dibuat: id=%d, judul=%s", id, berita.Judul)
		utils.Success(w, http.StatusCreated, "Berita berhasil dibuat", map[string]interface{}{
			"id": id,
		})
	}
}

// DeleteBerita handles DELETE /api/admin/berita
func DeleteBerita(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
		}

		// Parse request body for ID
		var req struct {
			ID int `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			// Try to get ID from query parameter
			idStr := r.URL.Query().Get("id")
			if idStr == "" {
				utils.Error(w, http.StatusBadRequest, "ID berita wajib disertakan")
				return
			}
			id, err := strconv.Atoi(idStr)
			if err != nil {
				utils.Error(w, http.StatusBadRequest, "ID tidak valid")
				return
			}
			req.ID = id
		}

		if req.ID == 0 {
			utils.Error(w, http.StatusBadRequest, "ID berita tidak valid")
			return
		}

		// Get berita to delete the image file
		berita, err := services.GetBeritaByID(db, req.ID)
		if err != nil {
			log.Printf("ERROR: Gagal mengambil data berita: %v", err)
			utils.Error(w, http.StatusInternalServerError, "Gagal menghapus berita")
			return
		}

		if berita == nil {
			utils.Error(w, http.StatusNotFound, "Berita tidak ditemukan")
			return
		}

		// Delete the image file if exists
		if berita.Gambar != "" {
			utils.HapusFoto(utils.BeritaImagePath, berita.Gambar)
		}

		// Delete from database
		err = services.DeleteBerita(db, req.ID)
		if err != nil {
			log.Printf("ERROR: Gagal menghapus berita id=%d: %v", req.ID, err)
			utils.Error(w, http.StatusInternalServerError, "Gagal menghapus berita")
			return
		}

		log.Printf("Berita id=%d berhasil dihapus", req.ID)
		utils.Success(w, http.StatusOK, "Berita berhasil dihapus", nil)
	}
}
