package main

import (
	"context"
	"fmt"

	"github.com/umutondersu/dockerfile-sources/internal/ghdocker"
	"github.com/umutondersu/dockerfile-sources/internal/input"
	"github.com/umutondersu/dockerfile-sources/internal/jsonoutput"
)

func main() {
	// TODO: turn these into an input
	url := "https://gist.githubusercontent.com/jmelis/c60e61a893248244dc4fa12b946585c4/raw/25d39f67f2405330a6314cad64fac423a171162c/sources.txt"
	access_token := ""

	body, err := input.GetHTTPResponseBody(url)
	if err != nil {
		fmt.Println("Error Getting Response Body: %w", err)
		return
	}

	sources := input.ParseRepositorySources(body)

	c := ghdocker.NewClient(access_token)

	dockerfiles, err := c.GetDockerFiles(context.Background(), sources)
	if err != nil {
		fmt.Println("Error Getting DockerFiles: %w", err)
		return
	}

	jsonStr, err := jsonoutput.GenerateJSONOutput(dockerfiles)
	if err != nil {
		fmt.Println("Error Parsing to JSON: %w", err)
		return
	}

	fmt.Println(jsonStr)
}
