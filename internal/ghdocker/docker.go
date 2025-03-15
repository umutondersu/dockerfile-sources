package ghdocker

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/umutondersu/dockerfile-sources/internal/input"
)

type DockerFile struct {
	Source *input.Source
	Path   string
	Images []string
}

func (c *Client) GetDockerFiles(ctx context.Context, sources []input.Source) ([]DockerFile, error) {
	var dockerfiles []DockerFile

	for _, source := range sources {
		tree, _, err := c.client.Git.GetTree(ctx, source.Owner, source.Repo, source.CommitSha, true)
		if err != nil {
			return nil, fmt.Errorf("failed to get repository tree for %s/%s: %w", source.Owner, source.Repo, err)
		}

		for _, entry := range tree.Entries {
			filePath := entry.GetPath()
			// if entry type is 'blob' it means the entry is a file and not a directory
			if entry.GetType() == "blob" && isDockerFile(filePath) {
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

func isDockerFile(path string) bool {
	return strings.HasSuffix(path, "/Dockerfile") || path == "Dockerfile"
}
