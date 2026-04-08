package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/dreynolds/ecs-helper/internal/aws"
	"github.com/dreynolds/ecs-helper/internal/ui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show running services and task counts",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := aws.NewClient(cmd.Context(), region)
		if err != nil {
			return err
		}

		out, err := client.DescribeClusters(cmd.Context(), &ecs.DescribeClustersInput{
			Clusters: []string{cluster},
		})
		if err != nil {
			return err
		}

		if len(out.Clusters) == 0 {
			fmt.Println(ui.TitleStyle.Render("ECS Status"))
			fmt.Println("No clusters returned for request.")
			return nil
		}

		c := out.Clusters[0]
		fmt.Println(ui.TitleStyle.Render("ECS Status"))
		fmt.Printf("%s %s\n", ui.KeyStyle.Render("Cluster:"), ui.ValueStyle.Render(cluster))
		fmt.Printf("%s %s\n", ui.KeyStyle.Render("Region:"), ui.ValueStyle.Render(region))
		fmt.Printf("%s %s\n", ui.KeyStyle.Render("Status:"), ui.SuccessStyle.Render(valueOrEmpty(c.Status)))
		fmt.Printf("%s %d\n", ui.KeyStyle.Render("Registered container instances:"), c.RegisteredContainerInstancesCount)
		fmt.Printf("%s %d\n", ui.KeyStyle.Render("Running tasks:"), c.RunningTasksCount)
		fmt.Printf("%s %d\n", ui.KeyStyle.Render("Pending tasks:"), c.PendingTasksCount)
		fmt.Printf("%s %d\n", ui.KeyStyle.Render("Active services:"), c.ActiveServicesCount)

		if len(out.Failures) > 0 {
			fmt.Printf("\n%s\n", ui.WarnStyle.Render("Failures"))
			for _, f := range out.Failures {
				fmt.Printf("- arn=%s reason=%s\n", valueOrEmpty(f.Arn), valueOrEmpty(f.Reason))
			}
		}

		return nil
	},
}

func valueOrEmpty(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
