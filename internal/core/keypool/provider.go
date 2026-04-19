package keypool

import (
	"context"
	"fmt"
	"gpt-load-mini/internal/data/model"
	"gpt-load-mini/internal/data/store"
	"gpt-load-mini/internal/pkg/errors"
	"gpt-load-mini/internal/pkg/utils"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	StatusActive   = "active"
	StatusInvalid  = "invalid"
	maxSelectDepth = 3
)

type Provider struct {
	db        *gorm.DB
	store     store.Store
	encryptor *utils.Encryptor
}

func NewProvider(db *gorm.DB, store store.Store, encryptor *utils.Encryptor) *Provider {
	return &Provider{db: db, store: store, encryptor: encryptor}
}

func (p *Provider) SelectKey(groupID uint) (*model.APIKey, error) {
	return p.selectKey(groupID, 0)
}

func (p *Provider) selectKey(groupID uint, depth int) (*model.APIKey, error) {
	if depth >= maxSelectDepth {
		return nil, errors.ErrNoActiveKeys
	}

	keyListKey := fmt.Sprintf("group:%d:keys", groupID)

	keyIDStr, err := p.store.Rotate(keyListKey)
	if err != nil {
		if err == store.ErrNotFound {
			return nil, errors.ErrNoActiveKeys
		}
		return nil, fmt.Errorf("failed to rotate key: %w", err)
	}

	keyID, err := strconv.ParseUint(keyIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse key id: %w", err)
	}

	key, err := p.getKeyByID(uint(keyID))
	if err != nil {
		return p.selectKey(groupID, depth+1)
	}

	decryptedKey, err := p.encryptor.Decrypt(key.KeyValue)
	if err != nil {
		logrus.WithFields(logrus.Fields{"keyID": key.ID, "error": err}).Warn("failed to decrypt key, using raw value")
		decryptedKey = key.KeyValue
	}
	key.KeyValue = decryptedKey

	return key, nil
}

func (p *Provider) ReleaseKey(groupID uint, keyID uint) error {
	keyListKey := fmt.Sprintf("group:%d:keys", groupID)
	return p.store.LPush(keyListKey, keyID)
}

func (p *Provider) RemoveKey(groupID uint, keyID uint) error {
	keyListKey := fmt.Sprintf("group:%d:keys", groupID)
	if err := p.store.LRem(keyListKey, 0, keyID); err != nil {
		return err
	}
	return p.db.Model(&model.APIKey{}).Where("id = ?", keyID).Update("status", StatusInvalid).Error
}

func (p *Provider) UpdateKeyStatus(key *model.APIKey, groupID uint, success bool, errorMessage string) {
	keyID := key.ID
	groupIDVal := groupID
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		keyIDStr := fmt.Sprintf("key:%d", keyID)

		if success {
			if err := p.db.WithContext(ctx).Model(&model.APIKey{}).Where("id = ?", keyID).Updates(map[string]any{
				"failure_count": 0,
				"last_used_at":  time.Now(),
			}); err != nil {
				logrus.WithFields(logrus.Fields{"error": err, "keyID": keyID}).Error("failed to update key status in db")
			}
			if err := p.store.HSet(keyIDStr, map[string]any{"failure_count": "0"}); err != nil {
				logrus.WithFields(logrus.Fields{"error": err, "keyID": keyID}).Error("failed to update key status in store")
			}
		} else {
			newCount, err := p.store.HIncrBy(keyIDStr, "failure_count", 1)
			if err != nil {
				logrus.WithFields(logrus.Fields{"error": err, "keyID": keyID}).Error("failed to increment failure count")
			}
			if err := p.db.WithContext(ctx).Model(&model.APIKey{}).Where("id = ?", keyID).Update("failure_count", newCount); err != nil {
				logrus.WithFields(logrus.Fields{"error": err, "keyID": keyID}).Error("failed to update failure count in db")
			}

			threshold := 5
			if newCount >= int64(threshold) {
				if err := p.RemoveKey(groupIDVal, keyID); err != nil {
					logrus.WithFields(logrus.Fields{"error": err, "keyID": keyID}).Error("failed to remove key")
				}
			} else {
				if err := p.ReleaseKey(groupIDVal, keyID); err != nil {
					logrus.WithFields(logrus.Fields{"error": err, "keyID": keyID}).Error("failed to release key")
				}
			}
		}
	}()
}

func (p *Provider) LoadKeysFromDB() error {
	var keys []model.APIKey
	if err := p.db.Where("status = ?", StatusActive).Find(&keys).Error; err != nil {
		return err
	}

	keyMap := make(map[uint][]uint)
	for _, key := range keys {
		keyMap[key.GroupID] = append(keyMap[key.GroupID], key.ID)
		keyIDStr := fmt.Sprintf("key:%d", key.ID)
		p.store.HSet(keyIDStr, map[string]any{
			"id":             fmt.Sprintf("%d", key.ID),
			"key_string":     key.KeyValue,
			"status":         key.Status,
			"failure_count":  fmt.Sprintf("%d", key.FailureCount),
			"group_id":       fmt.Sprintf("%d", key.GroupID),
		})
	}

	for groupID, keyIDs := range keyMap {
		keyListKey := fmt.Sprintf("group:%d:keys", groupID)
		if err := p.store.LPush(keyListKey, toAnys(keyIDs)...); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) getKeyByID(id uint) (*model.APIKey, error) {
	var key model.APIKey
	if err := p.db.First(&key, id).Error; err != nil {
		return nil, err
	}
	return &key, nil
}

func toAnys(vals []uint) []any {
	result := make([]any, len(vals))
	for i, v := range vals {
		result[i] = v
	}
	return result
}
