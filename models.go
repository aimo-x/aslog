package aslog

import (
	"time"
)

// HTTPLog http日志记录
type HTTPLog struct {
	Model
	Level         uint   `gorm:"index"` // 1,2,3,4,5,6,7,8,9
	Host          string `gorm:"index"`
	Path          string `gorm:"index"`
	Method        string `gorm:"index"`
	Body          string `gorm:"type:text(15535)"`
	RawQuery      string `gorm:"type:text(15535)"`
	status        int
	Error         string `gorm:"type:text(15535)"` // 错误描述
	Authorization string `gorm:"type:text(15535)"`
	ContentType   string `gorm:"type:text(15535)"`
	Origin        string `gorm:"type:text(15535)"`
	Cookie        string `gorm:"type:text(15535)"`
	UserAgent     string `gorm:"type:text(15535)"`
	IP            string
}

// ErrLog 错误日志
type ErrLog struct {
	Model
	Level      uint `gorm:"index"` // 1,2,3,4,5,6,7,8,9
	HTTPLogID  uint `gorm:"index"`
	HTTPStatus int
	Error      string `gorm:"type:varchar(255)"` // 错误描述
}

// Model ...
type Model struct {
	ID        uint `gorm:"primary_key" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
