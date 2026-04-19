package handler

import (
	"net/http"

	"gpt-load-mini/internal/data/model"
	"gpt-load-mini/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GroupHandler struct {
	db *gorm.DB
}

func NewGroupHandler(db *gorm.DB) *GroupHandler {
	return &GroupHandler{db: db}
}

func (h *GroupHandler) Create(c *gin.Context) {
	var group model.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}
	if err := h.db.Create(&group).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create group")
		return
	}
	utils.Success(c, group)
}

func (h *GroupHandler) List(c *gin.Context) {
	var groups []model.Group
	if err := h.db.Order("sort ASC").Find(&groups).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to list groups")
		return
	}
	utils.Success(c, groups)
}

func (h *GroupHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var group model.Group
	if err := h.db.First(&group, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Error(c, http.StatusNotFound, "Group not found")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "Failed to get group")
		return
	}
	utils.Success(c, group)
}

func (h *GroupHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var group model.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}
	if err := h.db.Where("id = ?", id).Select("name, display_name, channel_type, upstream_url, test_model, sort, proxy_api_key").Updates(&group).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update group")
		return
	}
	h.db.First(&group, id)
	utils.Success(c, group)
}

func (h *GroupHandler) GetByName(c *gin.Context) {
	name := c.Param("name")
	var group model.Group
	if err := h.db.Where("name = ?", name).First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Error(c, http.StatusNotFound, "Group not found")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "Failed to get group")
		return
	}
	utils.Success(c, group)
}

func (h *GroupHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&model.Group{}, id).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to delete group")
		return
	}
	utils.SuccessWithMessage(c, "Group deleted")
}

