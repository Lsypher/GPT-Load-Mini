package handler

import (
	"gpt-load-mini/internal/data/model"
	"gpt-load-mini/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StatsHandler struct {
	db *gorm.DB
}

func NewStatsHandler(db *gorm.DB) *StatsHandler {
	return &StatsHandler{db: db}
}

type Stats struct {
	TotalKeys     int64   `json:"total_keys"`
	ActiveKeys    int64   `json:"active_keys"`
	TotalRequests int64   `json:"total_requests"`
	ErrorRate     float64 `json:"error_rate"`
}

func (h *StatsHandler) GetStats(c *gin.Context) {
	var stats Stats

	h.db.Model(&model.APIKey{}).Count(&stats.TotalKeys)
	h.db.Model(&model.APIKey{}).Where("status = ?", "active").Count(&stats.ActiveKeys)
	h.db.Model(&model.RequestLog{}).Count(&stats.TotalRequests)

	var total, failed int64
	h.db.Model(&model.RequestLog{}).Count(&total)
	h.db.Model(&model.RequestLog{}).Where("is_success = ?", false).Count(&failed)
	if total > 0 {
		stats.ErrorRate = float64(failed) / float64(total)
	}

	utils.Success(c, stats)
}
