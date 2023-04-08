package command

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zoetrope/gh-mop/pkg/github"
	"gopkg.in/h2non/gock.v1"
)

//go:embed testdata/operation.json
var operationResponse []byte

func TestGetOperation(t *testing.T) {
	t.Cleanup(gock.Off)

	gock.New("https://api.github.com").
		Get("/repos/zoetrope/gh-mop/issues/9999").
		Reply(200).
		JSON(operationResponse)

	client, err := github.NewClient("zoetrope", "gh-mop")
	op, err := GetOperation(client, 9999)
	assert.NoError(t, err)
	assert.True(t, gock.IsDone())
	assert.Equal(t, 3, len(op.Commands))
	assert.Equal(t, "kubectl get pod", op.Commands[0])
}
