package version

import (
	"testing"
)

func Test_get(t *testing.T) {
	tests := []struct {
		name string
		want Info
	}{
		// TODO: Add test cases.
		{
			name: "test-name",
			want: Info{
				Version:      "",
				GitTag:       "",
				GitCommit:    "",
				GitTreeState: "",
				BuildDate:    "",
				GoVersion:    "",
				Compiler:     "",
				Platform:     "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := get()
			t.Log(got)
		})
	}
}
