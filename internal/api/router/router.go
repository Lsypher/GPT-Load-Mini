package router

import (
	"net/http"

	"gpt-load-mini/internal/api/handler"
	"gpt-load-mini/internal/core/proxy"
	"gpt-load-mini/internal/pkg/config"
	"gpt-load-mini/internal/pkg/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	groupHandler *handler.GroupHandler,
	keyHandler *handler.KeyHandler,
	statsHandler *handler.StatsHandler,
	logHandler *handler.LogHandler,
	healthHandler *handler.HealthHandler,
	proxyServer *proxy.ProxyServer,
	cfg *config.Config,
	configHandler *handler.ConfigHandler,
	proxyHandler *handler.ProxyHandler,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "X-Auth-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	api.Use(authMiddleware(cfg.AuthKey))
	{
		api.POST("/groups", groupHandler.Create)
		api.GET("/groups", groupHandler.List)
		api.GET("/groups/name/:name", groupHandler.GetByName)
		api.GET("/groups/:id", groupHandler.Get)
		api.PUT("/groups/:id", groupHandler.Update)
		api.DELETE("/groups/:id", groupHandler.Delete)
		api.POST("/keys", keyHandler.Add)
		api.GET("/keys", keyHandler.List)
		api.PUT("/keys/:id", keyHandler.Update)
		api.DELETE("/keys/:id", keyHandler.Delete)
		api.POST("/keys/:id/restore", keyHandler.Restore)
		api.GET("/keys/export", keyHandler.Export)
		api.POST("/keys/import", keyHandler.Import)

		api.GET("/stats", statsHandler.GetStats)

		api.GET("/logs", logHandler.List)

		api.POST("/admin/reload-config", configHandler.ReloadConfig)

		api.POST("/proxy/test", proxyHandler.Test)
	}

	r.GET("/api/health", healthHandler.Health)

	proxyGroup := r.Group("/proxy/:group_name")
	proxyGroup.Any("/*path", proxyServer.HandleProxy)

	return r
}

func authMiddleware(authKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if authKey == "" {
			c.Next()
			return
		}
		key := c.GetHeader("X-Auth-Key")
		if key != authKey {
			utils.Error(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}
