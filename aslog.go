package aslog

import (
	"bytes"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // 这是正确的包
)

// Level 日志级别 Debug、info、warn、error、fatal panic
type Level uint

// Debug ...
const Debug Level = 1

// Info ...
const Info Level = 2

// Warn ...
const Warn Level = 3

// Fatal ...
const Fatal Level = 4

// Panic ...
const Panic Level = 5

// AsLog ...
type AsLog struct {
	Option *Option
}

// Option 配置选项
type Option struct {
	FilePath string
	FileName string
}

// Init 使用此包 必须先初始化改
func Init(o Option) (err error) {
	db, err := gorm.Open("sqlite3", o.FilePath+o.FileName)
	if err != nil {
		return
	}
	defer db.Close()
	err = db.AutoMigrate(&HTTPLog{}, &ErrLog{}).Error
	return
}

// DB ...
func (al AsLog) DB() (db *gorm.DB, err error) {
	db, err = gorm.Open("sqlite3", al.Option.FilePath+al.Option.FileName)
	// defer db.Close()
	if err != nil {
		return
	}
	return
}

// New Aslog
func New(o *Option) *AsLog {
	return &AsLog{Option: o}
}

// Write 日志写入
func (al AsLog) Write(c *gin.Context, level Level, errinfo string) (err error) {
	var data []byte
	data, err = ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	db, err := al.DB()
	defer db.Close()
	if err != nil {
		return
	}
	var hl HTTPLog
	hl.Body = string(data)
	hl.Error = errinfo
	hl.Host = c.Request.Host
	hl.Level = uint(level)
	hl.Method = c.Request.Method
	hl.Path = c.Request.URL.EscapedPath()
	hl.RawQuery = c.Request.URL.RawQuery
	hl.Authorization = c.GetHeader("Authorization")
	hl.ContentType = c.GetHeader("Content-Type")
	hl.Cookie = c.GetHeader("Cookie")
	hl.UserAgent = c.Request.UserAgent()
	hl.IP = c.ClientIP()
	hl.Origin = c.GetHeader("Origin")
	body := ioutil.NopCloser(bytes.NewBuffer(data))
	c.Request.Body = body
	err = db.Create(&hl).Error
	if err != nil {
		return
	}
	return
}
