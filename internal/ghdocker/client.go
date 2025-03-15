package ghdocker

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/google/go-github/v69/github"
	"github.com/umutondersu/dockerfile-sources/internal/input"
	"golang.org/x/oauth2"
)

type Client struct {
	client *github.Client
}

// NewClient Generates a new Github Client with the given OAUTH2 access token
// If the token is an empty a new http.Client will be used
func NewClient(token string) *Client {
	var httpClient *github.Client

	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		httpClient = github.NewClient(tc)
	} else {
		httpClient = github.NewClient(nil)
	}

	return &Client{
		client: httpClient,
	}
}

func (c *Client) getFileContent(ctx context.Context, source input.Source, path string) (string, error) {
	content, _, _, err := c.client.Repositories.GetContents(ctx, source.Owner, source.Repo, path, &github.RepositoryContentGetOptions{
		Ref: source.CommitSha,
	})
	if err != nil {
		return "", err
	}

	decoded, err := base64.StdEncoding.DecodeString(*content.Content)
	if err != nil {
		return "", fmt.Errorf("failed to decode content: %w", err)
	}

	return string(decoded), nil
}
