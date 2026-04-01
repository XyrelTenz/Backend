package usecase

import (
	"time"

	chat_ws "backend/internal/chat/ws"
)

type BroadcastLocationUsecase interface {
	Execute(rideID, driverID string, lat, lng float64)
}

type broadcastLocationUsecase struct {
	hub *chat_ws.Hub
}

func NewBroadcastLocationUsecase(hub *chat_ws.Hub) BroadcastLocationUsecase {
	return &broadcastLocationUsecase{
		hub: hub,
	}
}

func (u *broadcastLocationUsecase) Execute(rideID, driverID string, lat, lng float64) {
	u.hub.Broadcast(rideID, map[string]interface{}{
		"type": "LOCATION_UPDATE",
		"payload": map[string]interface{}{
			"driver_id": driverID,
			"lat":       lat,
			"lng":       lng,
		},
		"timestamp": time.Now(),
	})
}
