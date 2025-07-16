package http

import (
	"duking/internal/models"
	"duking/internal/usecase"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RoomHandler struct {
	roomService usecase.RoomService
	logger      *zap.Logger
}

func NewRoomHandler(roomService usecase.RoomService, logger *zap.Logger) *RoomHandler {
	return &RoomHandler{
		roomService: roomService,
		logger:      logger,
	}
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var input models.RoomCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("failed to bind JSON for room creation", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	room, err := h.roomService.CreateRoom(c.Request.Context(), input)
	if err != nil {
		h.logger.Error("failed to create room", zap.Error(err), zap.Any("input", input))
		if strings.Contains(err.Error(), "room number already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room"})
		return
	}
	c.JSON(http.StatusCreated, room)
}

func (h *RoomHandler) GetRoom(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Error("invalid room ID parameter", zap.String("id", idStr), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room id"})
		return
	}
	room, err := h.roomService.GetRoom(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("failed to get room", zap.Error(err), zap.Int64("room_id", id))
		if err.Error() == "room not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve room"})
		return
	}
	c.JSON(http.StatusOK, room)
}

func (h *RoomHandler) GetRoomsByHotel(c *gin.Context) {
	hotelIDStr := c.Param("hotel_id")
	hotelID, err := strconv.ParseInt(hotelIDStr, 10, 64)
	if err != nil {
		h.logger.Error("invalid hotel ID parameter for rooms", zap.String("hotel_id", hotelIDStr), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID"})
		return
	}

	rooms, err := h.roomService.GetRoomsByHotel(c.Request.Context(), hotelID)
	if err != nil {
		h.logger.Error("failed to get rooms by hotel ID", zap.Error(err), zap.Int64("hotel_id", hotelID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve rooms"})
		return
	}
	c.JSON(http.StatusOK, rooms)
}

func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Error("invalid room ID parameter for update", zap.String("id", idStr), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	var input models.RoomUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("failed to bind JSON for room update", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	updatedRoom, err := h.roomService.UpdateRoom(c.Request.Context(), id, input)
	if err != nil {
		h.logger.Error("failed to update room", zap.Error(err), zap.Int64("room_id", id), zap.Any("input", input))
		if err.Error() == "room not found for update" || err.Error() == "room not found or was deleted during update" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room"})
		return
	}
	c.JSON(http.StatusOK, updatedRoom)
}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Error("invalid room ID parameter for delete", zap.String("id", idStr), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	if err := h.roomService.DeleteRoom(c.Request.Context(), id); err != nil {
		h.logger.Error("failed to delete room", zap.Error(err), zap.Int64("room_id", id))
		if err.Error() == "room not found for deletion" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete room"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Room deleted successfully"})
}
