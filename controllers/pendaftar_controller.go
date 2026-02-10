package controllers

import (
	"backend/dto"
	"backend/models"
	"backend/services"
	"backend/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

// SubmitPendaftaran handles POST /api/pendaftar with multipart form data
func SubmitPendaftaran(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
		}

		// Parse multipart form (max 32MB)
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Printf("ERROR: Gagal parse multipart form: %v", err)
			utils.Error(w, http.StatusBadRequest, "Gagal memproses form data")
			return
		}

		// Extract text fields
		pendaftar := models.Pendaftar{
			NamaLengkap:   r.FormValue("nama_lengkap"),
			NISN:          r.FormValue("nisn"),
			TempatLahir:   r.FormValue("tempat_lahir"),
			TanggalLahir:  r.FormValue("tanggal_lahir"), // Already in YYYY-MM-DD format from frontend
			Alamat:        r.FormValue("alamat"),
			AsalSekolah:   r.FormValue("asal_sekolah"),
			NamaAyah:      r.FormValue("nama_ayah"),
			PekerjaanAyah: r.FormValue("pekerjaan_ayah"),
			NamaIbu:       r.FormValue("nama_ibu"),
			PekerjaanIbu:  r.FormValue("pekerjaan_ibu"),
			NoHpOrtu:      r.FormValue("no_hp_ortu"),
		}

		log.Printf("Received pendaftaran: nama=%s, nisn=%s, tanggal_lahir=%s",
			pendaftar.NamaLengkap, pendaftar.NISN, pendaftar.TanggalLahir)

		// Validate required fields
		if pendaftar.NamaLengkap == "" || pendaftar.NISN == "" {
			utils.Error(w, http.StatusBadRequest, "Nama lengkap dan NISN wajib diisi")
			return
		}

		// Handle file uploads
		// Foto Profil
		if file, header, err := r.FormFile("foto_profil"); err == nil {
			filename, uploadErr := utils.UploadFoto(file, header, utils.FotoPendaftarPath)
			if uploadErr != nil {
				log.Printf("ERROR: Gagal upload foto profil: %v", uploadErr)
				utils.Error(w, http.StatusInternalServerError, "Gagal upload foto profil")
				return
			}
			pendaftar.FotoProfil = filename
			log.Printf("Foto profil uploaded: %s", filename)
		}

		// File SKL (Surat Keterangan Lulus)
		if file, header, err := r.FormFile("scan_skl"); err == nil {
			filename, uploadErr := utils.UploadFoto(file, header, utils.DokumenPendaftarPath)
			if uploadErr != nil {
				log.Printf("ERROR: Gagal upload file SKL: %v", uploadErr)
				utils.Error(w, http.StatusInternalServerError, "Gagal upload file SKL")
				return
			}
			pendaftar.FileSKL = filename
		}

		// File KK (Kartu Keluarga)
		if file, header, err := r.FormFile("scan_kk"); err == nil {
			filename, uploadErr := utils.UploadFoto(file, header, utils.DokumenPendaftarPath)
			if uploadErr != nil {
				log.Printf("ERROR: Gagal upload file KK: %v", uploadErr)
				utils.Error(w, http.StatusInternalServerError, "Gagal upload file KK")
				return
			}
			pendaftar.FileKK = filename
		}

		// File Akte (Akte Kelahiran)
		if file, header, err := r.FormFile("scan_akte"); err == nil {
			filename, uploadErr := utils.UploadFoto(file, header, utils.DokumenPendaftarPath)
			if uploadErr != nil {
				log.Printf("ERROR: Gagal upload file Akte: %v", uploadErr)
				utils.Error(w, http.StatusInternalServerError, "Gagal upload file Akte")
				return
			}
			pendaftar.FileAkte = filename
		}

		// File PIP/PKH (optional)
		if file, header, err := r.FormFile("scan_pkh"); err == nil {
			filename, uploadErr := utils.UploadFoto(file, header, utils.DokumenPendaftarPath)
			if uploadErr != nil {
				log.Printf("ERROR: Gagal upload file PIP: %v", uploadErr)
				utils.Error(w, http.StatusInternalServerError, "Gagal upload file PIP")
				return
			}
			pendaftar.FilePIP = filename
		}

		// Save to database
		id, registrasiID, err := services.CreatePendaftar(db, pendaftar)
		if err != nil {
			log.Printf("ERROR: Gagal menyimpan ke database: %v", err)
			utils.Error(w, http.StatusInternalServerError, "Gagal menyimpan data pendaftaran")
			return
		}

		log.Printf("Pendaftaran berhasil: id=%d, registrasi_id=%s", id, registrasiID)

		utils.Success(w, http.StatusCreated, "Pendaftaran berhasil dikirim", map[string]interface{}{
			"id":            id,
			"registrasi_id": registrasiID,
		})
	}
}

// GetAllPendaftar handles GET /api/admin/pendaftar
func GetAllPendaftar(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
		}

		pendaftarList, err := services.GetAllPendaftar(db)
		if err != nil {
			log.Printf("ERROR: Gagal mengambil data pendaftar: %v", err)
			utils.Error(w, http.StatusInternalServerError, "Gagal mengambil data pendaftar")
			return
		}

		utils.Success(w, http.StatusOK, "Data pendaftar berhasil diambil", pendaftarList)
	}
}

// UpdateStatusPendaftar handles PATCH /api/admin/pendaftar/status
func UpdateStatusPendaftar(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
		}

		var req dto.UpdateStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.Error(w, http.StatusBadRequest, "Format JSON tidak valid")
			return
		}

		if req.ID == 0 || req.Status == "" {
			utils.Error(w, http.StatusBadRequest, "ID dan status wajib diisi")
			return
		}

		// Validate status value
		validStatuses := map[string]bool{
			"Menunggu Verifikasi": true,
			"Terverifikasi":       true,
			"Ditolak":             true,
		}
		if !validStatuses[req.Status] {
			utils.Error(w, http.StatusBadRequest, "Status tidak valid")
			return
		}

		err := services.UpdatePendaftarStatus(db, req.ID, req.Status)
		if err != nil {
			log.Printf("ERROR: Gagal update status pendaftar id=%d: %v", req.ID, err)
			utils.Error(w, http.StatusInternalServerError, "Gagal memperbarui status")
			return
		}

		log.Printf("Status pendaftar id=%d diupdate menjadi: %s", req.ID, req.Status)
		utils.Success(w, http.StatusOK, "Status berhasil diperbarui", nil)
	}
}
