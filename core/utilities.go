package core

type Utility struct {
	Title    string
	Content  string
	Commands []string
}

func GetUtilities(client *Client, issue int) ([]Utility, error) {
	content, err := client.GetIssueContent(issue)
	if err != nil {
		return nil, err
	}

	sections, err := getSections(([]byte)(content))
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
