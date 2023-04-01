package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zoetrope/gh-mop/core"
)

var uploadOffset int64

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
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := core.NewClient(mopConfig.Owner, mopConfig.Repository)
		if err != nil {
			return err
		}
		issue, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		url, offset, err := core.UploadResult(client, issue, args[1], uploadOffset)
		if err != nil {
			return err
		}
		fmt.Println(url)
		fmt.Println(offset)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	fs := uploadCmd.Flags()
	fs.Int64VarP(&uploadOffset, "offset", "o", 0, "writes content only after the specified byte number")
}
