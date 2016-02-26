package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func homepage(c *gin.Context) {
	repositories := []Repository{}
	db.Find(&repositories)

	out := "<h1>Howl's super lazy artifact building system</h1><br><pre>Repos which artifacts are built on this machine:\n"
	for _, repo := range repositories {
		out += fmt.Sprintf("<a href='/view/%s/%s'>%[1]s/%[2]s</a>\n", repo.User, repo.Repo)
	}
	out += "</pre>"

	c.Data(200, "text/html; charset=UTF-8", []byte(out))
}
