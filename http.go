package main

import (
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
)

func httpServer(c *cli.Context) {
	conf := NewConf()
	MakeDB(conf.DBUsername, conf.DBPassword, conf.DBName, conf.DBHost, conf.DBType)

	r := gin.Default()

	r.GET("/", homepage)
	r.GET("/download", download)
	r.POST("/webhook", webhook)
	r.GET("/view/:user/:reponame", serveRepo)

	r.Run(conf.ListenTo)
}
