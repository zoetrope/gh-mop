package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zoetrope/gh-mop/core"
)

var uploadOffset int64
var removeAnsiEscape bool

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload ISSUE_NUMBER RESULT_FILE",
	Short: "Post command results to Issue comments",
	Long: `Post the command execution results to the specified Issue's comments.

Arguments:
  ISSUE_NUMBER: Issue number defining an operation on GitHub
  RESULT_FILE: File containing the command execution results

Examples:
  $ mop upload 1 typescript.txt

Constraints:
  Execute "mop start" command before using this command.
  Files larger than 65,000 bytes cannot be uploaded.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := core.NewClient(mopConfig.Owner, mopConfig.Repository)
		if err != nil {
			return err
		}
		issue, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		offset, err := core.UploadResult(client, issue,
			fmt.Sprintf("%s/%s/%d/typescript.txt", mopConfig.DataDir, mopConfig.Repository, issue), uploadOffset, removeAnsiEscape)
		if err != nil {
			return err
		}
		fmt.Println(offset)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	fs := uploadCmd.Flags()
	fs.Int64VarP(&uploadOffset, "offset", "o", 0, "writes content only after the specified byte number")
	fs.BoolVarP(&removeAnsiEscape, "remove-ansi-escape", "r", true, "removes ANSI escape sequences from the result file")
}
