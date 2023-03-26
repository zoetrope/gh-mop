package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/cli/go-gh"
)

func UploadResult(owner, repo string, issue int, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	if fileInfo.Size() > 65000 {
		return fmt.Errorf("file is too large")
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	cleanedContent := removeANSIEscapeSequences(string(content))

	client, err := gh.RESTClient(nil)
	if err != nil {
		return err
	}

	comment := "```\n" + cleanedContent + "\n```"
	var body = struct {
		Body string `json:"body"`
	}{
		Body: comment,
	}
	var response = struct {
		URL string `json:"html_url"`
	}{}
	j, err := json.Marshal(body)
	r := bytes.NewReader(j)
	err = client.Post(fmt.Sprintf("repos/%s/%s/issues/%d/comments", owner, repo, issue), r, &response)
	if err != nil {
		return err
	}

	return nil
}
func removeANSIEscapeSequences(input string) string {
	ansiEscapeRegex := regexp.MustCompile(`\x1B\[[0-?]*[ -/]*[@-~]`)
	return ansiEscapeRegex.ReplaceAllString(input, "")
}
