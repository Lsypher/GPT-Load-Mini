package handler

import (
	"net/http"
	"strconv"

	"gpt-load-mini/internal/data/model"
	"gpt-load-mini/internal/data/store"
	"gpt-load-mini/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type KeyHandler struct {
	db        *gorm.DB
	encryptor *utils.Encryptor
	store     store.Store
}

func NewKeyHandler(db *gorm.DB, encryptor *utils.Encryptor, store store.Store) *KeyHandler {
	return &KeyHandler{db: db, encryptor: encryptor, store: store}
}

func (h *KeyHandler) Add(c *gin.Context) {
	var input struct {
		GroupID  uint   `json:"group_id" binding:"required"`
		KeyValue string `json:"key_value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	encrypted, err := h.encryptor.Encrypt(input.KeyValue)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to encrypt key")
		return
	}

	key := &model.APIKey{
		GroupID:  input.GroupID,
		KeyValue: encrypted,
		KeyHash:  h.encryptor.Hash(input.KeyValue),
		Status:   "active",
	}

	if err := h.db.Create(key).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to add key")
		return
	}

	keyListKey := store.FormatGroupKeysKey(input.GroupID)
	if err := h.store.LPush(keyListKey, key.ID); err != nil {
		logrus.WithFields(logrus.Fields{"error": err, "keyID": key.ID}).Error("failed to push key to store")
	}

	utils.Success(c, key)
}

func (h *KeyHandler) List(c *gin.Context) {
	groupIDStr := c.Query("group_id")
	var keys []model.APIKey
	query := h.db.Model(&model.APIKey{})

	if groupIDStr != "" {
		groupID, _ := strconv.ParseUint(groupIDStr, 10, 64)
		query = query.Where("group_id = ?", groupID)
	}

	if err := query.Find(&keys).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to list keys")
		return
	}
	utils.Success(c, keys)
}

func (h *KeyHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	var key model.APIKey
	if err := h.db.First(&key, id).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to find key")
		return
	}

	if err := h.db.Delete(&model.APIKey{}, id).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to delete key")
		return
	}

	keyListKey := store.FormatGroupKeysKey(key.GroupID)
	if err := h.store.LRem(keyListKey, 0, key.ID); err != nil {
		logrus.WithFields(logrus.Fields{"error": err, "keyID": key.ID}).Error("failed to remove key from store")
	}
	utils.SuccessWithMessage(c, "Key deleted")
}

func (h *KeyHandler) Restore(c *gin.Context) {
	id := c.Param("id")
	var key model.APIKey
	if err := h.db.First(&key, id).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to find key")
		return
	}

	if err := h.db.Model(&model.APIKey{}).Where("id = ?", id).Updates(map[string]any{
		"status":        "active",
		"failure_count": 0,
	}).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to restore key")
		return
	}

	keyListKey := store.FormatGroupKeysKey(key.GroupID)
	if err := h.store.LPush(keyListKey, key.ID); err != nil {
		logrus.WithFields(logrus.Fields{"error": err, "keyID": key.ID}).Error("failed to push key to store")
	}
	utils.SuccessWithMessage(c, "Key restored")
}

func (h *KeyHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var key model.APIKey
	if err := h.db.First(&key, id).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to find key")
		return
	}

	var input struct {
		GroupID  *uint  `json:"group_id"`
		KeyValue string `json:"key_value"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	updates := map[string]any{}
	if input.GroupID != nil {
		updates["group_id"] = *input.GroupID
	}
	if input.KeyValue != "" {
		encrypted, err := h.encryptor.Encrypt(input.KeyValue)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "Failed to encrypt key")
			return
		}
		updates["key_value"] = encrypted
		updates["key_hash"] = h.encryptor.Hash(input.KeyValue)
	}

	if err := h.db.Model(&model.APIKey{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update key")
		return
	}

	h.db.First(&key, id)
	utils.Success(c, key)
}

type ExportKey struct {
	GroupID   uint   `json:"group_id"`
	KeyHash   string `json:"key_hash"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func (h *KeyHandler) Export(c *gin.Context) {
	var keys []model.APIKey
	query := h.db.Model(&model.APIKey{})

	if groupIDStr := c.Query("group_id"); groupIDStr != "" {
		groupID, _ := strconv.ParseUint(groupIDStr, 10, 64)
		query = query.Where("group_id = ?", groupID)
	}

	if err := query.Find(&keys).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to list keys")
		return
	}

	exportKeys := make([]ExportKey, len(keys))
	for i, k := range keys {
		exportKeys[i] = ExportKey{
			GroupID:   k.GroupID,
			KeyHash:   k.KeyHash,
			Status:    k.Status,
			CreatedAt: k.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	utils.Success(c, exportKeys)
}

type ImportKey struct {
	GroupID  uint   `json:"group_id" binding:"required"`
	KeyValue string `json:"key_value" binding:"required"`
}

func (h *KeyHandler) Import(c *gin.Context) {
	var input []ImportKey
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	imported := 0
	failed := 0
	for _, item := range input {
		encrypted, err := h.encryptor.Encrypt(item.KeyValue)
		if err != nil {
			failed++
			continue
		}

		key := &model.APIKey{
			GroupID:  item.GroupID,
			KeyValue: encrypted,
			KeyHash:  h.encryptor.Hash(item.KeyValue),
			Status:   "active",
		}

		if err := h.db.Create(key).Error; err != nil {
			failed++
			continue
		}

		keyListKey := store.FormatGroupKeysKey(item.GroupID)
		if err := h.store.LPush(keyListKey, key.ID); err != nil {
			logrus.WithFields(logrus.Fields{"error": err, "keyID": key.ID}).Error("failed to push key to store")
		}
		imported++
	}

	utils.Success(c, map[string]int{"imported": imported, "failed": failed})
}
