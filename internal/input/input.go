package input

import (
	"bufio"
	"io"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Source struct {
	Repository string
	CommitSha  string
}

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func GetHTTPResponseBody(url string) (string, error) {
	// Send a GET request
	response, err := http.Get(url)
	if err != nil {
		logger.Error("Error fetching the URL", slog.String("error", err.Error()))
		return "", err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error("Eror reading the response body", slog.String("error", err.Error()))
		return "", err
	}

	return string(body), nil
}

func ParseRepositorySources(body string) []Source {
	// Define the regex pattern
	pattern := `^https://github\.com/[a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+\.git\s+[0-9a-f]{40}$`
	re := regexp.MustCompile(pattern)

	// Create a slice to hold the source structs
	var sources []Source

	// Scan the body line by line
	scanner := bufio.NewScanner(strings.NewReader(body))
	for scanner.Scan() {
		line := scanner.Text()
		if re.MatchString(line) {
			// Split the line into repository and commit SHA
			parts := strings.Fields(line)
			if len(parts) == 2 {
				s := Source{
					Repository: parts[0],
					CommitSha:  parts[1],
				}
				sources = append(sources, s)
			}
		}
	}

	return sources
}
