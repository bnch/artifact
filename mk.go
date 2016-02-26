package main

import (
	"github.com/codegangsta/cli"
)

func mkconf(c *cli.Context) {
	NewConf()
}
func mkdb(c *cli.Context) {
	conf := NewConf()
	err := MakeDB(conf.DBUsername, conf.DBPassword, conf.DBName, conf.DBHost, conf.DBType)
	if err != nil {
		panic(err)
	}
	db.CreateTable(&Repository{}, &Build{})
}
