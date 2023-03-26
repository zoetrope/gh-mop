package core

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func getSections(markdownText []byte) ([]Utility, error) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	var sections []Utility
	reader := text.NewReader(markdownText)
	rootNode := md.Parser().Parse(reader)

	err := ast.Walk(rootNode, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && n.Kind() == ast.KindHeading {
			heading := n.(*ast.Heading)
			start := heading.Lines().At(0).Start
			end := heading.Lines().At(heading.Lines().Len() - 1).Stop
			headingText := string(markdownText[start:end])

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

			content := string(markdownText[start:contentEnd])

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
func getCommands(markdownText []byte) ([]string, error) {
	var lineRegex = regexp.MustCompile("\r\n|\n")

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	var codeBlocks []string
	err := ast.Walk(md.Parser().Parse(text.NewReader(markdownText)), func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if n.Kind() == ast.KindFencedCodeBlock {
				codeBlock := n.(*ast.FencedCodeBlock)
				start := codeBlock.Lines().At(0).Start
				end := codeBlock.Lines().At(codeBlock.Lines().Len() - 1).Stop
				code := bytes.TrimSpace(markdownText[start:end])
				codeBlocks = append(codeBlocks, string(code))
			}
			if n.Kind() == ast.KindParagraph {
				paragraph := n.(*ast.Paragraph)
				lines := paragraph.Lines()
				for i := 0; i < lines.Len(); i++ {
					line := lines.At(i)
					lineText := string(markdownText[line.Start:line.Stop])
					if strings.HasPrefix(strings.TrimSpace(lineText), "ref:") {
						//DetectCommand()
					}
				}
			}
		}
		return ast.WalkContinue, nil
	})

	var commands []string
	for _, codeBlock := range codeBlocks {
		var command string
		for _, line := range lineRegex.Split(codeBlock, -1) {
			if strings.HasPrefix(line, "#") {
				continue
			}
			if strings.TrimSpace(line) == "" {
				continue
			}
			if strings.HasPrefix(line, "$") {
				line = strings.TrimLeft(line, "$ ")
			}

			command += line
			if strings.HasSuffix(line, "\\") {
				command += "\n"
				continue
			}
			commands = append(commands, command)
			command = ""
		}
	}

	return commands, err
}
