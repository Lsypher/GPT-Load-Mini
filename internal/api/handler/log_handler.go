package handler

import (
	"net/http"
	"strconv"
	"time"

	"gpt-load-mini/internal/data/model"
	"gpt-load-mini/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LogHandler struct {
	db *gorm.DB
}

func NewLogHandler(db *gorm.DB) *LogHandler {
	return &LogHandler{db: db}
}

type PaginatedLogs struct {
	Data     []model.RequestLog `json:"data"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
	Total    int64              `json:"total"`
}

func (h *LogHandler) List(c *gin.Context) {
	var logs []model.RequestLog
	query := h.db.Model(&model.RequestLog{})

	page := 1
	pageSize := 100
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 1000 {
			pageSize = parsed
		}
	}
	offset := (page - 1) * pageSize

	if groupID := c.Query("group_id"); groupID != "" {
		query = query.Where("group_id = ?", groupID)
	}
	if model := c.Query("model"); model != "" {
		query = query.Where("model = ?", model)
	}
	if isSuccess := c.Query("is_success"); isSuccess != "" {
		query = query.Where("is_success = ?", isSuccess == "true")
	}
	if sourceIP := c.Query("source_ip"); sourceIP != "" {
		query = query.Where("source_ip = ?", sourceIP)
	}
	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("timestamp >= ?", t)
		}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			query = query.Where("timestamp <= ?", t.Add(24*time.Hour))
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to count logs")
		return
	}

	if err := query.Order("timestamp DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to list logs")
		return
	}
	utils.Success(c, PaginatedLogs{Data: logs, Page: page, PageSize: pageSize, Total: total})
}
