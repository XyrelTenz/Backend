package http

import (
	"fmt"
	"net/http"

	"backend/internal/driver/delivery/http/dto"
	"backend/internal/driver/usecase"
	"backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type DriverController struct {
	acceptRideUsecase       usecase.AcceptRideUsecase
	updateRideStatusUsecase usecase.UpdateRideStatusUsecase
	updateLocationUsecase   usecase.UpdateLocationUsecase
	getNearbyRidesUsecase   usecase.GetNearbyRidesUsecase
}

func NewDriverController(
	acceptRideUsecase usecase.AcceptRideUsecase,
	updateRideStatusUsecase usecase.UpdateRideStatusUsecase,
	updateLocationUsecase usecase.UpdateLocationUsecase,
	getNearbyRidesUsecase usecase.GetNearbyRidesUsecase,
) *DriverController {
	return &DriverController{
		acceptRideUsecase:       acceptRideUsecase,
		updateRideStatusUsecase: updateRideStatusUsecase,
		updateLocationUsecase:   updateLocationUsecase,
		getNearbyRidesUsecase:   getNearbyRidesUsecase,
	}
}

func (h *DriverController) Accept(c *gin.Context) {
	rideID := c.Param("id")
	userID := c.GetString("user_id")

	if err := h.acceptRideUsecase.Execute(c.Request.Context(), rideID, userID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Ride accepted successfully"})
}

func (h *DriverController) UpdateStatus(c *gin.Context) {
	rideID := c.Param("id")
	var req dto.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if err := h.updateRideStatusUsecase.Execute(
		c.Request.Context(),
		rideID,
		req.Status,
	); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update status: "+err.Error())
		return
	}

	response.Success(c, gin.H{"message": fmt.Sprintf("Ride status updated to %s", req.Status)})
}

func (h *DriverController) GetNearby(c *gin.Context) {
	lat := 7.8285
	lng := 123.4344

	rides, err := h.getNearbyRidesUsecase.Execute(c.Request.Context(), lat, lng)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	res := make([]dto.RideDTO, len(rides))
	for i, r := range rides {
		res[i] = dto.RideDTO{
			ID:             r.ID,
			PickupAddress:  r.PickupAddress,
			DropoffAddress: r.DropoffAddress,
			DistanceKM:     r.DistanceKM,
			Fare:           r.EstimatedFareAmount,
			Status:         string(r.Status),
		}
	}

	response.Success(c, dto.NearbyRidesResponse{Rides: res})
}

func (h *DriverController) UpdateLocation(c *gin.Context) {
	var req dto.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	userID := c.GetString("user_id")
	if err := h.updateLocationUsecase.Execute(
		c.Request.Context(),
		userID,
		req.Lat,
		req.Lng,
	); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Location updated"})
}
