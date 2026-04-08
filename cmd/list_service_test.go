package cmd

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func TestSummarizeService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		service     types.Service
		wantName    string
		wantStatus  string
		wantDesired int32
		wantRunning int32
		wantPending int32
	}{
		{
			name: "all fields populated",
			service: types.Service{
				ServiceName:  strPtr("api-service"),
				Status:       strPtr("ACTIVE"),
				DesiredCount: 3,
				RunningCount: 2,
				PendingCount: 1,
			},
			wantName:    "api-service",
			wantStatus:  "ACTIVE",
			wantDesired: 3,
			wantRunning: 2,
			wantPending: 1,
		},
		{
			name: "nil pointers fall back to empty strings",
			service: types.Service{
				ServiceName:  nil,
				Status:       nil,
				DesiredCount: 0,
				RunningCount: 0,
				PendingCount: 0,
			},
			wantName:    "",
			wantStatus:  "",
			wantDesired: 0,
			wantRunning: 0,
			wantPending: 0,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			name, status, desired, running, pending := summarizeService(tc.service)
			if name != tc.wantName {
				t.Fatalf("name=%q want %q", name, tc.wantName)
			}
			if status != tc.wantStatus {
				t.Fatalf("status=%q want %q", status, tc.wantStatus)
			}
			if desired != tc.wantDesired {
				t.Fatalf("desired=%d want %d", desired, tc.wantDesired)
			}
			if running != tc.wantRunning {
				t.Fatalf("running=%d want %d", running, tc.wantRunning)
			}
			if pending != tc.wantPending {
				t.Fatalf("pending=%d want %d", pending, tc.wantPending)
			}
		})
	}
}
