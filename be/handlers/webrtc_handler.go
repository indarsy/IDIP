package handlers

import (
	"net/http"

	"videodb/be/models"
	"videodb/be/services"

	"github.com/gin-gonic/gin"
)

type WebRTCHandler struct {
	webrtcService *services.WebRTCService
}

func NewWebRTCHandler(webrtcService *services.WebRTCService) *WebRTCHandler {
	return &WebRTCHandler{
		webrtcService: webrtcService,
	}
}

func (h *WebRTCHandler) HandleWebRTC(c *gin.Context) {
	var req models.WebRTCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.WebRTCResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	answer, err := h.webrtcService.HandleRTSP(req.RTSPURL, req.SDP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.WebRTCResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, // 前端期望的成功状态码
		"data": models.WebRTCResponse{
			Success: true,
			SDP:     answer.SDP,
		},
		"message": "success",
	})
}
