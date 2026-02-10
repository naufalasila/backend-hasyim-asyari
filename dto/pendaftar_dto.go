package dto

// UpdateStatusRequest is used for PATCH /api/admin/pendaftar/status
type UpdateStatusRequest struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}
