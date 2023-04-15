package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/charmbracelet/glamour"
	"github.com/cli/go-gh/pkg/markdown"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
	"github.com/zoetrope/gh-mop/pkg/command"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list ISSUE_NUMBER",
	Short: "List commands in the current operation",
	Long: `List commands in the current operation using fuzzy search like fzf.

Arguments:
  ISSUE_NUMBER: Issue number defining an operation on GitHub

Examples:
  $ mop list 1

Constraints:
  Execute "mop start" command before using this command.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		opDir := filepath.Join(mopConfig.DataDir, mopConfig.Repository, args[0])
		op, err := command.LoadOperation(filepath.Join(opDir, "operation.json"))
		if err != nil {
			return err
		}

		idx, err := fuzzyfinder.Find(
			op.Commands,
			func(i int) string {
				return op.Commands[i]
			},
			fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
				if i == -1 {
					return ""
				}
				md, err := markdown.Render(op.Commands[i], glamour.WithAutoStyle())
				if err != nil {
					return ""
				}
				return md
			}))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(op.Commands[idx])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
