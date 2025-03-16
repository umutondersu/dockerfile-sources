package jsonoutput_test

import (
	"reflect"
	"testing"

	"github.com/umutondersu/dockerfile-sources/internal/ghdocker"
	"github.com/umutondersu/dockerfile-sources/internal/jsonoutput"
	"github.com/umutondersu/dockerfile-sources/tests/testdata"
)

func TestConvert(t *testing.T) {
	dockerFiles := testdata.TestDockerFiles

	tests := []struct {
		name     string
		input    []ghdocker.DockerFile
		expected map[string]map[string][]string
	}{
		{
			name:  "Single Dockerfile, Single Repo",
			input: dockerFiles[0:1],
			expected: map[string]map[string][]string{
				"https://github.com/app-sre/qontract-reconcile.git:30af65af14a2dce962df923446afff24dd8f123e": {
					"dockerfiles/Dockerfile": {"quay.io/app-sre/qontract-reconcile-base:0.2.1"},
				},
			},
		},
		{
			name:  "Multiple Dockerfiles, Single Repo",
			input: dockerFiles[1:3],
			expected: map[string]map[string][]string{
				"https://github.com/app-sre/container-images.git:c260deaf135fc0efaab365ea234a5b86b3ead404": {
					"jiralert/Dockerfile": {
						"registry.access.redhat.com/ubi8/go-toolset:latest",
						"registry.access.redhat.com/ubi8-minimal:8.2",
					},
					"qontract-reconcile-base/Dockerfile": {
						"registry.access.redhat.com/ubi8/ubi:8.2",
						"registry.access.redhat.com/ubi8/ubi:8.2",
						"registry.access.redhat.com/ubi8/ubi:8.2",
					},
				},
			},
		},
		{
			name:  "Multiple Dockerfiles, Multiple Repos",
			input: dockerFiles[0:3],
			expected: map[string]map[string][]string{
				"https://github.com/app-sre/qontract-reconcile.git:30af65af14a2dce962df923446afff24dd8f123e": {
					"dockerfiles/Dockerfile": {"quay.io/app-sre/qontract-reconcile-base:0.2.1"},
				},
				"https://github.com/app-sre/container-images.git:c260deaf135fc0efaab365ea234a5b86b3ead404": {
					"jiralert/Dockerfile": {
						"registry.access.redhat.com/ubi8/go-toolset:latest",
						"registry.access.redhat.com/ubi8-minimal:8.2",
					},
					"qontract-reconcile-base/Dockerfile": {
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
			output, err := jsonoutput.Convert(tt.input)
			if err != nil {
				t.Errorf("Convert() error = %v", err)
				return
			}

			if !reflect.DeepEqual(output.Data, tt.expected) {
				t.Errorf("Convert() = %v, want %v", output.Data, tt.expected)
			}
		})
	}
}
