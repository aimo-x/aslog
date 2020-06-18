package aslog

import (
	"bytes"
	"io/ioutil"
	"net/http"

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
	db, err := gorm.Open(o.FilePath, o.FileName)
	if err != nil {
		return
	}
	defer db.Close()
	err = db.AutoMigrate(&HTTPLog{}, &ErrLog{}).Error
	return
}

// DB ...
func (al AsLog) DB() (db *gorm.DB, err error) {
	db, err := gorm.Open(al.Option.FilePath.FilePath, al.Option.FileName)
	defer db.Close()
	if err != nil {
		return
	}
	return
}

// New Aslog
func New(o Option) *Aslog {
	return &AsLog{Option: &o}
}

// Write 日志写入
func (al AsLog) Write(r *http.Request, level Level, err error) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	db, err := al.DB()
	defer db.Close()
	if err != nil {
		return
	}
	var hl HTTPLog
	hl.status = r.Response.Status
	hl.Body = string(data)
	hl.Error = err.Error()
	hl.Host = r.Host
	hl.Level = level
	hl.Method = r.Method
	hl.Path = r.URL.RawPath
	hl.RawQuery = r.URL.RawQuery
	hl.Scheme = r.URL.Scheme
	body := ioutil.NopCloser(bytes.NewBuffer(data))
	r.Body = body
	err = db.Create(&hl).Error
	if err != nil {
		return
	}
	return
}
