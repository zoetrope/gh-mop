package cmd

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zoetrope/gh-mop/parser"
)

// nextCmd represents the next command
var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		opDir := filepath.Join(mopConfig.DataDir, mopConfig.Repository, args[0])
		op, err := parser.LoadOperation(filepath.Join(opDir, "operation.json"))
		if err != nil {
			return err
		}

		index, err := strconv.Atoi(args[1])
		fmt.Print(op.Commands[index])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(nextCmd)
}
