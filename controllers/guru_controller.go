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

// GetGuru handles GET /api/guru - public endpoint, filter by ?jenjang=mts|ma
func GetGuru(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
		}

		jenjang := r.URL.Query().Get("jenjang")
		if jenjang == "" {
			jenjang = "mts"
		}

		guruList, err := services.GetGuruByJenjang(db, jenjang)
		if err != nil {
			log.Printf("ERROR: Gagal mengambil data guru: %v", err)
			utils.Error(w, http.StatusInternalServerError, "Gagal mengambil data guru")
			return
		}

		utils.Success(w, http.StatusOK, "Data guru berhasil diambil", guruList)
	}
}

// CreateGuru handles POST /api/admin/guru with multipart form data
func CreateGuru(db *sql.DB) http.HandlerFunc {
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

		// Parse urutan
		urutan := 0
		if urutanStr := r.FormValue("urutan"); urutanStr != "" {
			if val, err := strconv.Atoi(urutanStr); err == nil {
				urutan = val
			}
		}

		guru := models.Guru{
			Nama:          r.FormValue("nama"),
			Jabatan:       r.FormValue("jabatan"),
			MataPelajaran: r.FormValue("mata_pelajaran"),
			Pendidikan:    r.FormValue("pendidikan"),
			Jenjang:       r.FormValue("jenjang"),
			Urutan:        urutan,
		}

		// Validate required fields
		if guru.Nama == "" {
			utils.Error(w, http.StatusBadRequest, "Nama guru wajib diisi")
			return
		}
		if guru.Jenjang == "" {
			guru.Jenjang = "mts"
		}

		// Handle image upload
		if file, header, err := r.FormFile("foto"); err == nil {
			filename, uploadErr := utils.UploadFoto(file, header, utils.GuruImagePath)
			if uploadErr != nil {
				log.Printf("ERROR: Gagal upload foto guru: %v", uploadErr)
				utils.Error(w, http.StatusInternalServerError, "Gagal upload foto")
				return
			}
			guru.Foto = filename
			log.Printf("Foto guru uploaded: %s", filename)
		}

		// Save to database
		id, err := services.CreateGuru(db, guru)
		if err != nil {
			log.Printf("ERROR: Gagal menyimpan guru ke database: %v", err)
			utils.Error(w, http.StatusInternalServerError, "Gagal menyimpan data guru")
			return
		}

		log.Printf("Guru berhasil ditambahkan: id=%d, nama=%s", id, guru.Nama)
		utils.Success(w, http.StatusCreated, "Guru berhasil ditambahkan", map[string]interface{}{
			"id": id,
		})
	}
}

// UpdateGuru handles PUT /api/admin/guru with multipart form data
func UpdateGuru(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
		}

		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			log.Printf("ERROR: Gagal parse multipart form: %v", err)
			utils.Error(w, http.StatusBadRequest, "Gagal memproses form data")
			return
		}

		// Parse ID
		idStr := r.FormValue("id")
		if idStr == "" {
			utils.Error(w, http.StatusBadRequest, "ID guru wajib disertakan")
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "ID tidak valid")
			return
		}

		// Get existing guru
		existing, err := services.GetGuruByID(db, id)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal mengambil data guru")
			return
		}
		if existing == nil {
			utils.Error(w, http.StatusNotFound, "Guru tidak ditemukan")
			return
		}

		// Parse urutan
		urutan := existing.Urutan
		if urutanStr := r.FormValue("urutan"); urutanStr != "" {
			if val, err := strconv.Atoi(urutanStr); err == nil {
				urutan = val
			}
		}

		guru := models.Guru{
			ID:            id,
			Nama:          r.FormValue("nama"),
			Jabatan:       r.FormValue("jabatan"),
			MataPelajaran: r.FormValue("mata_pelajaran"),
			Pendidikan:    r.FormValue("pendidikan"),
			Jenjang:       r.FormValue("jenjang"),
			Urutan:        urutan,
			Foto:          existing.Foto, // keep old photo by default
		}

		if guru.Nama == "" {
			utils.Error(w, http.StatusBadRequest, "Nama guru wajib diisi")
			return
		}
		if guru.Jenjang == "" {
			guru.Jenjang = existing.Jenjang
		}

		// Handle new image upload
		if file, header, err := r.FormFile("foto"); err == nil {
			filename, uploadErr := utils.UploadFoto(file, header, utils.GuruImagePath)
			if uploadErr != nil {
				log.Printf("ERROR: Gagal upload foto guru: %v", uploadErr)
				utils.Error(w, http.StatusInternalServerError, "Gagal upload foto")
				return
			}
			// Delete old photo
			if existing.Foto != "" {
				utils.HapusFoto(utils.GuruImagePath, existing.Foto)
			}
			guru.Foto = filename
		}

		err = services.UpdateGuru(db, guru)
		if err != nil {
			log.Printf("ERROR: Gagal mengupdate guru id=%d: %v", id, err)
			utils.Error(w, http.StatusInternalServerError, "Gagal mengupdate data guru")
			return
		}

		log.Printf("Guru id=%d berhasil diupdate", id)
		utils.Success(w, http.StatusOK, "Guru berhasil diupdate", nil)
	}
}

// DeleteGuru handles DELETE /api/admin/guru
func DeleteGuru(db *sql.DB) http.HandlerFunc {
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
			// Try query parameter
			idStr := r.URL.Query().Get("id")
			if idStr == "" {
				utils.Error(w, http.StatusBadRequest, "ID guru wajib disertakan")
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
			utils.Error(w, http.StatusBadRequest, "ID guru tidak valid")
			return
		}

		// Get guru to delete the image file
		guru, err := services.GetGuruByID(db, req.ID)
		if err != nil {
			log.Printf("ERROR: Gagal mengambil data guru: %v", err)
			utils.Error(w, http.StatusInternalServerError, "Gagal menghapus guru")
			return
		}

		if guru == nil {
			utils.Error(w, http.StatusNotFound, "Guru tidak ditemukan")
			return
		}

		// Delete the image file if exists
		if guru.Foto != "" {
			utils.HapusFoto(utils.GuruImagePath, guru.Foto)
		}

		// Delete from database
		err = services.DeleteGuru(db, req.ID)
		if err != nil {
			log.Printf("ERROR: Gagal menghapus guru id=%d: %v", req.ID, err)
			utils.Error(w, http.StatusInternalServerError, "Gagal menghapus guru")
			return
		}

		log.Printf("Guru id=%d berhasil dihapus", req.ID)
		utils.Success(w, http.StatusOK, "Guru berhasil dihapus", nil)
	}
}
