package input

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

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
	pattern := `^https://github\.com/([a-zA-Z0-9_.-]+)/([a-zA-Z0-9_.-]+)\.git\s+([0-9a-f]{40})$`
	re := regexp.MustCompile(pattern)

	var sources []Source

	// Scan the body line by line
	scanner := bufio.NewScanner(strings.NewReader(body))
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) == 4 {
			s := Source{
				Owner:     matches[1],
				Repo:      matches[2],
				CommitSha: matches[3],
			}
			sources = append(sources, s)
		}
	}

	return sources
}
