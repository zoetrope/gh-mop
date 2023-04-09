package command

import (
	"fmt"
	"os"

	"github.com/zoetrope/gh-mop/pkg/ansi"
	"github.com/zoetrope/gh-mop/pkg/github"
)

// UploadResult uploads the content as a formatted comment to the given issue.
// The content will be read from the specified file with the given offset.
// If removeEscSequences is true, it will remove ANSI escape sequences from the content.
// Returns the sum of the offset and the length of the content read
func UploadResult(client *github.Client, issue int, filepath string, offset int64) (string, int64, error) {
	if err := checkFileSize(filepath); err != nil {
		return "", 0, err
	}

	content, err := ansi.ProcessFile(filepath)
	if err != nil {
		return "", 0, err
	}

	length := int64(len(content))
	if offset > length {
		return "", 0, fmt.Errorf("offset is too large")
	}
	content = content[offset:]
	if len(content) == 0 {
		return "", 0, fmt.Errorf("no content to upload")
	}

	comment := formatAsCodeBlock(content)
	url, err := client.PostComment(issue, comment)
	if err != nil {
		return "", 0, err
	}
	return url, length, nil
}

func checkFileSize(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	// The max size of a Issue comment is 65536 characters
	if fileInfo.Size() > 65000 {
		return fmt.Errorf("file is too large")
	}

	return nil
}

func formatAsCodeBlock(content string) string {
	return "```\n" + content + "```\n"
}
