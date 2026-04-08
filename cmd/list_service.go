package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/davidtiberius/ecs-helper/internal/aws"
	"github.com/davidtiberius/ecs-helper/internal/ui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listServiceCmd)
}

var listServiceCmd = &cobra.Command{
	Use:   "list-service",
	Short: "List ECS services",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := aws.NewClient(cmd.Context(), region)
		if err != nil {
			return fmt.Errorf("create ECS client: %w", err)
		}
		out, err := client.ListServices(cmd.Context(), &ecs.ListServicesInput{
			Cluster: &cluster,
		})
		if err != nil {
			return fmt.Errorf("list services: %w", err)
		}

		return describeService(cmd, client, out.ServiceArns)
	},
}

func describeService(cmd *cobra.Command, client *ecs.Client, services []string) error {
	out, err := client.DescribeServices(cmd.Context(), &ecs.DescribeServicesInput{
		Cluster:  &cluster,
		Services: services,
	})
	if err != nil {
		return fmt.Errorf("describe services: %w", err)
	}
	for _, service := range out.Services {
		name, status, desired, running, pending := summarizeService(service)
		fmt.Printf("%s %s\n", ui.KeyStyle.Render("Service:"), ui.ValueStyle.Render(name))
		fmt.Printf("%s %s\n", ui.KeyStyle.Render("Status:"), ui.ValueStyle.Render(status))
		fmt.Printf("%s %s\n", ui.KeyStyle.Render("Desired Count:"), ui.ValueStyle.Render(fmt.Sprintf("%d", desired)))
		fmt.Printf("%s %s\n", ui.KeyStyle.Render("Running Count:"), ui.ValueStyle.Render(fmt.Sprintf("%d", running)))
		fmt.Printf("%s %s\n", ui.KeyStyle.Render("Pending Count:"), ui.ValueStyle.Render(fmt.Sprintf("%d", pending)))
		fmt.Println("-------------------------------")
	}
	return nil
}

func summarizeService(service types.Service) (name, status string, desired, running, pending int32) {
	return valueOrEmpty(service.ServiceName), valueOrEmpty(service.Status), service.DesiredCount, service.RunningCount, service.PendingCount
}
