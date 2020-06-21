package aslog

import (
	"bytes"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"  // mysql
	_ "github.com/jinzhu/gorm/dialects/sqlite" // sqlite
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

// Mysql struct
type Mysql struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// Option 配置选项
type Option struct {
	Mysql     Mysql
	DataDrive string // mysql,sqlite3
	FilePath  string
	FileName  string
}

// Init 使用此包 必须先初始化改
func Init(o Option) (err error) {
	var db *gorm.DB
	if o.DataDrive == "mysql" {
		db, err = gorm.Open("mysql", o.Mysql.User+":"+o.Mysql.Password+"@tcp("+o.Mysql.Host+":"+o.Mysql.Port+")/"+o.Mysql.Name+"?charset=utf8mb4&parseTime=True&loc=Local")
		if err != nil {
			return
		}
	} else {
		db, err = gorm.Open("sqlite3", o.FilePath+o.FileName)
		if err != nil {
			return
		}
	}
	if err != nil {
		return
	}
	defer db.Close()
	err = db.AutoMigrate(&HTTPLog{}, &ErrLog{}).Error
	return
}

// DB ...
func (al AsLog) DB() (db *gorm.DB, err error) {
	if al.Option.DataDrive == "mysql" {
		db, err = gorm.Open("mysql", al.Option.Mysql.User+":"+al.Option.Mysql.Password+"@tcp("+al.Option.Mysql.Host+":"+al.Option.Mysql.Port+")/"+al.Option.Mysql.Name+"?charset=utf8mb4&parseTime=True&loc=Local")
		if err != nil {
			return
		}
	} else {
		db, err = gorm.Open("sqlite3", al.Option.FilePath+al.Option.FileName)
		if err != nil {
			return
		}
	}
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
