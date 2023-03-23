package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/cli/go-gh"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
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

func GetOperation(repo string, issue int) (*Operation, error) {
	client, err := gh.RESTClient(nil)
	if err != nil {
		return nil, err
	}
	response := &Operation{}
	err = client.Get(fmt.Sprintf("repos/%s/issues/%d", repo, issue), &response)
	if err != nil {
		return nil, err
	}

	commands, err := ParseMarkdown(([]byte)(response.Body))
	if err != nil {
		return nil, err
	}
	response.Commands = commands
	return response, nil
}

var lineRegex = regexp.MustCompile("\r\n|\n")

func ParseMarkdown(source []byte) ([]string, error) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	var codeBlocks []string
	err := ast.Walk(md.Parser().Parse(text.NewReader(source)), func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if n.Kind() == ast.KindFencedCodeBlock {
				codeBlock := n.(*ast.FencedCodeBlock)
				start := codeBlock.Lines().At(0).Start
				end := codeBlock.Lines().At(codeBlock.Lines().Len() - 1).Stop
				code := bytes.TrimSpace(source[start:end])
				codeBlocks = append(codeBlocks, string(code))
			}
			if n.Kind() == ast.KindParagraph {
				paragraph := n.(*ast.Paragraph)
				lines := paragraph.Lines()
				for i := 0; i < lines.Len(); i++ {
					line := lines.At(i)
					lineText := string(source[line.Start:line.Stop])
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
