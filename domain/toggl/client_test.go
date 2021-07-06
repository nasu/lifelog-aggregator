package toggl

import (
	"os"
	"testing"
)

func TestGetDetails(t *testing.T) {
	tests := []struct {
		since string
		until string
		want  int
	}{
		// inspect with my data in the past directly.
		{"2020-03-01", "2020-03-01", 0},
		{"2020-03-01", "2020-03-31", 28},
		{"2020-04-01", "2020-04-30", 42},
		{"2020-05-01", "2020-05-31", 27},
		{"2020-06-01", "2020-06-30", 24},
		{"2020-03-01", "2020-06-30", 121},
	}
	c := &Client{
		UserAgent:   os.Getenv("TOGGL_USER_AGENT"),
		WorkSpaceID: os.Getenv("TOGGL_WORKSPACE_ID"),
		ApiToken:    os.Getenv("TOGGL_API_TOKEN"),
	}
	for _, tt := range tests {
		data, err := c.GetDetails(tt.since, tt.until)
		if err != nil {
			t.Fatalf("error occured: %s", err)
		}
		if len(data) != tt.want {
			t.Errorf("wrong length. got=%d, want=%d", len(data), tt.want)
		}
	}
}

func TestGetDetailsInvalidEnvs(t *testing.T) {
	tests := []struct {
		title       string
		userAgent   string
		workspaceID string
		apiToken    string
	}{
		{"no user_agent", "", os.Getenv("TOGGL_WORKSPACE_ID"), os.Getenv("TOGGL_API_TOKEN")},
		{"no workspace_id", os.Getenv("TOGGL_USER_AGENT"), "", os.Getenv("TOGGL_API_TOKEN")},
		{"no api_token", os.Getenv("TOGGL_USERAGENT"), os.Getenv("TOGGL_WORKSPACE_ID"), ""},
	}
	for _, tt := range tests {
		c := &Client{
			UserAgent:   tt.userAgent,
			WorkSpaceID: tt.workspaceID,
			ApiToken:    tt.apiToken,
		}
		_, err := c.GetDetails("2021-06-01", "2021-06-01")
		if err == nil {
			t.Error(tt.title)
		}
	}
}
