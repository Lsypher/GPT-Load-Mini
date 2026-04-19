package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"gpt-load-mini/internal/api/handler"
	"gpt-load-mini/internal/core/keypool"
	"gpt-load-mini/internal/core/proxy"
	"gpt-load-mini/internal/data/db"
	"gpt-load-mini/internal/data/store"
	"gpt-load-mini/internal/pkg/config"
	"gpt-load-mini/internal/pkg/utils"
	"gpt-load-mini/internal/api/router"

	"github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load()

	cfg := config.Load()

	container, err := buildContainer(cfg)
	if err != nil {
		log.Fatalf("Failed to build container: %v", err)
	}

	var srv struct {
		Router *gin.Engine
	}
	if err := container.Invoke(func(s *gin.Engine) { srv.Router = s }); err != nil {
		log.Fatalf("Failed to get router: %v", err)
	}

	var db *gorm.DB
	container.Invoke(func(d *gorm.DB) { db = d })

	var redisStore store.Store
	container.Invoke(func(s store.Store) { redisStore = s })

	httpSrv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ServerPort),
		Handler: srv.Router,
	}

	go func() {
		logrus.Infof("Server starting on :%s", cfg.ServerPort)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(ctx); err != nil {
		logrus.WithError(err).Error("Server forced to shutdown")
	}

	sqlDB, err := db.DB()
	if err == nil {
		if err := sqlDB.Close(); err != nil {
			logrus.WithError(err).Error("Failed to close database")
		}
	}

	if err := redisStore.Close(); err != nil {
		logrus.WithError(err).Error("Failed to close redis store")
	}

	logrus.Info("Server exited")
}

func buildContainer(cfg *config.Config) (*dig.Container, error) {
	container := dig.New()

	container.Provide(func() *config.Config { return cfg })

	database, err := db.NewDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	container.Provide(func() *gorm.DB { return database })

	redisStore, err := store.NewRedisStore(fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}
	container.Provide(func() store.Store { return redisStore })

	encryptor := utils.NewEncryptor(cfg.EncryptionKey)
	container.Provide(func() *utils.Encryptor { return encryptor })

	container.Provide(handler.NewGroupHandler)
	container.Provide(handler.NewKeyHandler)
	container.Provide(handler.NewStatsHandler)
	container.Provide(handler.NewLogHandler)
	container.Provide(handler.NewHealthHandler)
	container.Provide(handler.NewProxyHandler)
	container.Provide(func(cfg *config.Config) *handler.ConfigHandler {
		return handler.NewConfigHandler(cfg)
	})

	container.Provide(keypool.NewProvider)

	container.Provide(func(kp *keypool.Provider, db *gorm.DB) *proxy.ProxyServer {
		return proxy.NewProxyServer(kp, db, "https://api.openai.com", 3)
	})

	container.Provide(router.NewRouter)

	return container, nil
}
