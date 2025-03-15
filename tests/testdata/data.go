package testdata

import (
	"github.com/umutondersu/dockerfile-sources/internal/ghdocker"
	"github.com/umutondersu/dockerfile-sources/internal/input"
)

var TestSources = []input.Source{
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

var TestDockerFiles = []ghdocker.DockerFile{
	{
		Source: &TestSources[0],
		Path:   "dockerfiles/Dockerfile",
		Images: []string{"quay.io/app-sre/qontract-reconcile-base:0.2.1"},
	},
	{
		Source: &TestSources[1],
		Path:   "jiralert/Dockerfile",
		Images: []string{"registry.access.redhat.com/ubi8/go-toolset:latest", "registry.access.redhat.com/ubi8-minimal:8.2"},
	},
	{
		Source: &TestSources[1],
		Path:   "qontract-reconcile-base/Dockerfile",
		Images: []string{
			"registry.access.redhat.com/ubi8/ubi:8.2",
			"registry.access.redhat.com/ubi8/ubi:8.2",
			"registry.access.redhat.com/ubi8/ubi:8.2",
		},
	},
}
