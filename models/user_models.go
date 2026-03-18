// models/user.go
package models

import "time"

type User struct {
	IDUser         int       `json:"id_user"`
	Username       string    `json:"username"`
	FullName       string    `json:"full_name"`
	Email          string    `json:"email"`
	Password       string    `json:"-"`
	ProfilePicture string    `json:"profile_picture"`
	Role           string    `json:"role"`
	IsVerified     bool      `json:"is_verified"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Pendaftar represents the pendaftar table for student registration
type Pendaftar struct {
	ID            int       `json:"id"`
	RegistrasiID  string    `json:"registrasi_id"`
	NamaLengkap   string    `json:"nama_lengkap"`
	NISN          string    `json:"nisn"`
	TempatLahir   string    `json:"tempat_lahir"`
	TanggalLahir  string    `json:"tanggal_lahir"`
	Alamat        string    `json:"alamat"`
	AsalSekolah   string    `json:"asal_sekolah"`
	NamaAyah      string    `json:"nama_ayah"`
	PekerjaanAyah string    `json:"pekerjaan_ayah"`
	NamaIbu       string    `json:"nama_ibu"`
	PekerjaanIbu  string    `json:"pekerjaan_ibu"`
	NoHpOrtu      string    `json:"no_hp_ortu"`
	FotoProfil    string    `json:"foto_profil"`
	FileSKL       string    `json:"file_skl"`
	FileKK        string    `json:"file_kk"`
	FileAkte      string    `json:"file_akte"`
	FilePIP       string    `json:"file_pip"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

// Album represents the album table for gallery
type Album struct {
	ID        int       `json:"id"`
	Judul     string    `json:"judul"`
	Gambar    string    `json:"gambar"`
	Kategori  string    `json:"kategori"`
	CreatedAt time.Time `json:"created_at"`
}
