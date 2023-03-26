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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		utilities, err := core.GetUtilities(mopConfig.Owner, mopConfig.Repository, mopConfig.Utilities)
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
