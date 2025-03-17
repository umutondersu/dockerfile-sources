package main

import (
	"context"
	"fmt"
	"os"

	"github.com/umutondersu/dockerfile-sources/internal/ghdocker"
	"github.com/umutondersu/dockerfile-sources/internal/input"
	"github.com/umutondersu/dockerfile-sources/internal/jsonoutput"
)

func main() {
	url := os.Getenv("REPOSITORY_LIST_URL")
	if url == "" {
		fmt.Println("Error: REPOSITORY_LIST_URL environment variable is not set")
		return
	}
	githubAccessToken := os.Getenv("GITHUB_ACCESS_TOKEN")

	body, err := input.GetHTTPResponseBody(url)
	if err != nil {
		fmt.Printf("Error Getting Response Body: %v\n", err)
		return
	}

	sources := input.ParseRepositorySources(body)

	c := ghdocker.NewClient(githubAccessToken)

	dockerfiles, err := c.GetDockerFiles(context.Background(), sources)
	if err != nil {
		fmt.Printf("Error Getting DockerFiles: %v\n", err)
		return
	}

	jsonStr, err := jsonoutput.GenerateJSONOutput(dockerfiles)
	if err != nil {
		fmt.Printf("Error Parsing to JSON: %v\n", err)
		return
	}

	fmt.Println(jsonStr)
}
