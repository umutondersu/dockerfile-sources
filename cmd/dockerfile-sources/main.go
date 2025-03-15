package main

import (
	"context"
	"fmt"

	"github.com/umutondersu/dockerfile-sources/internal/ghdocker"
	"github.com/umutondersu/dockerfile-sources/internal/input"
)

func main() {
	url := "https://gist.githubusercontent.com/jmelis/c60e61a893248244dc4fa12b946585c4/raw/25d39f67f2405330a6314cad64fac423a171162c/sources.txt" // TODO: turn this into an input

	body, err := input.GetHTTPResponseBody(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	sources := input.ParseRepositorySources(body)
	source := sources[0]

	ctx := context.Background()
	c := ghdocker.NewClient("")

	dockerfiles, err := c.GetDockerfiles(ctx, source)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(dockerfiles)
}
