package model

import "time"

type Group struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);uniqueIndex" json:"name"`
	DisplayName string    `gorm:"type:varchar(255)" json:"display_name"`
	ChannelType string    `gorm:"type:varchar(50)" json:"channel_type"`
	UpstreamURL string    `gorm:"type:varchar(512)" json:"upstream_url"`
	TestModel   string    `gorm:"type:varchar(255)" json:"test_model"`
	Sort        int       `gorm:"default:0" json:"sort"`
	ProxyAPIKey string    `gorm:"type:varchar(255)" json:"proxy_api_key"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
