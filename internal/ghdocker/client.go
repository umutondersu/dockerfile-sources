package ghdocker

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/go-github/v69/github"
	"github.com/umutondersu/dockerfile-sources/internal/input"
)

// Custom error types for better error handling
type ErrGitHub struct {
	StatusCode int
	Message    string
	Err        error
}

func (e *ErrGitHub) Error() string {
	return fmt.Sprintf("github API error (status %d): %s: %v", e.StatusCode, e.Message, e.Err)
}

type ErrRateLimit struct {
	ResetTime time.Time
}

func (e *ErrRateLimit) Error() string {
	return fmt.Sprintf("github API rate limit exceeded, resets at %v", e.ResetTime)
}

type Client struct {
	client *github.Client
}

// NewClient Generates a new Github Client with the given OAUTH2 access token
// If the token is an empty a new http.Client will be used
func NewClient(token string) *Client {
	var httpClient *github.Client

	if token != "" {
		httpClient = github.NewClientWithEnvProxy().WithAuthToken(token)
	} else {
		httpClient = github.NewClient(nil)
	}

	return &Client{
		client: httpClient,
	}
}

func isRetryableError(resp *github.Response) bool {
	if resp == nil {
		// Network errors should be retried
		return true
	}

	// Don't retry client errors except for 429 (Too Many Requests)
	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return resp.StatusCode == http.StatusTooManyRequests
	}

	// Retry server errors
	return resp.StatusCode >= 500
}

// withBackoff executes an operation with exponential backoff
func (c *Client) withBackoff(operation func() (*github.Response, error)) error {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 30 * time.Second
	b.InitialInterval = 100 * time.Millisecond

	return backoff.Retry(func() error {
		resp, err := operation()
		if err == nil {
			return nil
		}

		if !isRetryableError(resp) {
			return backoff.Permanent(err)
		}

		return err
	}, b)
}

func (c *Client) getFileContent(ctx context.Context, source input.Source, path string) (string, error) {
	var content *github.RepositoryContent
	var resp *github.Response
	var err error

	operation := func() (*github.Response, error) {
		content, _, resp, err = c.client.Repositories.GetContents(
			ctx,
			source.Owner,
			source.Repo,
			path,
			&github.RepositoryContentGetOptions{
				Ref: source.CommitSha,
			},
		)
		return resp, err
	}

	if err := c.withBackoff(operation); err != nil {
		notFoundMsg := fmt.Sprintf("file not found: %s/%s:%s Path:%s", source.Owner, source.Repo, source.CommitSha, path)
		defaultErrMsg := "failed to get content"
		if err := handleGitHubResponseError(resp, err, notFoundMsg, defaultErrMsg); err != nil {
			return "", err
		}
	}

	decoded, err := base64.StdEncoding.DecodeString(*content.Content)
	if err != nil {
		return "", fmt.Errorf("failed to decode content: %w", err)
	}

	return string(decoded), nil
}

// handleGitHubResponseError processes GitHub API responses and standardizes error handling
func handleGitHubResponseError(resp *github.Response, err error, notFoundMsg string, defaultErrMsg string) error {
	if err == nil {
		return nil
	}

	if resp != nil {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return fmt.Errorf("%s", notFoundMsg)
		case http.StatusForbidden:
			if resp.Rate.Remaining == 0 {
				return &ErrRateLimit{
					ResetTime: resp.Rate.Reset.Time,
				}
			}
			return &ErrGitHub{
				StatusCode: resp.StatusCode,
				Message:    "access forbidden",
				Err:        err,
			}
		default:
			return &ErrGitHub{
				StatusCode: resp.StatusCode,
				Message:    defaultErrMsg,
				Err:        err,
			}
		}
	}
	return fmt.Errorf("%s: %w", defaultErrMsg, err)
}
