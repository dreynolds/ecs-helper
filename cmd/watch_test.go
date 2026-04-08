package cmd

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func TestSummarizePrimaryDeployment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		deployments  []types.Deployment
		wantLinePart string
		wantComplete bool
		wantFailed   bool
	}{
		{
			name: "completed rollout",
			deployments: []types.Deployment{
				{
					Status:             strPtr("PRIMARY"),
					RolloutState:       types.DeploymentRolloutStateCompleted,
					RolloutStateReason: strPtr("all tasks healthy"),
					DesiredCount:       4,
					RunningCount:       4,
					PendingCount:       0,
				},
			},
			wantLinePart: "rollout=COMPLETED",
			wantComplete: true,
			wantFailed:   false,
		},
		{
			name: "failed rollout",
			deployments: []types.Deployment{
				{
					Status:             strPtr("PRIMARY"),
					RolloutState:       types.DeploymentRolloutStateFailed,
					RolloutStateReason: strPtr("service failed ELB health checks"),
					DesiredCount:       2,
					RunningCount:       1,
					PendingCount:       1,
				},
			},
			wantLinePart: "rollout=FAILED",
			wantComplete: false,
			wantFailed:   true,
		},
		{
			name: "in progress rollout",
			deployments: []types.Deployment{
				{
					Status:             strPtr("PRIMARY"),
					RolloutState:       types.DeploymentRolloutStateInProgress,
					RolloutStateReason: strPtr(""),
					DesiredCount:       3,
					RunningCount:       2,
					PendingCount:       1,
				},
			},
			wantLinePart: "rollout=IN_PROGRESS",
			wantComplete: false,
			wantFailed:   false,
		},
		{
			name: "no primary deployment",
			deployments: []types.Deployment{
				{
					Status:             strPtr("ACTIVE"),
					RolloutState:       types.DeploymentRolloutStateInProgress,
					RolloutStateReason: strPtr(""),
				},
			},
			wantLinePart: "No PRIMARY deployment found yet",
			wantComplete: false,
			wantFailed:   false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			line, complete, failed := summarizePrimaryDeployment(tc.deployments)
			if !strings.Contains(line, tc.wantLinePart) {
				t.Fatalf("line %q does not contain %q", line, tc.wantLinePart)
			}
			if complete != tc.wantComplete {
				t.Fatalf("complete=%v want %v", complete, tc.wantComplete)
			}
			if failed != tc.wantFailed {
				t.Fatalf("failed=%v want %v", failed, tc.wantFailed)
			}
		})
	}
}

func TestValueOrEmpty(t *testing.T) {
	t.Parallel()

	if got := valueOrEmpty(nil); got != "" {
		t.Fatalf("valueOrEmpty(nil)=%q want empty string", got)
	}

	s := "hello"
	if got := valueOrEmpty(&s); got != "hello" {
		t.Fatalf("valueOrEmpty(&s)=%q want hello", got)
	}
}

func strPtr(v string) *string {
	return &v
}
