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

// GetAlbums handles GET /api/album - public endpoint, filter by ?kategori=MTS|MA
func GetAlbums(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
		}

		kategori := r.URL.Query().Get("kategori")
		if kategori == "" {
			kategori = "MTS"
		}

		albumList, err := services.GetAlbumsByKategori(db, kategori)
		if err != nil {
			log.Printf("ERROR: Gagal mengambil data album: %v", err)
			utils.Error(w, http.StatusInternalServerError, "Gagal mengambil data album")
			return
		}

		utils.Success(w, http.StatusOK, "Data album berhasil diambil", albumList)
	}
}

// CreateAlbum handles POST /api/admin/album with multipart form data
func CreateAlbum(db *sql.DB) http.HandlerFunc {
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

		album := models.Album{
			Judul:    r.FormValue("judul"),
			Kategori: r.FormValue("kategori"),
		}

		if album.Judul == "" {
			utils.Error(w, http.StatusBadRequest, "Judul album wajib diisi")
			return
		}
		if album.Kategori == "" {
			utils.Error(w, http.StatusBadRequest, "Kategori album wajib diisi")
			return
		}

		// Handle image upload
		if file, header, err := r.FormFile("gambar"); err == nil {
			filename, uploadErr := utils.UploadFoto(file, header, utils.AlbumImagePath)
			if uploadErr != nil {
				log.Printf("ERROR: Gagal upload gambar album: %v", uploadErr)
				utils.Error(w, http.StatusInternalServerError, "Gagal upload gambar")
				return
			}
			album.Gambar = filename
			log.Printf("Gambar album uploaded: %s", filename)
		} else {
			utils.Error(w, http.StatusBadRequest, "Gambar album wajib diupload")
			return
		}

		// Save to database
		id, err := services.CreateAlbum(db, album)
		if err != nil {
			log.Printf("ERROR: Gagal menyimpan album ke database: %v", err)
			utils.Error(w, http.StatusInternalServerError, "Gagal menyimpan data album")
			return
		}

		log.Printf("Album berhasil ditambahkan: id=%d, judul=%s", id, album.Judul)
		utils.Success(w, http.StatusCreated, "Album berhasil ditambahkan", map[string]interface{}{
			"id":     id,
			"gambar": album.Gambar,
		})
	}
}

// DeleteAlbum handles DELETE /api/admin/album
func DeleteAlbum(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
		}

		var req struct {
			ID int `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			idStr := r.URL.Query().Get("id")
			if idStr == "" {
				utils.Error(w, http.StatusBadRequest, "ID album wajib disertakan")
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
			utils.Error(w, http.StatusBadRequest, "ID album tidak valid")
			return
		}

		// Get album to delete the image file
		album, err := services.GetAlbumByID(db, req.ID)
		if err != nil {
			log.Printf("ERROR: Gagal mengambil data album: %v", err)
			utils.Error(w, http.StatusInternalServerError, "Gagal menghapus album")
			return
		}

		if album == nil {
			utils.Error(w, http.StatusNotFound, "Album tidak ditemukan")
			return
		}

		// Delete the image file if exists
		if album.Gambar != "" {
			utils.HapusFoto(utils.AlbumImagePath, album.Gambar)
		}

		// Delete from database
		err = services.DeleteAlbum(db, req.ID)
		if err != nil {
			log.Printf("ERROR: Gagal menghapus album id=%d: %v", req.ID, err)
			utils.Error(w, http.StatusInternalServerError, "Gagal menghapus album")
			return
		}

		log.Printf("Album id=%d berhasil dihapus", req.ID)
		utils.Success(w, http.StatusOK, "Album berhasil dihapus", nil)
	}
}
