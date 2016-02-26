package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"net/url"
	"strings"
)

func addRepo(c *cli.Context) {
	conf := NewConf()
	err := MakeDB(conf.DBUsername, conf.DBPassword, conf.DBName, conf.DBHost, conf.DBType)
	if err != nil {
		panic(err)
	}
	repoURL, err := url.Parse(c.Args().First())
	if err != nil {
		panic(err)
	}
	repoInfo := strings.Split(repoURL.Path, "/")
	if len(repoInfo) < 3 {
		panic("must give repo path having username and repo name")
	}
	repoSecret := genRandomString(100)
	repo := Repository{
		User:      repoInfo[1],
		Repo:      repoInfo[2],
		SecretKey: repoSecret,
	}
	db.Create(&repo)
	fmt.Printf("=== CREATED NEW REPO===\nSecret key: %s\nID: %d\n", repoSecret, repo.ID)
}
