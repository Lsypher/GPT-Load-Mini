package handler

import (
	"gpt-load-mini/internal/pkg/config"
	"gpt-load-mini/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

type ConfigHandler struct {
	cfg *config.Config
}

func NewConfigHandler(cfg *config.Config) *ConfigHandler {
	return &ConfigHandler{cfg: cfg}
}

func (h *ConfigHandler) ReloadConfig(c *gin.Context) {
	h.cfg.Reload()
	utils.SuccessWithMessage(c, "Configuration reloaded")
}
