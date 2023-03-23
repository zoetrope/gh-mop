package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zoetrope/gh-mop/parser"
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
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		issue, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		op, err := parser.GetOperation(mopConfig.Owner, mopConfig.Repository, issue)
		if err != nil {
			return err
		}
		out, err := json.Marshal(op)
		if err != nil {
			return err
		}
		opDir := filepath.Join(mopConfig.DataDir, mopConfig.Repository, args[0])
		err = os.MkdirAll(opDir, 0755)
		if err != nil {
			return err
		}
		err = os.WriteFile(filepath.Join(opDir, "operation.json"), out, 0644)
		return err
	},
}

// parseMarkdown parses the markdown string and returns a list of tasks

func init() {
	rootCmd.AddCommand(startCmd)
}
