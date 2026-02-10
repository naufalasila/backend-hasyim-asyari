// dto/response.go
package dto

type SuccessResponse struct {
    Success bool        `json:"success"`
    Status  int         `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
    Success bool   `json:"success"`
    Status  int    `json:"status"`
    Message string `json:"message"`
}