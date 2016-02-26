package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

var validPlatforms = []string{
	"linux_386",
	"linux_amd64",
	"windows_386",
	"windows_amd64",
	"darwin_386",
	"darwin_amd64",
}

func download(c *gin.Context) {
	buildID := c.Query("id")
	platform := c.Query("p")

	if buildID == "" || platform == "" {
		c.String(400, "invalid platform or repository ID. how did you get here?")
		return
	}

	if !isValidPlatform(platform) {
		c.String(400, "not a valid platform")
		return
	}

	build := Build{}
	db.First(&build, buildID)
	if db.NewRecord(build) {
		c.String(400, "build ID does not exist")
		return
	}

	repo := Repository{}
	db.First(&repo, build.RepositoryID)
	if db.NewRecord(repo) {
		c.String(400, "repository does not exist")
		return
	}

	if build.Status != statusBuilt {
		c.String(400, "the artifact ain't ready (yet?)")
		return
	}

	ext := windowsExt(platform)
	fileName := fmt.Sprintf("data/%d/%s%s", build.ID, platform, ext)
	if _, err := os.Stat(fileName); err != nil {
		c.String(400, "whoops, we couldn't fetch that artifact. perhaps it got deleted?")
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s-%s%s", repo.Repo, platform, ext))
	c.Header("Content-Type", "application/octet-stream")

	sauceFile, err := os.Open(fileName)
	if err != nil {
		c.String(400, "can't open file")
		return
	}
	defer sauceFile.Close()

	io.Copy(c.Writer, sauceFile)
}
func isValidPlatform(a string) bool {
	for _, b := range validPlatforms {
		if b == a {
			return true
		}
	}
	return false
}
func windowsExt(platform string) string {
	if len(platform) > 7 && platform[:7] == "windows" {
		return ".exe"
	}
	return ""
}
