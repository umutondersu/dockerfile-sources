package ghdocker_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/umutondersu/dockerfile-sources/internal/ghdocker"
	"github.com/umutondersu/dockerfile-sources/internal/input"
)

func TestGetDockerFiles(t *testing.T) {
	allSources := []input.Source{
		{
			Owner:     "app-sre",
			Repo:      "qontract-reconcile",
			CommitSha: "30af65af14a2dce962df923446afff24dd8f123e",
		},
		{
			Owner:     "app-sre",
			Repo:      "container-images",
			CommitSha: "c260deaf135fc0efaab365ea234a5b86b3ead404",
		},
	}

	tests := []struct {
		name          string
		sources       []input.Source
		expectedFiles []ghdocker.DockerFile
	}{
		{
			name:    "Single Dockerfile, Single repo",
			sources: allSources[0:1], expectedFiles: []ghdocker.DockerFile{
				{
					Source: &allSources[0],
					Path:   "dockerfiles/Dockerfile",
					Images: []string{"quay.io/app-sre/qontract-reconcile-base:0.2.1"},
				},
			},
		},
		{
			name:    "Multiple dockerfiles, Single repo, Duplicate Images",
			sources: allSources[1:2],
			expectedFiles: []ghdocker.DockerFile{
				{
					Source: &allSources[1],
					Path:   "jiralert/Dockerfile",
					Images: []string{"registry.access.redhat.com/ubi8/go-toolset:latest", "registry.access.redhat.com/ubi8-minimal:8.2"},
				},
				{
					Source: &allSources[1],
					Path:   "qontract-reconcile-base/Dockerfile",
					Images: []string{
						"registry.access.redhat.com/ubi8/ubi:8.2",
						"registry.access.redhat.com/ubi8/ubi:8.2",
						"registry.access.redhat.com/ubi8/ubi:8.2",
					},
				},
			},
		},
		{
			name:    "Multiple dockerfiles, Multiple repos, Duplicate Images",
			sources: allSources[0:2],
			expectedFiles: []ghdocker.DockerFile{
				{
					Source: &allSources[0],
					Path:   "dockerfiles/Dockerfile",
					Images: []string{"quay.io/app-sre/qontract-reconcile-base:0.2.1"},
				},

				{
					Source: &allSources[1],
					Path:   "jiralert/Dockerfile",
					Images: []string{"registry.access.redhat.com/ubi8/go-toolset:latest", "registry.access.redhat.com/ubi8-minimal:8.2"},
				},
				{
					Source: &allSources[1],
					Path:   "qontract-reconcile-base/Dockerfile",
					Images: []string{
						"registry.access.redhat.com/ubi8/ubi:8.2",
						"registry.access.redhat.com/ubi8/ubi:8.2",
						"registry.access.redhat.com/ubi8/ubi:8.2",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := ghdocker.NewClient("")

			files, err := client.GetDockerFiles(context.Background(), tt.sources)
			if err != nil {
				t.Errorf("Error encountered: %v", err)
				return
			}
			if !reflect.DeepEqual(files, tt.expectedFiles) {
				t.Errorf("GetDockerFiles() = %v, want %v", files, tt.expectedFiles)
			}
		})
	}
}
