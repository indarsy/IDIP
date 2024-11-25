package handlers

import (
	"net/http"
	"strconv"
	"videodb/be/models"
	"videodb/be/services"

	"github.com/gin-gonic/gin"
)

type CaptureHandler struct {
	captureService *services.CaptureService
}

func NewCaptureHandler(captureService *services.CaptureService) *CaptureHandler {
	return &CaptureHandler{
		captureService: captureService,
	}
}

// Create 创建采集任务
func (h *CaptureHandler) Create(c *gin.Context) {
	var capture models.Capture
	if err := c.ShouldBindJSON(&capture); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的请求数据",
			"error":   err.Error(),
		})
		return
	}

	if err := h.captureService.Create(&capture); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建采集任务失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "创建成功",
		"data":    capture,
	})
}

// List 获取采集任务列表
func (h *CaptureHandler) List(c *gin.Context) {
	var captures []models.Capture
	var err error

	// 获取可选的 workshopId 参数
	workshopIDStr := c.Query("workshopId")
	if workshopIDStr != "" {
		// 如果提供了 workshopId，则按车间筛选
		workshopID, err := strconv.ParseUint(workshopIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "无效的车间ID",
				"error":   err.Error(),
			})
			return
		}
		captures, err = h.captureService.List(uint(workshopID))
	} else {
		// 如果没有提供 workshopId，则获取所有任务
		captures, err = h.captureService.ListAll()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取采集任务列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": captures,
	})
}

// Get 获取单个采集任务
func (h *CaptureHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的任务ID",
			"error":   err.Error(),
		})
		return
	}

	capture, err := h.captureService.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取采集任务失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": capture,
	})
}

// Cancel 取消采集任务
func (h *CaptureHandler) Cancel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的任务ID",
			"error":   err.Error(),
		})
		return
	}

	if err := h.captureService.Cancel(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "取消采集任务失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "取消成功",
	})
}
