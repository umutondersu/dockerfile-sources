package input

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

var RegexPattern = regexp.MustCompile(`^https://github\.com/([a-zA-Z0-9_.-]+)/([a-zA-Z0-9_.-]+)\.git\s+([0-9a-f]{40})$`)

type Source struct {
	Owner     string
	Repo      string
	CommitSha string
}

func GetHTTPResponseBody(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Error fetching the URL: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("Eror reading the response body: %w", err)
	}

	return string(body), nil
}

func ParseRepositorySources(body string) []Source {
	var sources []Source
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		if matches := RegexPattern.FindStringSubmatch(line); len(matches) == 4 {
			sources = append(sources, Source{
				Owner:     matches[1],
				Repo:      matches[2],
				CommitSha: matches[3],
			})
		}
	}

	return sources
}
