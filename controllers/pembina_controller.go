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

func GetPembina(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
		}

		list, err := services.GetAllPembina(db)
		if err != nil {
			log.Printf("ERROR: Gagal mengambil data pembina: %v", err)
			utils.Error(w, http.StatusInternalServerError, "Gagal mengambil data pembina")
			return
		}

		utils.Success(w, http.StatusOK, "Data pembina berhasil diambil", list)
	}
}

func CreatePembina(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
		}

		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			log.Printf("ERROR: Gagal parse multipart form: %v", err)
			utils.Error(w, http.StatusBadRequest, "Gagal memproses form data")
			return
		}

		urutan := 0
		if urutanStr := r.FormValue("urutan"); urutanStr != "" {
			if val, err := strconv.Atoi(urutanStr); err == nil {
				urutan = val
			}
		}

		pembina := models.Pembina{
			Nama:       r.FormValue("nama"),
			Jabatan:    r.FormValue("jabatan"),
			Pendidikan: r.FormValue("pendidikan"),
			Urutan:     urutan,
		}

		if pembina.Nama == "" {
			utils.Error(w, http.StatusBadRequest, "Nama pembina wajib diisi")
			return
		}

		if file, header, err := r.FormFile("foto"); err == nil {
			filename, uploadErr := utils.UploadFoto(file, header, utils.PembinaImagePath)
			if uploadErr != nil {
				log.Printf("ERROR: Gagal upload foto pembina: %v", uploadErr)
				utils.Error(w, http.StatusInternalServerError, "Gagal upload foto")
				return
			}
			pembina.Foto = filename
		}

		id, err := services.CreatePembina(db, pembina)
		if err != nil {
			log.Printf("ERROR: Gagal menyimpan pembina ke database: %v", err)
			utils.Error(w, http.StatusInternalServerError, "Gagal menyimpan data pembina")
			return
		}

		log.Printf("Pembina berhasil ditambahkan: id=%d, nama=%s", id, pembina.Nama)
		utils.Success(w, http.StatusCreated, "Pembina berhasil ditambahkan", map[string]interface{}{
			"id": id,
		})
	}
}

func UpdatePembina(db *sql.DB) http.HandlerFunc {
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

		idStr := r.FormValue("id")
		if idStr == "" {
			utils.Error(w, http.StatusBadRequest, "ID pembina wajib disertakan")
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "ID tidak valid")
			return
		}

		existing, err := services.GetPembinaByID(db, id)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal mengambil data pembina")
			return
		}
		if existing == nil {
			utils.Error(w, http.StatusNotFound, "Pembina tidak ditemukan")
			return
		}

		urutan := existing.Urutan
		if urutanStr := r.FormValue("urutan"); urutanStr != "" {
			if val, err := strconv.Atoi(urutanStr); err == nil {
				urutan = val
			}
		}

		pembina := models.Pembina{
			ID:         id,
			Nama:       r.FormValue("nama"),
			Jabatan:    r.FormValue("jabatan"),
			Pendidikan: r.FormValue("pendidikan"),
			Urutan:     urutan,
			Foto:       existing.Foto,
		}

		if pembina.Nama == "" {
			utils.Error(w, http.StatusBadRequest, "Nama pembina wajib diisi")
			return
		}

		if file, header, err := r.FormFile("foto"); err == nil {
			filename, uploadErr := utils.UploadFoto(file, header, utils.PembinaImagePath)
			if uploadErr != nil {
				log.Printf("ERROR: Gagal upload foto pembina: %v", uploadErr)
				utils.Error(w, http.StatusInternalServerError, "Gagal upload foto")
				return
			}
			if existing.Foto != "" {
				utils.HapusFoto(utils.PembinaImagePath, existing.Foto)
			}
			pembina.Foto = filename
		}

		err = services.UpdatePembina(db, pembina)
		if err != nil {
			log.Printf("ERROR: Gagal mengupdate pembina id=%d: %v", id, err)
			utils.Error(w, http.StatusInternalServerError, "Gagal mengupdate data pembina")
			return
		}

		log.Printf("Pembina id=%d berhasil diupdate", id)
		utils.Success(w, http.StatusOK, "Pembina berhasil diupdate", nil)
	}
}

func DeletePembina(db *sql.DB) http.HandlerFunc {
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
				utils.Error(w, http.StatusBadRequest, "ID pembina wajib disertakan")
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
			utils.Error(w, http.StatusBadRequest, "ID pembina tidak valid")
			return
		}

		pembina, err := services.GetPembinaByID(db, req.ID)
		if err != nil {
			log.Printf("ERROR: Gagal mengambil data pembina: %v", err)
			utils.Error(w, http.StatusInternalServerError, "Gagal menghapus pembina")
			return
		}

		if pembina == nil {
			utils.Error(w, http.StatusNotFound, "Pembina tidak ditemukan")
			return
		}

		if pembina.Foto != "" {
			utils.HapusFoto(utils.PembinaImagePath, pembina.Foto)
		}

		err = services.DeletePembina(db, req.ID)
		if err != nil {
			log.Printf("ERROR: Gagal menghapus pembina id=%d: %v", req.ID, err)
			utils.Error(w, http.StatusInternalServerError, "Gagal menghapus pembina")
			return
		}

		log.Printf("Pembina id=%d berhasil dihapus", req.ID)
		utils.Success(w, http.StatusOK, "Pembina berhasil dihapus", nil)
	}
}
