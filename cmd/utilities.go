package cmd

import (
	"fmt"
	"log"

	"github.com/charmbracelet/glamour"
	"github.com/cli/go-gh/pkg/markdown"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
	"github.com/zoetrope/gh-mop/core"
)

// utilitiesCmd represents the utilities command
var utilitiesCmd = &cobra.Command{
	Use:   "utilities",
	Short: "Search and return useful commands from Issues",
	Long: `Search and return useful commands defined in the utilities GitHub Issues using fuzzy search like fzf.
Detailed explanations can be displayed during the search process.

Arguments:
  None

Examples:
  $ mop utilities

Constraints:
  Prepare the utilities GitHub Issues in advance.
  Define the Issue number for the utilities in the configuration file.
  No need to execute "mop start" before using this command.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := core.NewClient(mopConfig.Owner, mopConfig.Repository)
		if err != nil {
			return err
		}
		utilities, err := core.GetUtilities(client, mopConfig.Utilities)
		if err != nil {
			return err
		}

		idx, err := fuzzyfinder.Find(
			utilities,
			func(i int) string {
				return utilities[i].Title
			},
			fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
				if i == -1 {
					return ""
				}
				md, err := markdown.Render(utilities[i].Content, glamour.WithAutoStyle())
				if err != nil {
					return ""
				}
				return md
			}))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(utilities[idx].Commands[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(utilitiesCmd)
}
