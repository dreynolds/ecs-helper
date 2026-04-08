package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/davidtiberius/ecs-helper/internal/aws"
	"github.com/davidtiberius/ecs-helper/internal/ui"
	"github.com/spf13/cobra"
)

var (
	watchService  string
	watchInterval time.Duration
	watchTimeout  time.Duration
)

func init() {
	watchCmd.Flags().StringVar(&watchService, "service", "", "ECS service name")
	watchCmd.Flags().DurationVar(&watchInterval, "interval", 5*time.Second, "Polling interval")
	watchCmd.Flags().DurationVar(&watchTimeout, "timeout", 15*time.Minute, "How long to wait before timing out")
	_ = watchCmd.MarkFlagRequired("service")
	rootCmd.AddCommand(watchCmd)
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch ECS service deployment rollout",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cluster == "" {
			return errors.New("cluster is required; pass --cluster")
		}
		if watchInterval <= 0 {
			return errors.New("interval must be > 0")
		}
		if watchTimeout <= 0 {
			return errors.New("timeout must be > 0")
		}

		client, err := aws.NewClient(cmd.Context(), region)
		if err != nil {
			return fmt.Errorf("create ECS client: %w", err)
		}

		fmt.Println(ui.TitleStyle.Render(fmt.Sprintf("Watching deploy: cluster=%s service=%s region=%s", cluster, watchService, region)))

		deadline := time.Now().Add(watchTimeout)
		for {
			line, complete, failed, err := describeDeployment(cmd, client)
			if err != nil {
				return err
			}

			fmt.Printf("\r%-140s", line)
			if complete {
				fmt.Printf("\n%s\n", ui.SuccessStyle.Render("Deployment completed."))
				return nil
			}
			if failed {
				fmt.Printf("\n%s\n", ui.ErrorStyle.Render("Deployment failed."))
				return errors.New(strings.TrimSpace(line))
			}

			if time.Now().After(deadline) {
				fmt.Printf("\n%s\n", ui.ErrorStyle.Render("Timed out waiting for deployment to complete."))
				return fmt.Errorf("timed out after %s", watchTimeout)
			}

			select {
			case <-cmd.Context().Done():
				fmt.Println()
				return cmd.Context().Err()
			case <-time.After(watchInterval):
			}
		}
	},
}

func describeDeployment(cmd *cobra.Command, client *ecs.Client) (line string, complete bool, failed bool, err error) {
	out, err := client.DescribeServices(cmd.Context(), &ecs.DescribeServicesInput{
		Cluster:  &cluster,
		Services: []string{watchService},
	})
	if err != nil {
		return "", false, false, fmt.Errorf("describe service: %w", err)
	}
	if len(out.Failures) > 0 {
		f := out.Failures[0]
		return "", false, true, fmt.Errorf("service lookup failed: arn=%s reason=%s", valueOrEmpty(f.Arn), valueOrEmpty(f.Reason))
	}
	if len(out.Services) == 0 {
		return "", false, true, fmt.Errorf("service not found: %s", watchService)
	}

	svc := out.Services[0]
	primaryFound := false
	for _, d := range svc.Deployments {
		if valueOrEmpty(d.Status) != "PRIMARY" {
			continue
		}

		primaryFound = true
		rollout := string(d.RolloutState)
		reason := valueOrEmpty(d.RolloutStateReason)
		reason = strings.TrimSpace(reason)
		if reason == "" {
			reason = "-"
		}

		line = fmt.Sprintf(
			"rollout=%s desired=%d running=%d pending=%d reason=%s",
			rollout,
			d.DesiredCount,
			d.RunningCount,
			d.PendingCount,
			reason,
		)

		if rollout == "FAILED" {
			return line, false, true, nil
		}
		if rollout == "COMPLETED" {
			return line, true, false, nil
		}
		return line, false, false, nil
	}

	if !primaryFound {
		return "No PRIMARY deployment found yet", false, false, nil
	}
	return line, false, false, nil
}
