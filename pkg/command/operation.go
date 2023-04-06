package command

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cli/go-gh"
	"github.com/zoetrope/gh-mop/pkg/markdown"
)

type Operation struct {
	Body     string   `json:"body"`
	Commands []string `json:"commands"`
}

func LoadOperation(filepath string) (*Operation, error) {
	b, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	op := &Operation{}
	err = json.Unmarshal(b, op)
	if err != nil {
		return nil, err
	}
	return op, nil
}

func SaveOperation() {

}

func GetOperation(owner, repo string, issue int) (*Operation, error) {
	client, err := gh.RESTClient(nil)
	if err != nil {
		return nil, err
	}
	response := &Operation{}
	err = client.Get(fmt.Sprintf("repos/%s/%s/issues/%d", owner, repo, issue), &response)
	if err != nil {
		return nil, err
	}

	commands, err := markdown.ExtractCommands(([]byte)(response.Body))
	if err != nil {
		return nil, err
	}
	response.Commands = commands
	return response, nil
}
