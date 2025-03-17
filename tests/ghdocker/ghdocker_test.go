package ghdocker_test

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/umutondersu/dockerfile-sources/internal/ghdocker"
	"github.com/umutondersu/dockerfile-sources/internal/input"
	"github.com/umutondersu/dockerfile-sources/tests/testdata"
)

func TestGetDockerFiles(t *testing.T) {
	sources := testdata.TestSources
	dockerFiles := testdata.TestDockerFiles

	tests := []struct {
		name          string
		sources       []input.Source
		expectedFiles []ghdocker.DockerFile
	}{
		{
			name:          "Single Dockerfile, Single repo",
			sources:       sources[0:1],
			expectedFiles: dockerFiles[0:1],
		},
		{
			name:          "Multiple dockerfiles, Single repo, Duplicate Images",
			sources:       sources[1:2],
			expectedFiles: dockerFiles[1:3],
		},
		{
			name:          "Multiple dockerfiles, Multiple repos, Duplicate Images",
			sources:       sources[0:2],
			expectedFiles: dockerFiles[0:3],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := ghdocker.NewClient(os.Getenv("GITHUB_ACCESS_TOKEN"))

			files, err := client.GetDockerFiles(context.Background(), tt.sources)
			if err != nil {
				t.Errorf("Error encountered: %v", err)
				return
			}
			if !reflect.DeepEqual(files, tt.expectedFiles) {
				t.Errorf("GetDockerFiles() = %v\n want %v", files, tt.expectedFiles)
			}
		})
	}
}
