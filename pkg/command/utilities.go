package command

import (
	"github.com/zoetrope/gh-mop/pkg/github"
	"github.com/zoetrope/gh-mop/pkg/markdown"
)

func GetUtilities(client *github.Client, issue int) ([]markdown.Section, error) {
	content, err := client.GetIssueContent(issue)
	if err != nil {
		return nil, err
	}

	sections, err := markdown.ExtractSections(([]byte)(content))
	if err != nil {
		return nil, err
	}

	for i, section := range sections {
		commands, err := markdown.ExtractCommands(section.Content, func(issue int) (string, error) {
			return "", nil
		}, []int{issue})
		if err != nil {
			return nil, err
		}
		sections[i].Commands = commands
	}

	return sections, nil
}
