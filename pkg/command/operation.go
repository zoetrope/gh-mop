package command

import (
	"encoding/json"
	"os"

	"github.com/zoetrope/gh-mop/pkg/github"
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

func GetOperation(client *github.Client, issue int) (*Operation, error) {
	content, err := client.GetIssueContent(issue)
	if err != nil {
		return nil, err
	}

	commands, err := markdown.ExtractCommands(content, func(issue int) (string, error) {
		c, err := client.GetIssueContent(issue)
		if err != nil {
			return "", err
		}
		return c, nil
	}, []int{issue})
	if err != nil {
		return nil, err
	}

	operation := &Operation{}
	operation.Commands = commands
	return operation, nil
}
