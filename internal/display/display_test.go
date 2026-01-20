package display

import (
	"testing"
	"time"
)

func TestRelativeTime(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{"just now", 0, "just now"},
		{"30 seconds", 30 * time.Second, "just now"},
		{"59 seconds", 59 * time.Second, "just now"},
		{"1 minute", 1 * time.Minute, "1 minute ago"},
		{"2 minutes", 2 * time.Minute, "2 minutes ago"},
		{"59 minutes", 59 * time.Minute, "59 minutes ago"},
		{"1 hour", 1 * time.Hour, "1 hour ago"},
		{"2 hours", 2 * time.Hour, "2 hours ago"},
		{"23 hours", 23 * time.Hour, "23 hours ago"},
		{"1 day", 24 * time.Hour, "1 day ago"},
		{"2 days", 48 * time.Hour, "2 days ago"},
		{"7 days", 7 * 24 * time.Hour, "7 days ago"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timestamp := time.Now().Add(-tt.duration)
			got := RelativeTime(timestamp)
			if got != tt.want {
				t.Errorf("RelativeTime() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRelativeTimeFuture(t *testing.T) {
	future := time.Now().Add(1 * time.Hour)
	got := RelativeTime(future)
	if got != "in the future" {
		t.Errorf("RelativeTime(future) = %q, want %q", got, "in the future")
	}
}

func TestRelativeTimeEdgeCases(t *testing.T) {
	// Exactly at boundary between minute and hour
	t.Run("60 minutes", func(t *testing.T) {
		timestamp := time.Now().Add(-60 * time.Minute)
		got := RelativeTime(timestamp)
		if got != "1 hour ago" {
			t.Errorf("RelativeTime(60min) = %q, want %q", got, "1 hour ago")
		}
	})

	// Exactly at boundary between hour and day
	t.Run("24 hours", func(t *testing.T) {
		timestamp := time.Now().Add(-24 * time.Hour)
		got := RelativeTime(timestamp)
		if got != "1 day ago" {
			t.Errorf("RelativeTime(24h) = %q, want %q", got, "1 day ago")
		}
	})
}
