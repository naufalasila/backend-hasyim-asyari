package models

// Guru represents a teacher/staff member
type Guru struct {
	ID            int    `json:"id"`
	Nama          string `json:"nama"`
	Jabatan       string `json:"jabatan"`
	MataPelajaran string `json:"mata_pelajaran"`
	Pendidikan    string `json:"pendidikan"`
	Foto          string `json:"foto"`
	Jenjang       string `json:"jenjang"`
	Urutan        int    `json:"urutan"`
	CreatedAt     string `json:"created_at,omitempty"`
}

/*
SQL Schema:

CREATE TABLE guru (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nama VARCHAR(150) NOT NULL,
    jabatan VARCHAR(100) DEFAULT '',
    mata_pelajaran VARCHAR(100) DEFAULT '',
    pendidikan VARCHAR(150) DEFAULT '',
    foto VARCHAR(255) DEFAULT '',
    jenjang ENUM('mts', 'ma') NOT NULL DEFAULT 'mts',
    urutan INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
*/
