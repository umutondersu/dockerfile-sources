package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/umutondersu/dockerfile-sources/internal/ghdocker"
	"github.com/umutondersu/dockerfile-sources/internal/input"
	"github.com/umutondersu/dockerfile-sources/internal/jsonoutput"
)

func main() {
	timeout := 5 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Validate and sanitize URL
	repoURL := os.Getenv("REPOSITORY_LIST_URL")
	if repoURL == "" {
		fmt.Println("Error: REPOSITORY_LIST_URL environment variable is not set")
		os.Exit(1)
	}
	parsedURL, err := url.Parse(repoURL)
	if err != nil || (!parsedURL.IsAbs() || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https")) {
		fmt.Printf("Error: Invalid URL format: %v\n", err)
		os.Exit(1)
	}

	// Validate GitHub token
	githubAccessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	if githubAccessToken == "" {
		fmt.Println("Warning: GITHUB_ACCESS_TOKEN not set. Rate limiting may occur.")
	}

	body, err := input.GetHTTPResponseBody(repoURL)
	if err != nil {
		fmt.Printf("Error Getting Response Body: %v\n", err)
		return
	}

	sources := input.ParseRepositorySources(body)

	c := ghdocker.NewClient(githubAccessToken)

	dockerfiles, err := c.GetDockerFiles(ctx, sources)
	if err != nil {
		switch e := err.(type) {
		case *ghdocker.ErrRateLimit:
			fmt.Printf("GitHub API rate limit exceeded. Reset time: %v\n", e.ResetTime)
			os.Exit(1)
		case *ghdocker.ErrGitHub:
			fmt.Printf("GitHub API error (Status %d): %v\n", e.StatusCode, e.Message)
			os.Exit(1)
		default:
			fmt.Printf("Error Getting DockerFiles: %v\n", err)
			os.Exit(1)
		}
	}

	jsonStr, err := jsonoutput.GenerateJSONOutput(dockerfiles)
	if err != nil {
		fmt.Printf("Error Parsing to JSON: %v\n", err)
		return
	}

	fmt.Println(jsonStr)
}
