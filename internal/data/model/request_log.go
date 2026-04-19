package model

import "time"

type RequestLog struct {
	ID           string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Timestamp    time.Time `gorm:"index" json:"timestamp"`
	GroupID      uint      `gorm:"index" json:"group_id"`
	GroupName    string    `gorm:"type:varchar(255)" json:"group_name"`
	KeyID        uint      `json:"key_id"`
	Model        string    `gorm:"type:varchar(255)" json:"model"`
	IsSuccess    bool      `json:"is_success"`
	SourceIP     string    `gorm:"type:varchar(64)" json:"source_ip"`
	StatusCode   int       `json:"status_code"`
	RequestPath  string    `gorm:"type:varchar(500)" json:"request_path"`
	DurationMs   int64     `json:"duration_ms"`
	ErrorMessage string    `gorm:"type:text" json:"error_message"`
	IsStream     bool      `json:"is_stream"`
	RequestType  string    `gorm:"type:varchar(20);default:final" json:"request_type"`
}
