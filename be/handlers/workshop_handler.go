package handlers

import (
	"fmt"
	"strconv"

	"videodb/be/models"
	"videodb/be/services"
	"videodb/be/utils"

	"github.com/gin-gonic/gin"
)

type WorkshopHandler struct {
	workshopService *services.WorkshopService
	rtspService     *services.RTSPService
}

func NewWorkshopHandler(ws *services.WorkshopService, rs *services.RTSPService) *WorkshopHandler {
	return &WorkshopHandler{
		workshopService: ws,
		rtspService:     rs,
	}
}

// @Summary 获取车间列表
// @Description 获取所有车间信息
// @Tags 车间管理
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/workshops [get]
func (h *WorkshopHandler) List(c *gin.Context) {
	workshops, err := h.workshopService.List()
	if err != nil {
		utils.Error(c, err)
		return
	}

	utils.Success(c, workshops)
}

// @Summary 创建车间
// @Description 创建新的车间信息
// @Tags 车间管理
// @Accept json
// @Produce json
// @Param body body models.Workshop true "车间信息"
// @Success 200 {object} utils.Response
// @Router /api/workshops [post]
func (h *WorkshopHandler) Create(c *gin.Context) {
	var workshop models.Workshop
	if err := c.ShouldBindJSON(&workshop); err != nil {
		utils.Error(c, err)
		return
	}

	// 验证RTSP地址
	// if err := h.rtspService.CheckRTSPStream(workshop.RTSPUrl); err != nil {
	// 	utils.Error(c, fmt.Errorf("invalid RTSP URL: %v", err))
	// 	return
	// }

	if err := h.workshopService.Create(&workshop); err != nil {
		utils.Error(c, err)
		return
	}

	utils.Success(c, workshop)
}

// @Summary 更新车间信息
// @Description 更新指定车间的信息
// @Tags 车间管理
// @Accept json
// @Produce json
// @Param id path int true "车间ID"
// @Param body body models.Workshop true "车间信息"
// @Success 200 {object} utils.Response
// @Router /api/workshops/{id} [put]
func (h *WorkshopHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, fmt.Errorf("invalid id format"))
		return
	}

	// 检查id是否存在
	if _, err := h.workshopService.GetByID(uint(id)); err != nil {
		utils.Error(c, err)
		return
	}

	var workshop models.Workshop
	if err := c.ShouldBindJSON(&workshop); err != nil {
		fmt.Println("err:", err)
		utils.Error(c, err)
		return
	}

	workshop.ID = uint(id)
	if err := h.workshopService.Update(uint(id), &workshop); err != nil {
		utils.Error(c, err)
		return
	}

	utils.Success(c, workshop)
}

// @Summary 删除车间
// @Description 删除指定的车间
// @Tags 车间管理
// @Accept json
// @Produce json
// @Param id path int true "车间ID"
// @Success 200 {object} utils.Response
// @Router /api/workshops/{id} [delete]
func (h *WorkshopHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, fmt.Errorf("invalid id format"))
		return
	}

	if err := h.workshopService.Delete(uint(id)); err != nil {
		utils.Error(c, err)
		return
	}

	utils.Success(c, nil)
}

// @Summary 获取车间预览流
// @Description 获取指定车间的视频预览流地址
// @Tags 车间管理
// @Accept json
// @Produce json
// @Param id path int true "车间ID"
// @Success 200 {object} utils.Response
// @Router /api/workshops/{id}/preview [get]
func (h *WorkshopHandler) GetPreview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, fmt.Errorf("invalid id format"))
		return
	}

	workshop, err := h.workshopService.GetByID(uint(id))
	if err != nil {
		utils.Error(c, err)
		return
	}

	previewUrl, err := h.rtspService.GetPreviewURL(workshop.ID)
	if err != nil {
		utils.Error(c, err)
		return
	}

	utils.Success(c, gin.H{"previewUrl": previewUrl})
}
