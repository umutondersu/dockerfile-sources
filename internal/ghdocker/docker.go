package ghdocker

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/umutondersu/dockerfile-sources/internal/input"
)

// DockerFile represents a Dockerfile found in a repository
// Multiple Dockerfiles can belong to a single Source
type DockerFile struct {
	Source *input.Source
	Path   string
	Images []string
}

func (c *Client) GetDockerFiles(ctx context.Context, sources []input.Source) ([]DockerFile, error) {
	var dockerfiles []DockerFile

	for _, source := range sources {
		fileTree, err := c.getFileTree(ctx, source)
		if err != nil {
			return nil, fmt.Errorf("Failed the get File Tree from %s/%s:%s: %w", source.Owner, source.Repo, source.CommitSha, err)
		}

		ch := make(chan string)
		var wg sync.WaitGroup

		for _, entry := range fileTree.Entries {
			filePath := entry.GetPath()
			// if entry type is 'blob' it means the entry is a file and not a directory
			if entry.GetType() == "blob" && isDockerFile(filePath) {
				wg.Add(1)
				go func(path string) {
					_, err := c.getFileContent(ctx, source, path, ch, &wg)
					if err != nil {
						fmt.Printf("Error getting content for %s in %s/%s: %v\n", path, source.Owner, source.Repo, err)
					}
				}(filePath)
			}
		}

		go func() {
			wg.Wait()
			close(ch)
		}()

		for content := range ch {
			Path, Images := extractData(content)

			dockerfiles = append(dockerfiles, DockerFile{
				Source: &source,
				Path:   Path,
				Images: Images,
			})
		}
	}

	// Sort dockerfiles by path to ensure consistent ordering
	sort.Slice(dockerfiles, func(i, j int) bool {
		return dockerfiles[i].Path < dockerfiles[j].Path
	})
	return dockerfiles, nil
}

func extractData(content string) (string, []string) {
	var images []string
	var path string
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		if i == 0 {
			path = line
			continue
		}

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

		images = append(images, curImage)
	}

	return path, images
}

func isDockerFile(path string) bool {
	return strings.HasSuffix(path, "/Dockerfile") || path == "Dockerfile"
}
