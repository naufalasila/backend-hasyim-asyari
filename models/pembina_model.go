package models

type Pembina struct {
	ID         int    `json:"id"`
	Nama       string `json:"nama"`
	Jabatan    string `json:"jabatan"`
	Pendidikan string `json:"pendidikan"`
	Foto       string `json:"foto"`
	Urutan     int    `json:"urutan"`
	CreatedAt  string `json:"created_at,omitempty"`
}
