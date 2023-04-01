package core

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
)

type Client struct {
	ghClient   api.RESTClient
	owner      string
	repository string
}

func NewClient(owner, repo string) (*Client, error) {
	ghClient, err := gh.RESTClient(nil)
	if err != nil {
		return nil, err
	}
	return &Client{
		ghClient:   ghClient,
		owner:      owner,
		repository: repo,
	}, nil
}

func (c *Client) GetIssueContent(issue int) (string, error) {
	response := &struct {
		Body string `json:"body"`
	}{}
	err := c.ghClient.Get(fmt.Sprintf("repos/%s/%s/issues/%d", c.owner, c.repository, issue), &response)
	if err != nil {
		return "", err
	}
	return response.Body, nil
}

func (c *Client) PostComment(issue int, comment string) (string, error) {
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
	err = c.ghClient.Post(fmt.Sprintf("repos/%s/%s/issues/%d/comments", c.owner, c.repository, issue), r, &response)
	if err != nil {
		return "", err
	}
	return response.URL, nil
}
