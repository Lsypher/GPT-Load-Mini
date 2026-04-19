package handler

import (
	"gpt-load-mini/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(c *gin.Context) {
	utils.Success(c, gin.H{"status": "ok"})
}
