package models

type Brosur struct {
	ID        int    `json:"id"`
	Filename  string `json:"filename"`
	CreatedAt string `json:"created_at,omitempty"`
}
