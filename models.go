package aslog

import (
	"time"
)

// HTTPLog http日志记录
type HTTPLog struct {
	Model
	Level    uint   `gorm:"index"` // 1,2,3,4,5,6,7,8,9
	Host     string `gorm:"index"`
	Path     string `gorm:"index"`
	Method   string `gorm:"index"`
	Scheme   string
	Body     string `gorm:"size:2048"`
	RawQuery string `gorm:"size:2048"`
	status   int
	Error    string `gorm:"size:2048"` // 错误描述
}

// ErrLog 错误日志
type ErrLog struct {
	Model
	Level      uint `gorm:"index"` // 1,2,3,4,5,6,7,8,9
	HTTPLogID  uint `gorm:"index"`
	HTTPStatus int
	Error      string `gorm:"size:2048"` // 错误描述
}

// Model ...
type Model struct {
	ID        uint `gorm:"primary_key" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
