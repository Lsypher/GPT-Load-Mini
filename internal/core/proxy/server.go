package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"gpt-load-mini/internal/core/keypool"
	"gpt-load-mini/internal/data/model"
	"gpt-load-mini/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProxyServer struct {
	keyProvider *keypool.Provider
	db          *gorm.DB
	upstreamURL string
	maxRetries  int
	httpClient  *http.Client
}

func NewProxyServer(keyProvider *keypool.Provider, db *gorm.DB, upstreamURL string, maxRetries int) *ProxyServer {
	return &ProxyServer{
		keyProvider: keyProvider,
		db:          db,
		upstreamURL: upstreamURL,
		maxRetries:  maxRetries,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (s *ProxyServer) HandleProxy(c *gin.Context) {
	groupName := c.Param("group_name")
	path := c.Param("path")
	fullPath := path

	startTime := time.Now()

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		utils.Error(c, 400, "Failed to read request body")
		return
	}

	var group model.Group
	if err := s.db.Where("name = ?", groupName).First(&group).Error; err != nil {
		utils.Error(c, 404, "Group not found")
		return
	}

	// Authenticate with ProxyAPIKey if set
	if group.ProxyAPIKey != "" {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || authHeader != "Bearer "+group.ProxyAPIKey {
			utils.Error(c, 401, "Invalid or missing API key")
			return
		}
	}

	var lastErr error
	var usedKey *model.APIKey

	for attempt := 0; attempt < s.maxRetries; attempt++ {
		apiKey, err := s.keyProvider.SelectKey(group.ID)
		if err != nil {
			s.logRequest(c, nil, groupName, startTime, 503, err.Error(), false, "final", group.ID, bodyBytes)
			utils.Error(c, 503, err.Error())
			return
		}

		upstreamURL := group.UpstreamURL
		// Remove /v1 suffix and all trailing slashes from upstream
		for strings.HasSuffix(upstreamURL, "/") {
			upstreamURL = strings.TrimSuffix(upstreamURL, "/")
		}
		upstreamURL = strings.TrimSuffix(upstreamURL, "/v1")
		// Remove leading slashes from path
		for strings.HasPrefix(fullPath, "/") {
			fullPath = fullPath[1:]
		}
		upstreamURL = upstreamURL + "/" + fullPath
		req, err := http.NewRequestWithContext(c.Request.Context(), "POST", upstreamURL, bytes.NewReader(bodyBytes))
		if err != nil {
			utils.Error(c, 500, "Failed to create upstream request")
			return
		}

		req.Header = c.Request.Header.Clone()
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey.KeyValue))
		req.Header.Del("X-Api-Key")

		resp, err := s.httpClient.Do(req)
		if err != nil {
			s.keyProvider.UpdateKeyStatus(apiKey, apiKey.GroupID, false, err.Error())
			lastErr = err
			usedKey = apiKey
			continue
		}
		if resp != nil {
			defer resp.Body.Close()
		}

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.Error(c, 500, "Failed to read upstream response")
			return
		}

		isSuccess := resp.StatusCode >= 200 && resp.StatusCode < 300
		if !isSuccess {
			s.keyProvider.UpdateKeyStatus(apiKey, apiKey.GroupID, false, string(respBody))
			lastErr = fmt.Errorf("upstream returned status %d", resp.StatusCode)
			usedKey = apiKey
			continue
		} else {
			s.keyProvider.UpdateKeyStatus(apiKey, apiKey.GroupID, true, "")
		}

		for k, v := range resp.Header {
			for _, val := range v {
				c.Header(k, val)
			}
		}
		c.Header("X-Accel-Buffering", "no")
		c.Data(resp.StatusCode, "application/json", respBody)

		s.logRequest(c, apiKey, groupName, startTime, resp.StatusCode, "", false, "final", group.ID, bodyBytes)
		return
	}

	// All retries exhausted
	s.keyProvider.UpdateKeyStatus(usedKey, group.ID, false, lastErr.Error())
	s.logRequest(c, usedKey, groupName, startTime, 502, lastErr.Error(), false, "final", group.ID, bodyBytes)
	utils.Error(c, 502, fmt.Sprintf("All retries failed: %v", lastErr))
}

func (s *ProxyServer) logRequest(c *gin.Context, key *model.APIKey, groupName string, startTime time.Time, statusCode int, errMsg string, isStream bool, reqType string, groupID uint, bodyBytes []byte) {
	// Capture values before spawning goroutine to avoid use of invalid gin.Context
	clientIP := c.ClientIP()
	requestPath := c.Request.URL.Path
	keyID := uint(0)
	if key != nil {
		keyID = key.ID
	}

	go func() {
		ctx := context.Background()
		modelName := extractModelFromBody(bodyBytes)

		logEntry := &model.RequestLog{
			ID:           uuid.New().String(),
			Timestamp:    startTime,
			GroupID:      groupID,
			GroupName:    groupName,
			KeyID:        keyID,
			Model:        modelName,
			IsSuccess:    statusCode >= 200 && statusCode < 400,
			SourceIP:     clientIP,
			StatusCode:   statusCode,
			RequestPath:  requestPath,
			DurationMs:   time.Since(startTime).Milliseconds(),
			ErrorMessage: errMsg,
			IsStream:     isStream,
			RequestType:  reqType,
		}

		if err := s.db.WithContext(ctx).Create(logEntry).Error; err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("failed to persist request log")
		}

		logrus.WithFields(logrus.Fields{
			"group":    groupName,
			"status":   statusCode,
			"duration": logEntry.DurationMs,
		}).Info("proxy request")
	}()
}

func extractModelFromBody(body []byte) string {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return ""
	}
	if model, ok := payload["model"].(string); ok {
		return model
	}
	return ""
}
