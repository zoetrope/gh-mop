package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/cli/go-gh"
	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		stdOut, stdErr, err := gh.Exec("issue", "view", "1", "--json", "body")
		if err != nil {
			fmt.Println(stdErr.String())
			return err
		}

		issue := struct{ Body string }{}
		err = json.Unmarshal(stdOut.Bytes(), &issue)
		if err != nil {
			return err
		}

		//fmt.Println(issue.Body)

		return parseMarkdown(([]byte)(issue.Body))
	},
}

// parseMarkdown parses the markdown string and returns a list of tasks
func parseMarkdown(source []byte) error {
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
		}
		return ast.WalkContinue, nil
	})
	fmt.Println("CodeBlocks:")
	for i, codeBlock := range codeBlocks {
		fmt.Printf("[%d] %s\n\n", i+1, codeBlock)
	}

	return err
}

func init() {
	rootCmd.AddCommand(startCmd)
}
