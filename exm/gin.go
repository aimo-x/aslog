package main

import (
	"github.com/aimo-x/aslog"
	"github.com/gin-gonic/gin"
)

func main() {

	en := gin.New()
	alo := aslog.Option{FilePath: "./", FileName: "sqlite3.db"}
	err := aslog.Init(alo)
	if err != nil {
		panic(err)
	}
	al := aslog.New(&alo)
	en.Use(func(c *gin.Context) {
		go al.Write(c.Copy(), aslog.Debug, "nil")
	})
	en.Run(":80")
}
