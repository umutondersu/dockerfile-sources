package input_test

import (
	"testing"

	"github.com/umutondersu/dockerfile-sources/internal/input"
)

func TestParseRepositorySources(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		url  string
		want []input.Source
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := input.ParseRepositorySources(tt.url)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("ParseRepositorySources() = %v, want %v", got, tt.want)
			}
		})
	}
}
