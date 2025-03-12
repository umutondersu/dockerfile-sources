package main

import (
	"fmt"

	"github.com/umutondersu/dockerfile-sources/internal/input"
)

func main() {
	url := "https://gist.githubusercontent.com/jmelis/c60e61a893248244dc4fa12b946585c4/raw/25d39f67f2405330a6314cad64fac423a171162c/sources.txt" // TODO: turn this into an input

	result := input.ParseRepositorySources(url)

	fmt.Println(result)
}
