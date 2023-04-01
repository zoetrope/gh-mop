package core

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

//go:embed testdata/expect_typescript1.txt
var expectTypescript1 string

//go:embed testdata/expect_typescript2.txt
var expectTypescript2 string

func TestUploadResult(t *testing.T) {
	t.Cleanup(gock.Off)
	//gock.Observe(gock.DumpRequest)

	gock.New("https://api.github.com").
		Post("/repos/zoetrope/gh-mop/issues/9999/comments").
		MatchType("application/json; charset=utf-8").
		JSON(map[string]string{"body": expectTypescript1}).
		Reply(200).
		BodyString(`{"html_url": "https://github.com/zoetrope/gh-mop/issues/9999#issuecomment-1234"}`)

	client, err := NewClient("zoetrope", "gh-mop")

	readBytes, err := UploadResult(client, 9999, "./testdata/typescript1.txt", 0, true)
	assert.NoError(t, err)
	assert.True(t, gock.IsDone())
	assert.Equal(t, int64(222), readBytes)
}

func TestUploadResultWithOffset(t *testing.T) {
	t.Cleanup(gock.Off)
	//gock.Observe(gock.DumpRequest)

	gock.New("https://api.github.com").
		Post("/repos/zoetrope/gh-mop/issues/9999/comments").
		MatchType("application/json; charset=utf-8").
		JSON(map[string]string{"body": expectTypescript2}).
		Reply(200).
		BodyString(`{"html_url": "https://github.com/zoetrope/gh-mop/issues/9999#issuecomment-1235"}`)

	client, err := NewClient("zoetrope", "gh-mop")

	readBytes, err := UploadResult(client, 9999, "./testdata/typescript2.txt", 222, true)
	assert.NoError(t, err)
	assert.True(t, gock.IsDone())
	assert.Equal(t, int64(346), readBytes)
}
