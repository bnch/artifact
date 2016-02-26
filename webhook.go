package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func webhook(c *gin.Context) {
	payloadRaw, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
		c.String(400, "error")
		return
	}

	if c.Query("repoID") == "" {
		c.String(400, "error")
		return
	}

	repoID, err := strconv.Atoi(c.Query("repoID"))
	if err != nil {
		fmt.Println(err)
		c.String(400, "error")
		return
	}

	repo := Repository{}
	db.First(&repo, repoID)
	if db.NewRecord(repo) {
		c.String(400, "no such repo")
		return
	}

	if !verifySignature([]byte(repo.SecretKey), c.Request.Header.Get("X-Hub-Signature"), payloadRaw) {
		c.String(400, "invalid signature")
		return
	}

	payload := pushData{}
	err = json.Unmarshal(payloadRaw, &payload)
	if err != nil {
		fmt.Println(err)
		c.String(400, "error")
		return
	}

	go job(payload.HeadCommit.ID, repo)

	c.String(200, "job enqueued (commit: %s)", payload.HeadCommit.ID)

	fmt.Printf("%#v\n", payload)
}

type pushData struct {
	HeadCommit struct {
		ID string `json:"id"`
	} `json:"head_commit"`
}

func verifySignature(secret []byte, signature string, body []byte) bool {

	const signaturePrefix = "sha1="
	const signatureLength = 45 // len(SignaturePrefix) + len(hex(sha1))

	if len(signature) != signatureLength || !strings.HasPrefix(signature, signaturePrefix) {
		return false
	}

	actual := make([]byte, 20)
	hex.Decode(actual, []byte(signature[5:]))

	return hmac.Equal(signBody(secret, body), actual)
}
func signBody(secret, body []byte) []byte {
	computed := hmac.New(sha1.New, secret)
	computed.Write(body)
	return []byte(computed.Sum(nil))
}

func job(commitID string, repo Repository) {
	build := Build{
		RepositoryID: repo.ID,
		Commit:       commitID,
		Status:       statusBuilding,
	}
	db.Create(&build)

	repoNameForGoGet := fmt.Sprintf("github.com/%s/%s", repo.User, repo.Repo)

	_, err := exec.Command("go", "get", "-d", "-u", repoNameForGoGet).Output()
	if err != nil {
		fmt.Println("job terminated prematurely")
		fmt.Println(err)
		build.Status = statusFailed
		db.Update(&build)
		return
	}

	buildDir := fmt.Sprintf("data/%d", build.ID)
	err = os.MkdirAll(buildDir, 755)
	if err != nil {
		fmt.Println("job terminated prematurely")
		fmt.Println(err)
		build.Status = statusFailed
		db.Update(&build)
		return
	}

	// BUILD IT!
	for _, goos := range []string{"darwin", "linux", "windows"} {
		for _, arch := range []string{"386", "amd64"} {
			ext := ""
			if goos == "windows" {
				ext = ".exe"
			}
			c := exec.Command("go", "build", "-o", buildDir+"/"+goos+"_"+arch+ext)
			c.Env = append(os.Environ(),
				"GOOS="+goos,
				"GOARCH="+arch,
			)
			_, err := c.Output()
			if err != nil {
				fmt.Println("job terminated prematurely")
				fmt.Println(err)
				build.Status = statusFailed
				db.Update(&build)
				return
			}
		}
	}

	build.Status = statusBuilt
	db.Update(&build)
}
