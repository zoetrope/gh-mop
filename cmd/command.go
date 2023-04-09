package cmd

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zoetrope/gh-mop/pkg/command"
)

// commandCmd represents the command command
var commandCmd = &cobra.Command{
	Use:   "command ISSUE_NUMBER COMMAND_INDEX",
	Short: "Retrieve command string from Issue",
	Long: `Retrieve command string from Issue.

Arguments:
  ISSUE_NUMBER: Issue number defining an operation
  COMMAND_INDEX: Index of the desired command

Examples:
  $ mop command 1 0

Constraints:
  Execute "mop start" command before using this command.
  Errors occur if the command index exceeds the number of commands defined in the operation.
`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		opDir := filepath.Join(mopConfig.DataDir, mopConfig.Repository, args[0])
		op, err := command.LoadOperation(filepath.Join(opDir, "operation.json"))
		if err != nil {
			return err
		}

		index, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}
		if index >= len(op.Commands) {
			return fmt.Errorf("command index is too large")
		}
		fmt.Print(op.Commands[index])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(commandCmd)
}
