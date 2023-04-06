package command

import (
	"github.com/zoetrope/gh-mop/pkg/github"
	"github.com/zoetrope/gh-mop/pkg/markdown"
)

type Utility struct {
	Title    string
	Content  string
	Commands []string
}

func GetUtilities(client *github.Client, issue int) ([]Utility, error) {
	content, err := client.GetIssueContent(issue)
	if err != nil {
		return nil, err
	}

	sections, err := markdown.ExtractSections(([]byte)(content))
	if err != nil {
		return nil, err
	}

	for i, section := range sections {
		commands, err := markdown.ExtractCommands(([]byte)(section.Content))
		if err != nil {
			return nil, err
		}
		sections[i].Commands = commands
	}

	return sections, nil
}
