package ghdocker

import (
	"context"
	"encoding/base64"
	"fmt"
	"slices"
	"strings"

	"github.com/google/go-github/v69/github"
	"github.com/umutondersu/dockerfile-sources/internal/input"
	"golang.org/x/oauth2"
)

type Client struct {
	client *github.Client
}

type DockerFile struct {
	Source *input.Source
	Path   string
	Images []string
}

// NewClient Generates a new Github Client with the given OAUTH2 access token
// If the token is an empty a new http.Client will be used
func NewClient(token string) *Client {
	var httpClient *github.Client

	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		httpClient = github.NewClient(tc)
	} else {
		httpClient = github.NewClient(nil)
	}

	return &Client{
		client: httpClient,
	}
}

func (c *Client) GetDockerfiles(ctx context.Context, sources ...input.Source) ([]DockerFile, error) {
	var dockerfiles []DockerFile

	for _, source := range sources {
		tree, _, err := c.client.Git.GetTree(ctx, source.Owner, source.Repo, source.CommitSha, true)
		if err != nil {
			return nil, fmt.Errorf("failed to get repository tree for %s/%s: %w", source.Owner, source.Repo, err)
		}

		for _, entry := range tree.Entries {
			filePath := entry.GetPath()
			// if entry type is 'blob' it means the entry is a file and not a directory
			if entry.GetType() == "blob" && isDockerfile(filePath) {
				content, err := c.getFileContent(ctx, source, filePath)
				if err != nil {
					return nil, fmt.Errorf("failed to get content for %s in %s/%s: %w", filePath, source.Owner, source.Repo, err)
				}

				Images := extractImages(content)

				dockerfiles = append(dockerfiles, DockerFile{
					Source: &source,
					Path:   filePath,
					Images: Images,
				})
			}
		}
	}

	return dockerfiles, nil
}

func (c *Client) getFileContent(ctx context.Context, source input.Source, path string) (string, error) {
	content, _, _, err := c.client.Repositories.GetContents(ctx, source.Owner, source.Repo, path, &github.RepositoryContentGetOptions{
		Ref: source.CommitSha,
	})
	if err != nil {
		return "", err
	}

	decoded, err := base64.StdEncoding.DecodeString(*content.Content)
	if err != nil {
		return "", fmt.Errorf("failed to decode content: %w", err)
	}

	return string(decoded), nil
}

func extractImages(content string) []string {
	var images []string
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if !strings.HasPrefix(strings.ToUpper(line), "FROM") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		curImage := parts[1]

		// Special case in Dockerfile
		if strings.ToLower(curImage) == "scratch" {
			continue
		}

		if !slices.Contains(images, curImage) {
			images = append(images, curImage)
		}
	}

	return images
}

func isDockerfile(path string) bool {
	return strings.HasSuffix(path, "/Dockerfile") || path == "Dockerfile"
}
