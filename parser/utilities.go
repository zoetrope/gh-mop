package parser

import (
	"fmt"
	"strings"

	"github.com/cli/go-gh"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
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

func getSections(source []byte) ([]Utility, error) {
	markdown := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	var sections []Utility
	reader := text.NewReader(source)
	rootNode := markdown.Parser().Parse(reader)

	err := ast.Walk(rootNode, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && n.Kind() == ast.KindHeading {
			heading := n.(*ast.Heading)
			start := heading.Lines().At(0).Start
			end := heading.Lines().At(heading.Lines().Len() - 1).Stop
			headingText := string(source[start:end])

			var contentEnd int

			for sibling := n.NextSibling(); sibling != nil; sibling = sibling.NextSibling() {
				if sibling.Kind() == ast.KindHeading && sibling.(*ast.Heading).Level > heading.Level {
					return ast.WalkContinue, nil
				}
				if sibling.Kind() == ast.KindHeading && sibling.(*ast.Heading).Level <= heading.Level {
					break
				}
				contentEnd = sibling.Lines().At(sibling.Lines().Len() - 1).Stop
			}

			content := string(source[start:contentEnd])

			sections = append(sections, Utility{
				Title:   headingText,
				Content: strings.Repeat("#", heading.Level) + " " + content,
			})
		}
		return ast.WalkContinue, nil
	})

	if err != nil {
		return nil, err
	}

	return sections, nil
}
