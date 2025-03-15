package main

import (
	"context"
	"fmt"

	"github.com/umutondersu/dockerfile-sources/internal/ghdocker"
	"github.com/umutondersu/dockerfile-sources/internal/input"
	"github.com/umutondersu/dockerfile-sources/internal/jsonoutput"
)

func main() {
	url := "https://gist.githubusercontent.com/jmelis/c60e61a893248244dc4fa12b946585c4/raw/25d39f67f2405330a6314cad64fac423a171162c/sources.txt" // TODO: turn this into an input

	body, err := input.GetHTTPResponseBody(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	sources := input.ParseRepositorySources(body)

	ctx := context.Background()
	c := ghdocker.NewClient("")

	dockerfiles, err := c.GetDockerFiles(ctx, sources)
	if err != nil {
		fmt.Println(err)
		return
	}

	output, err := jsonoutput.Convert(dockerfiles)
	if err != nil {
		fmt.Println(err)
		return
	}

	jsonStr, err := output.ToJSON()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(jsonStr)
}
