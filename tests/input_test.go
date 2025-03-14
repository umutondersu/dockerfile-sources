package input_test

import (
	"reflect"
	"testing"

	"github.com/umutondersu/dockerfile-sources/internal/input"
)

func TestParseRepositorySources(t *testing.T) {
	tests := []struct {
		body string
		want []input.Source
	}{
		{
			body: "https://github.com/app-sre/qontract-reconcile.git 30af65af14a2dce962df923446afff24dd8f123e",
			want: []input.Source{
				{Repository: "https://github.com/app-sre/qontract-reconcile.git", CommitSha: "30af65af14a2dce962df923446afff24dd8f123e"},
			},
		},
		{
			body: "https://github.com/app-sre/qontract-reconcile.git 30af65af14a2dce962df923446afff24dd8f123e\nhttps://github.com/app-sre/container-images.git c260deaf135fc0efaab365ea234a5b86b3ead404",
			want: []input.Source{
				{Repository: "https://github.com/app-sre/qontract-reconcile.git", CommitSha: "30af65af14a2dce962df923446afff24dd8f123e"},
				{Repository: "https://github.com/app-sre/container-images.git", CommitSha: "c260deaf135fc0efaab365ea234a5b86b3ead404"},
			},
		},
		{
			body: "",
			want: []input.Source{},
		},
		{
			// typo in the link
			body: "https://gitub.com/app-sre/qontract-reconcile.git 30af65af14a2dce962df923446afff24dd8f123e\nhttps://github.com/app-sre/container-images.git c260deaf135fc0efaab365ea234a5b86b3ead404",
			want: []input.Source{
				{Repository: "https://github.com/app-sre/container-images.git", CommitSha: "c260deaf135fc0efaab365ea234a5b86b3ead404"},
			},
		},
		{
			// Wrong CommitSha format
			body: "https://github.com/app-sre/qontract-reconcile.git 30af65af14a2dce962df923446afff24dd8f123\nhttps://github.com/app-sre/container-images.git c260deaf135fc0efaab365ea234a5b86b3ead404",
			want: []input.Source{
				{Repository: "https://github.com/app-sre/container-images.git", CommitSha: "c260deaf135fc0efaab365ea234a5b86b3ead404"},
			},
		},
	}
	for _, tt := range tests {
		got := input.ParseRepositorySources(tt.body)
		if !reflect.DeepEqual(got, tt.want) && !(len(got) == 0 && len(tt.want) == 0) {
			t.Errorf("ParseRepositorySources() = %v, want %v", got, tt.want)
		}
	}
}
