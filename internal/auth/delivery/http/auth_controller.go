package http

import (
	"net/http"

	"backend/internal/auth/delivery/http/dto"
	"backend/internal/auth/usecase"
	"backend/internal/domain"
	"backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	signupUsecase usecase.SignupUsecase
	loginUsecase  usecase.LoginUsecase
}

func NewAuthController(
	signupUsecase usecase.SignupUsecase,
	loginUsecase usecase.LoginUsecase,
) *AuthController {
	return &AuthController{
		signupUsecase: signupUsecase,
		loginUsecase:  loginUsecase,
	}
}

func (h *AuthController) Signup(c *gin.Context) {
	var req dto.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	role := domain.RolePassenger
	if req.Role == "driver" {
		role = domain.RoleDriver
	}

	token, err := h.signupUsecase.Execute(c.Request.Context(), usecase.SignupInput{
		Email:         req.Email,
		Password:      req.Password,
		FullName:      req.FullName,
		Phone:         req.Phone,
		Role:          role,
		PlateNumber:   req.PlateNumber,
		VehicleType:   req.VehicleType,
		VehicleColor:  req.VehicleColor,
		LicenseNumber: req.LicenseNumber,
	})
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Signup failed: "+err.Error())
		return
	}

	response.Created(c, dto.AuthResponse{Token: token})
}

func (h *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	token, err := h.loginUsecase.Execute(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.Unauthorized(c, "Invalid credentials")
		return
	}

	response.Success(c, dto.AuthResponse{Token: token})
}
