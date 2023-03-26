package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zoetrope/gh-mop/core"
)

// operationCmd represents the operation command
var operationCmd = &cobra.Command{
	Use:   "operation ISSUE_NUMBER",
	Short: "Fetch and save operation data from Issue",
	Long: `Fetches operation information from the specified Issue and saves it to a local directory.

Arguments:
  ISSUE_NUMBER: Issue number defining an operation on GitHub

Examples:
  $ mop operation 1

Constraints: 
  Must be executed before running "mop command" or "mop upload" commands.
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		issue, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		op, err := core.GetOperation(mopConfig.Owner, mopConfig.Repository, issue)
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

func init() {
	rootCmd.AddCommand(operationCmd)
}
