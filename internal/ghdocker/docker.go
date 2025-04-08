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

// processDockerFile processes a single Dockerfile and sends it to the result channel
func (c *Client) processDockerFile(ctx context.Context, src input.Source, path string, resultCh chan<- DockerFile) error {
	// Check if context is already canceled
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	content, err := c.getFileContent(ctx, src, path)
	if err != nil {
		return fmt.Errorf("error getting content for %s in %s/%s: %w", path, src.Owner, src.Repo, err)
	}

	images := extractImages(content)

	select {
	case resultCh <- DockerFile{
		Source: &src,
		Path:   path,
		Images: images,
	}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// processSource processes a single source repository and finds all Dockerfiles
func (c *Client) processSource(ctx context.Context, src input.Source, resultCh chan<- DockerFile) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	fileTree, _, err := c.client.Git.GetTree(ctx, src.Owner, src.Repo, src.CommitSha, true)
	if err != nil {
		return fmt.Errorf("failed to get File Tree from %s/%s:%s: %w", src.Owner, src.Repo, src.CommitSha, err)
	}

	var fileWg sync.WaitGroup
	for _, entry := range fileTree.Entries {
		filePath := entry.GetPath()
		// if entry type is 'blob' it means the entry is a file and not a directory
		if entry.GetType() == "blob" && isDockerFile(filePath) {
			fileWg.Add(1)
			go func(path string, s input.Source) {
				defer fileWg.Done()

				err := c.processDockerFile(ctx, s, path, resultCh)
				if err != nil {
					fmt.Printf("%v\n", err)
				}
			}(filePath, src)
		}
	}

	fileWg.Wait()
	return nil
}

// getChanResult collects results from the channel and returns them as a slice
func getChanResult(ctx context.Context, resultCh <-chan DockerFile) ([]DockerFile, error) {
	var dockerfiles []DockerFile

	for {
		select {
		case dockerfile, ok := <-resultCh:
			if !ok {
				// Sort dockerfiles by path to ensure consistent ordering
				sort.Slice(dockerfiles, func(i, j int) bool {
					return dockerfiles[i].Path < dockerfiles[j].Path
				})
				return dockerfiles, nil
			}
			dockerfiles = append(dockerfiles, dockerfile)
		case <-ctx.Done():
			// Context canceled
			return dockerfiles, ctx.Err()
		}
	}
}

func (c *Client) GetDockerFiles(ctx context.Context, sources []input.Source) ([]DockerFile, error) {
	resultCh := make(chan DockerFile)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var wg sync.WaitGroup

	// Process each source in a separate goroutine
	for _, source := range sources {
		wg.Add(1)
		go func(src input.Source) {
			defer wg.Done()

			err := c.processSource(ctx, src, resultCh)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		}(source)
	}

	// Close the result channel when all sources are processed
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Collect and return results
	return getChanResult(ctx, resultCh)
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

		images = append(images, curImage)
	}

	return images
}

func isDockerFile(path string) bool {
	return strings.HasSuffix(path, "/Dockerfile") || path == "Dockerfile"
}
