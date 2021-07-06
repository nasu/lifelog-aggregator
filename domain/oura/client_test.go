package oura

import (
	"os"
	"testing"
	"time"
)

func TestGetDetails(t *testing.T) {
	tests := []struct {
		title string
	}{
		{"user info"},
	}
	c := &Client{
		AccessToken: os.Getenv("OURA_ACCESS_TOKEN"),
	}
	for _, tt := range tests {
		_, err := c.UserInfo()
		if err != nil {
			t.Fatalf("error occured: title=%s, err=%s", tt.title, err)
		}
	}
}

func TestReadinessOneDay(t *testing.T) {
	tests := []struct {
		title string
		date  string
	}{
		{"today", time.Now().Format("2006-01-02")},
		{"yesterday", time.Now().Add(time.Hour * 24 * -1).Format("2006-01-02")},
	}
	c := &Client{
		AccessToken: os.Getenv("OURA_ACCESS_TOKEN"),
	}
	for _, tt := range tests {
		_, err := c.ReadinessOneDay(tt.date)
		if err != nil {
			t.Fatalf("error occured: title=%s, error=%s", tt.title, err)
		}
	}
}

func TestSleepOneDay(t *testing.T) {
	tests := []struct {
		title string
		date  string
	}{
		{"today", time.Now().Format("2006-01-02")},
		{"yesterday", time.Now().Add(time.Hour * 24 * -1).Format("2006-01-02")},
	}
	c := &Client{
		AccessToken: os.Getenv("OURA_ACCESS_TOKEN"),
	}
	for _, tt := range tests {
		_, err := c.SleepOneDay(tt.date)
		if err != nil {
			t.Fatalf("error occured: title=%s, error=%s", tt.title, err)
		}
	}
}

func TestGetDetailsInvalidEnvs(t *testing.T) {
	tests := []struct {
		title       string
		accessToken string
	}{
		{"no access token", ""},
	}
	for _, tt := range tests {
		c := &Client{
			AccessToken: "",
		}
		_, err := c.UserInfo()
		if err == nil {
			t.Error(tt.title)
		}
	}
}
