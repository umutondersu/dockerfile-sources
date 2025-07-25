package input_test

import (
	"reflect"
	"testing"

	"github.com/umutondersu/dockerfile-sources/internal/input"
)

func TestParseRepositorySources(t *testing.T) {
	tests := []struct {
		name string
		body string
		want []input.Source
	}{
		{
			name: "Basic",
			body: "https://github.com/app-sre/qontract-reconcile.git 30af65af14a2dce962df923446afff24dd8f123e",
			want: []input.Source{
				{Owner: "app-sre", Repo: "qontract-reconcile", CommitSha: "30af65af14a2dce962df923446afff24dd8f123e"},
			},
		},
		{
			name: "2 Sources",
			body: "https://github.com/app-sre/qontract-reconcile.git 30af65af14a2dce962df923446afff24dd8f123e\nhttps://github.com/app-sre/container-images.git c260deaf135fc0efaab365ea234a5b86b3ead404",
			want: []input.Source{
				{Owner: "app-sre", Repo: "qontract-reconcile", CommitSha: "30af65af14a2dce962df923446afff24dd8f123e"},
				{Owner: "app-sre", Repo: "container-images", CommitSha: "c260deaf135fc0efaab365ea234a5b86b3ead404"},
			},
		},
		{
			name: "Empty body",
			body: "",
			want: nil,
		},
		{
			name: "Typo in the link",
			body: "https://gitub.com/app-sre/qontract-reconcile.git 30af65af14a2dce962df923446afff24dd8f123e\nhttps://github.com/app-sre/container-images.git c260deaf135fc0efaab365ea234a5b86b3ead404",
			want: []input.Source{
				{Owner: "app-sre", Repo: "container-images", CommitSha: "c260deaf135fc0efaab365ea234a5b86b3ead404"},
			},
		},
		{
			name: "Wrong Sha Format",
			body: "https://github.com/app-sre/qontract-reconcile.git 30af65af14a2dce962df923446afff24dd8f123\nhttps://github.com/app-sre/container-images.git c260deaf135fc0efaab365ea234a5b86b3ead404",
			want: []input.Source{
				{Owner: "app-sre", Repo: "container-images", CommitSha: "c260deaf135fc0efaab365ea234a5b86b3ead404"},
			},
		},
		{
			name: "Url without .git suffix",
			body: "https://github.com/kubernetes/kubernetes 8f0b92c0512afb25c8b2667ddfd1c7d5409903d3",
			want: nil,
		},
		{
			name: "Multiple Invalid entries",
			body: "https://gitlab.com/user/repo.git abc123\nhttps://bitbucket.org/user/repo def456",
			want: nil,
		},
		{
			name: "Extra Spaces",
			body: "  https://github.com/test/repo.git   abcdef1234567890abcdef1234567890abcdef12",
			want: nil,
		},
		{
			name: "Mixed valid and invalid formats",
			body: "https://github.com/valid/repo.git 1234567890123456789012345678901234567890\nmalformed/input\nhttps://github.com/another/valid.git 2234567890123456789012345678901234567890",
			want: []input.Source{
				{Owner: "valid", Repo: "repo", CommitSha: "1234567890123456789012345678901234567890"},
				{Owner: "another", Repo: "valid", CommitSha: "2234567890123456789012345678901234567890"},
			},
		},
		{
			name: "Invalid commit SHA lengths",
			body: "https://github.com/test/short.git abc123\nhttps://github.com/test/long.git 1234567890123456789012345678901234567890extra",
			want: nil,
		},
		{
			name: "Case sensitivity in URL",
			body: "HTTPS://GITHUB.COM/org/repo.git abcdef1234567890abcdef1234567890abcdef12",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := input.ParseRepositorySources(tt.body)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRepositorySources() = %v, want %v", got, tt.want)
			}
		})
	}
}
