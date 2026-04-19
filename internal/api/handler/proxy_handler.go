package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"gpt-load-mini/internal/core/keypool"
	"gpt-load-mini/internal/data/model"
	"gpt-load-mini/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProxyHandler struct {
	db          *gorm.DB
	keyProvider *keypool.Provider
	httpClient  *http.Client
}

func NewProxyHandler(db *gorm.DB, keyProvider *keypool.Provider) *ProxyHandler {
	return &ProxyHandler{
		db:          db,
		keyProvider: keyProvider,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type ProxyTestRequest struct {
	GroupName string `json:"group_name" binding:"required"`
	Path      string `json:"path" binding:"required"`
	Method    string `json:"method"`
	Body      string `json:"body"`
}

type ProxyTestResponse struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

func (h *ProxyHandler) Test(c *gin.Context) {
	var req ProxyTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	if req.Method == "" {
		req.Method = "POST"
	}

	var group model.Group
	if err := h.db.Where("name = ?", req.GroupName).First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Error(c, http.StatusNotFound, "Group not found")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "Failed to find group")
		return
	}

	apiKey, err := h.keyProvider.SelectKey(group.ID)
	if err != nil {
		utils.Error(c, http.StatusServiceUnavailable, "No active keys available: "+err.Error())
		return
	}

	upstreamURL := group.UpstreamURL
	for strings.HasSuffix(upstreamURL, "/") {
		upstreamURL = strings.TrimSuffix(upstreamURL, "/")
	}
	upstreamURL = strings.TrimSuffix(upstreamURL, "/v1")
	path := req.Path
	for strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	upstreamURL = upstreamURL + "/" + path

	var reqBody io.Reader
	if req.Body != "" {
		reqBody = bytes.NewReader([]byte(req.Body))
	} else if req.Method == "POST" {
		reqBody = bytes.NewReader([]byte(`{"model":"test","messages":[{"role":"user","content":"test"}]}`))
	}

	httpReq, err := http.NewRequestWithContext(c.Request.Context(), req.Method, upstreamURL, reqBody)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create request")
		return
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey.KeyValue))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Del("X-Api-Key")

	resp, err := h.httpClient.Do(httpReq)
	if err != nil {
		utils.Error(c, http.StatusBadGateway, "Upstream request failed: "+err.Error())
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to read response body")
		return
	}

	headers := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	h.keyProvider.UpdateKeyStatus(apiKey, group.ID, resp.StatusCode >= 200 && resp.StatusCode < 300, "")

	utils.Success(c, ProxyTestResponse{
		Status:  resp.StatusCode,
		Headers: headers,
		Body:    string(respBody),
	})
}
