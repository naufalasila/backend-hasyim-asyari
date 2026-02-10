// dto/profile.go
package dto

type UpdateProfileRequest struct {
    FullName string `json:"full_name" validate:"required,min=2,max=100"`
}

type UserProfileResponse struct {
    FullName            string `json:"full_name"`
    Email               string `json:"email"`
    ProfilePicture      string `json:"profile_picture"`
    TanggalBergabung    string `json:"tanggal_bergabung"`
    StatusKeanggotaan   string `json:"status_keanggotaan"`
}