// utils/upload.go
package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

const FotoPendaftarPath = "uploads/foto_pendaftar"
const ProfilePhotoPath = "uploads/profile"
const DokumenPendaftarPath = "uploads/dokumen_pendaftar"
const BeritaImagePath = "uploads/berita"

func UploadFoto(file multipart.File, header *multipart.FileHeader, uploadDir string) (string, error) {
	defer file.Close()

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}

	ext := filepath.Ext(header.Filename)
	newName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	path := filepath.Join(uploadDir, newName)

	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	return newName, nil
}

func HapusFoto(uploadDir, filename string) error {
	if filename == "" {
		return nil
	}
	path := filepath.Join(uploadDir, filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(path)
}
