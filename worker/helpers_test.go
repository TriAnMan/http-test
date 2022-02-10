package worker

import "testing"

func TestPrepareUrl(t *testing.T) {
	tests := []struct {
		name   string
		rawUrl string
		want   string
	}{
		{
			"empty string",
			"",
			"http://",
		},
		{
			"no schema",
			"addr",
			"http://addr",
		},
		{
			"regular http",
			"http://addr",
			"http://addr",
		},
		{
			"regular https",
			"https://addr",
			"https://addr",
		},
		{
			"schema in the middle",
			"some_http://addr",
			"http://some_http://addr",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PrepareUrl(tt.rawUrl); got != tt.want {
				t.Errorf("PrepareUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
