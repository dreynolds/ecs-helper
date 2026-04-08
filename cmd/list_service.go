package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
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
		fmt.Printf("%s %s\n", ui.KeyStyle.Render("Service:"), ui.ValueStyle.Render(valueOrEmpty(service.ServiceName)))
		fmt.Printf("%s %s\n", ui.KeyStyle.Render("Status:"), ui.ValueStyle.Render(valueOrEmpty(service.Status)))
		fmt.Printf("%s %s\n", ui.KeyStyle.Render("Desired Count:"), ui.ValueStyle.Render(fmt.Sprintf("%d", service.DesiredCount)))
		fmt.Printf("%s %s\n", ui.KeyStyle.Render("Running Count:"), ui.ValueStyle.Render(fmt.Sprintf("%d", service.RunningCount)))
		fmt.Printf("%s %s\n", ui.KeyStyle.Render("Pending Count:"), ui.ValueStyle.Render(fmt.Sprintf("%d", service.PendingCount)))
		fmt.Println("-------------------------------")
	}
	return nil
}
