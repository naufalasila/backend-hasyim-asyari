package models

// Berita represents news/article in the system
type Berita struct {
	ID        int    `json:"id"`
	Tema      string `json:"tema"`
	Judul     string `json:"judul"`
	Gambar    string `json:"gambar"`
	Isi       string `json:"isi"`
	Tanggal   string `json:"tanggal"`
	CreatedAt string `json:"created_at,omitempty"`
}

/*
SQL Schema:

CREATE TABLE berita (
    id INT AUTO_INCREMENT PRIMARY KEY,
    tema VARCHAR(50) NOT NULL,
    judul VARCHAR(255) NOT NULL,
    gambar VARCHAR(255),
    isi TEXT NOT NULL,
    tanggal DATE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
*/
