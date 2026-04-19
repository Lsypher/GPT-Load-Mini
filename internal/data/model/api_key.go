package model

import "time"

type APIKey struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	GroupID      uint      `gorm:"index" json:"group_id"`
	KeyValue     string    `gorm:"type:text" json:"key_value"`
	KeyHash      string    `gorm:"type:varchar(128);index" json:"key_hash"`
	Status       string    `gorm:"type:varchar(50);default:active" json:"status"`
	FailureCount int64     `gorm:"default:0" json:"failure_count"`
	LastUsedAt   *time.Time `json:"last_used_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
