package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"videodb/be/services"
	"videodb/be/utils"

	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	videoService    *services.VideoService
	rtspService     *services.RTSPService
	workshopService *services.WorkshopService
}

func NewVideoHandler(vs *services.VideoService, rs *services.RTSPService, ws *services.WorkshopService) *VideoHandler {
	return &VideoHandler{
		videoService:    vs,
		rtspService:     rs,
		workshopService: ws,
	}
}

// @Summary 获取视频列表
// @Description 分页获取视频列表
// @Tags 视频管理
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param workshopId query int false "车间ID"
// @Param startTime query string false "开始时间"
// @Param endTime query string false "结束时间"
// @Success 200 {object} utils.Response
// @Router /api/videos [get]
func (h *VideoHandler) List(c *gin.Context) {
	var query struct {
		Page       int    `form:"page" binding:"required,min=1"`
		PageSize   int    `form:"pageSize" binding:"required,min=1,max=100"`
		WorkshopID uint   `form:"workshopId"`
		StartTime  string `form:"startTime"`
		EndTime    string `form:"endTime"`
	}

	if err := c.ShouldBindQuery(&query); err != nil {
		utils.Error(c, err)
		return
	}

	// 解析时间
	var startTime, endTime time.Time
	var err error
	if query.StartTime != "" {
		startTime, err = utils.ParseTime(query.StartTime)
		if err != nil {
			utils.Error(c, fmt.Errorf("invalid start time format"))
			return
		}
	}
	if query.EndTime != "" {
		endTime, err = utils.ParseTime(query.EndTime)
		if err != nil {
			utils.Error(c, fmt.Errorf("invalid end time format"))
			return
		}
	}

	videos, total, err := h.videoService.List(services.VideoQuery{
		WorkshopID: query.WorkshopID,
		StartTime:  startTime,
		EndTime:    endTime,
		Page:       query.Page,
		PageSize:   query.PageSize,
		Preload:    []string{"Workshop"},
	})

	if err != nil {
		utils.Error(c, err)
		return
	}

	utils.Success(c, gin.H{
		"list":  videos,
		"total": total,
		"page":  query.Page,
		"size":  query.PageSize,
	})
}

// @Summary 获取视频详情
// @Description 根据ID获取视频详情
// @Tags 视频管理
// @Accept json
// @Produce json
// @Param id path int true "视频ID"
// @Success 200 {object} utils.Response
// @Router /api/videos/{id} [get]
func (h *VideoHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, fmt.Errorf("invalid id format"))
		return
	}

	video, err := h.videoService.GetByID(uint(id))
	if err != nil {
		utils.Error(c, err)
		return
	}

	utils.Success(c, video)
}

// @Summary 删除视频
// @Description 根据ID删除视频
// @Tags 视频管理
// @Accept json
// @Produce json
// @Param id path int true "视频ID"
// @Success 200 {object} utils.Response
// @Router /api/videos/{id} [delete]
func (h *VideoHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, fmt.Errorf("invalid id format"))
		return
	}

	if err := h.videoService.Delete(uint(id)); err != nil {
		utils.Error(c, err)
		return
	}

	utils.Success(c, nil)
}

// @Summary 开始录制视频
// @Description 开始录制指定车间的视频
// @Tags 视频管理
// @Accept json
// @Produce json
// @Param body body StartRecordingRequest true "录制参数"
// @Success 200 {object} utils.Response
// @Router /api/videos/start-recording [post]
func (h *VideoHandler) StartRecording(c *gin.Context) {
	var req struct {
		WorkshopID uint   `json:"workshopId" binding:"required"`
		Duration   string `json:"duration" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err)
		return
	}
	// 获取车间信息
	workshop, err := h.workshopService.GetByID(req.WorkshopID)
	if err != nil {
		utils.Error(c, err)
		return
	}

	// 生成输出文件路径
	outputPath := h.videoService.GenerateVideoPath(req.WorkshopID)

	// 开始录制
	if err := h.rtspService.StartRecording(c, workshop.RTSPUrl, outputPath, req.WorkshopID); err != nil {
		utils.Error(c, err)
		return
	}

	utils.Success(c, gin.H{"message": "Recording started"})
}

// @Summary 停止录制视频
// @Description 停止录制指定车间的视频
// @Tags 视频管理
// @Accept json
// @Produce json
// @Param workshopId path int true "车间ID"
// @Success 200 {object} utils.Response
// @Router /api/videos/stop-recording/{workshopId} [post]
func (h *VideoHandler) StopRecording(c *gin.Context) {
	workshopID, err := strconv.ParseUint(c.Param("workshopId"), 10, 32)
	if err != nil {
		utils.Error(c, fmt.Errorf("invalid workshop id format"))
		return
	}

	if err := h.rtspService.StopRecording(uint(workshopID)); err != nil {
		utils.Error(c, err)
		return
	}

	utils.Success(c, gin.H{"message": "Recording stopped"})
}

// @Summary 下载视频
// @Description 下载指定ID的视频文件
// @Tags 视频管理
// @Accept json
// @Produce json
// @Param id path int true "视频ID"
// @Success 200 {file} binary
// @Router /api/videos/{id}/download [get]
func (h *VideoHandler) Download(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, fmt.Errorf("invalid id format"))
		return
	}

	video, err := h.videoService.GetByID(uint(id))
	if err != nil {
		utils.Error(c, err)
		return
	}

	if !utils.FileExists(video.FilePath) {
		utils.Error(c, fmt.Errorf("video file not found"))
		return
	}

	c.FileAttachment(video.FilePath, video.FileName)
}

// 添加批量删除的请求结构
type BatchDeleteRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}

// @Summary 批量删除视频
// @Description 批量删除多个视频
// @Tags 视频管理
// @Accept json
// @Produce json
// @Param body body BatchDeleteRequest true "视频ID列表"
// @Success 200 {object} utils.Response
// @Router /api/videos/batch [delete]
func (h *VideoHandler) BatchDelete(c *gin.Context) {
	var req BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err)
		return
	}

	if err := h.videoService.BatchDelete(req.IDs); err != nil {
		utils.Error(c, err)
		return
	}

	utils.Success(c, nil)
}

func (h *VideoHandler) StreamVideo(c *gin.Context) {
	filePath := c.Query("path")

	// 安全检查：确保路径在允许的目录内
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path"})
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// 设置响应头
	c.Header("Content-Type", "video/mp4")
	c.Header("Content-Disposition", "inline")

	// 流式传输视频
	c.File(absPath)
}
