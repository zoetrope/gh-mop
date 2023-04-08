package command

import (
	_ "embed"
	"testing"

	"github.com/cli/go-gh/pkg/config"

	"github.com/stretchr/testify/assert"
	"github.com/zoetrope/gh-mop/pkg/github"
	"gopkg.in/h2non/gock.v1"
)

//go:embed testdata/expect1.txt
var expectTypescript1 string

//go:embed testdata/expect2.txt
var expectTypescript2 string

func TestUploadResult(t *testing.T) {
	stubConfig(t)
	t.Cleanup(gock.Off)
	//gock.Observe(gock.DumpRequest)

	gock.New("https://api.github.com").
		Post("/repos/zoetrope/gh-mop/issues/9999/comments").
		MatchType("application/json; charset=utf-8").
		JSON(map[string]string{"body": expectTypescript1}).
		Reply(200).
		BodyString(`{"html_url": "https://github.com/zoetrope/gh-mop/issues/9999#issuecomment-1234"}`)

	client, err := github.NewClient("zoetrope", "gh-mop")
	assert.NoError(t, err)

	url, readBytes, err := UploadResult(client, 9999, "./testdata/input1.txt", 0)
	assert.NoError(t, err)
	assert.True(t, gock.IsDone())
	assert.Equal(t, int64(52), readBytes)
	assert.Equal(t, "https://github.com/zoetrope/gh-mop/issues/9999#issuecomment-1234", url)
}

func TestUploadResultWithOffset(t *testing.T) {
	stubConfig(t)
	t.Cleanup(gock.Off)
	//gock.Observe(gock.DumpRequest)

	gock.New("https://api.github.com").
		Post("/repos/zoetrope/gh-mop/issues/9999/comments").
		MatchType("application/json; charset=utf-8").
		JSON(map[string]string{"body": expectTypescript2}).
		Reply(200).
		BodyString(`{"html_url": "https://github.com/zoetrope/gh-mop/issues/9999#issuecomment-1235"}`)

	client, err := github.NewClient("zoetrope", "gh-mop")
	assert.NoError(t, err)

	url, readBytes, err := UploadResult(client, 9999, "./testdata/input2.txt", 52)
	assert.NoError(t, err)
	assert.True(t, gock.IsDone())
	assert.Equal(t, int64(75), readBytes)
	assert.Equal(t, "https://github.com/zoetrope/gh-mop/issues/9999#issuecomment-1235", url)
}

func stubConfig(t *testing.T) {
	t.Helper()
	old := config.Read
	cfg := `
hosts:
  github.com:
    user: dummy
    oauth_token: abc123
    git_protocol: ssh
`
	config.Read = func() (*config.Config, error) {
		return config.ReadFromString(cfg), nil
	}
	t.Cleanup(func() {
		config.Read = old
	})
}
