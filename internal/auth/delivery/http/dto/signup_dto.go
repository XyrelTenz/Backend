package dto

type SignupRequest struct {
	Email         string `json:"email" binding:"required,email"`
	Password      string `json:"password" binding:"required,min=6"`
	FullName      string `json:"full_name" binding:"required"`
	Phone         string `json:"phone" binding:"required"`
	Role          string `json:"role" binding:"required"` // "passenger" or "driver"
	PlateNumber   string `json:"plate_number"`
	VehicleType   string `json:"vehicle_type"`
	VehicleColor  string `json:"vehicle_color"`
	LicenseNumber string `json:"license_number"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
