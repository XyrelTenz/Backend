package http

import (
	"net/http"

	"backend/internal/domain"
	"backend/internal/passenger/delivery/http/dto"
	"backend/internal/passenger/usecase"
	"backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type PassengerController struct {
	requestRideUsecase    usecase.RequestRideUsecase
	getRideUsecase        usecase.GetRideUsecase
	getHistoryUsecase     usecase.GetPassengerHistoryUsecase
	addSavedPlaceUsecase  usecase.AddSavedPlaceUsecase
	getSavedPlacesUsecase usecase.GetSavedPlacesUsecase
}

func NewPassengerController(
	requestRideUsecase usecase.RequestRideUsecase,
	getRideUsecase usecase.GetRideUsecase,
	getHistoryUsecase usecase.GetPassengerHistoryUsecase,
	addSavedPlaceUsecase usecase.AddSavedPlaceUsecase,
	getSavedPlacesUsecase usecase.GetSavedPlacesUsecase,
) *PassengerController {
	return &PassengerController{
		requestRideUsecase:    requestRideUsecase,
		getRideUsecase:        getRideUsecase,
		getHistoryUsecase:     getHistoryUsecase,
		addSavedPlaceUsecase:  addSavedPlaceUsecase,
		getSavedPlacesUsecase: getSavedPlacesUsecase,
	}
}

func (h *PassengerController) Request(c *gin.Context) {
	var req dto.RequestRideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	userID, _ := c.Get("user_id")

	ride, err := h.requestRideUsecase.Execute(c.Request.Context(), usecase.RequestRideInput{
		PassengerID:    userID.(string),
		PickupAddress:  req.PickupAddress,
		PickupLat:      req.PickupLat,
		PickupLng:      req.PickupLng,
		DropoffAddress: req.DropoffAddress,
		DropoffLat:     req.DropoffLat,
		DropoffLng:     req.DropoffLng,
		VehicleType:    req.VehicleType,
		PaymentMethod:  req.PaymentMethod,
	})

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Created(c, mapToRideResponse(ride))
}

func (h *PassengerController) GetRide(c *gin.Context) {
	id := c.Param("id")
	ride, err := h.getRideUsecase.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Ride not found")
		return
	}

	response.Success(c, mapToRideResponse(ride))
}

func (h *PassengerController) GetHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	rides, err := h.getHistoryUsecase.Execute(c.Request.Context(), userID.(string))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	res := make([]dto.RideResponse, len(rides))
	for i, r := range rides {
		res[i] = mapToRideResponse(r)
	}

	response.Success(c, res)
}

func (h *PassengerController) AddSavedPlace(c *gin.Context) {
	var req dto.SavedPlaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	err := h.addSavedPlaceUsecase.Execute(c.Request.Context(), &domain.SavedPlace{
		UserID:  userID.(string),
		Name:    req.Name,
		Address: req.Address,
		Lat:     req.Lat,
		Lng:     req.Lng,
		Type:    req.Type,
	})

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, "Place saved successfully")
}

func (h *PassengerController) GetSavedPlaces(c *gin.Context) {
	userID, _ := c.Get("user_id")
	places, err := h.getSavedPlacesUsecase.Execute(c.Request.Context(), userID.(string))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	res := make([]dto.SavedPlaceResponse, len(places))
	for i, p := range places {
		res[i] = dto.SavedPlaceResponse{
			ID:      p.ID,
			Name:    p.Name,
			Address: p.Address,
			Lat:     p.Lat,
			Lng:     p.Lng,
			Type:    p.Type,
		}
	}

	response.Success(c, res)
}

func mapToRideResponse(ride *domain.Ride) dto.RideResponse {
	return dto.RideResponse{
		ID:                    ride.ID,
		PickupAddress:         ride.PickupAddress,
		DropoffAddress:        ride.DropoffAddress,
		DistanceKM:            ride.DistanceKM,
		EstimatedFareAmount:   ride.EstimatedFareAmount,
		EstimatedDurationMins: ride.EstimatedDurationMins,
		Status:                string(ride.Status),
		VehicleType:           ride.VehicleType,
	}
}
