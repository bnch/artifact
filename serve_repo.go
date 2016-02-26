package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func serveRepo(c *gin.Context) {
	repo := Repository{}
	db.Where(&Repository{
		User: c.Param("user"),
		Repo: c.Param("reponame"),
	}).First(&repo)
	if db.NewRecord(repo) {
		c.String(400, "error: no such repository was found. How did you get here?")
		return
	}

	builds := []Build{}
	db.Where(&Build{
		RepositoryID: repo.ID,
	}).Order("id desc").Find(&builds)

	var retStr string
	for _, build := range builds {
		add := ""
		switch build.Status {
		case statusBuilding:
			add = "(still building)"
		case statusBuilt:
			add = fmt.Sprintf(
				"<a href='/download?id=%[1]d&p=linux_386'>linux_386</a> - "+
					"<a href='/download?id=%[1]d&p=linux_amd64'>linux_amd64</a> - "+
					"<a href='/download?id=%[1]d&p=windows_386'>windows_386</a> - "+
					"<a href='/download?id=%[1]d&p=windows_amd64'>windows_amd64</a> - "+
					"<a href='/download?id=%[1]d&p=darwin_386'>darwin_386 (osx)</a> - "+
					"<a href='/download?id=%[1]d&p=darwin_amd64'>darwin_amd64 (osx)</a>",
				build.ID,
			)
		case statusFailed:
			add = "(failed)"
		}
		retStr += fmt.Sprintf("commit %s %s\n", build.Commit, add)
	}
	endStr := fmt.Sprintf(`<h1>%s/%s</h1>
<pre>builds: %d

%s</pre>`, c.Param("user"), c.Param("reponame"), len(builds), retStr)

	c.Data(200, "text/html; charset=UTF-8", []byte(endStr))
}
