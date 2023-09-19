package controllers

import (
	"testing"
	"time"
)

func TestUpdateInterval(t *testing.T) {
	currentSyncInterval := 5 * time.Minute
	rep := ReplicationConfigReconciler{
		SyncInterval: currentSyncInterval,
	}

	for _, tt := range []struct {
		name     string
		interval string
		expected time.Duration
	}{
		{
			name:     "use default interval",
			interval: "",
			expected: currentSyncInterval,
		},
		{
			name:     "use custom interval",
			interval: "10",
			expected: 10 * time.Minute,
		},
		{
			name:     "use none valid custom interval",
			interval: "3",
			expected: 5 * time.Minute,
		},
		{
			name:     "use none valid custom interval",
			interval: "undefined",
			expected: 5 * time.Minute,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			actual := rep.parseSyncInterval(tt.interval)
			if actual != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, actual)
			}
		})
	}
}
