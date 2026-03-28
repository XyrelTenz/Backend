package http

import (
	"log"
	"net/http"

	"backend/internal/chat/delivery/http/dto"
	"backend/internal/chat/usecase"
	"backend/internal/chat/ws"
	"backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ChatController struct {
	sendMessageUsecase    usecase.SendMessageUsecase
	getChatHistoryUsecase usecase.GetChatHistoryUsecase
	hub                   *chat_ws.Hub
}

func NewChatController(
	sendMessageUsecase usecase.SendMessageUsecase,
	getChatHistoryUsecase usecase.GetChatHistoryUsecase,
	hub *chat_ws.Hub,
) *ChatController {
	return &ChatController{
		sendMessageUsecase:    sendMessageUsecase,
		getChatHistoryUsecase: getChatHistoryUsecase,
		hub:                   hub,
	}
}

func (h *ChatController) HandleWebSocket(c *gin.Context) {
	rideID := c.Param("id")
	userID := c.Query("user_id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &chat_ws.Client{
		RideID: rideID,
		Conn:   conn,
		UserID: userID,
	}

	h.hub.Register(client)
	defer h.hub.Unregister(client)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (h *ChatController) SendMessage(c *gin.Context) {
	rideID := c.Param("id")
	var req dto.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	userID := c.GetString("user_id")
	msg, err := h.sendMessageUsecase.Execute(c.Request.Context(), rideID, userID, req.Content)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Created(c, dto.MessageResponse{
		ID:        msg.ID,
		RideID:    msg.RideID,
		SenderID:  msg.SenderID,
		Content:   msg.Message,
		CreatedAt: msg.CreatedAt,
	})
}

func (h *ChatController) GetHistory(c *gin.Context) {
	rideID := c.Param("id")
	messages, err := h.getChatHistoryUsecase.Execute(c.Request.Context(), rideID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	res := make([]dto.MessageResponse, len(messages))
	for i, msg := range messages {
		res[i] = dto.MessageResponse{
			ID:        msg.ID,
			RideID:    msg.RideID,
			SenderID:  msg.SenderID,
			Content:   msg.Message,
			CreatedAt: msg.CreatedAt,
		}
	}

	response.Success(c, res)
}
