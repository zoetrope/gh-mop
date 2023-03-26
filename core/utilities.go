package core

import (
	"fmt"

	"github.com/cli/go-gh"
)

type Utility struct {
	Title    string
	Content  string
	Commands []string
}

func GetUtilities(owner, repo string, issue int) ([]Utility, error) {
	client, err := gh.RESTClient(nil)
	if err != nil {
		return nil, err
	}
	response := &struct {
		Body string `json:"body"`
	}{}
	err = client.Get(fmt.Sprintf("repos/%s/%s/issues/%d", owner, repo, issue), &response)
	if err != nil {
		return nil, err
	}

	sections, err := getSections(([]byte)(response.Body))
	if err != nil {
		return nil, err
	}

	for i, section := range sections {
		commands, err := getCommands(([]byte)(section.Content))
		if err != nil {
			return nil, err
		}
		sections[i].Commands = commands
	}

	return sections, nil
}
