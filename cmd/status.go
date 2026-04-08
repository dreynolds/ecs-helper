package cmd

import (
    "fmt"

    "github.com/aws/aws-sdk-go-v2/service/ecs"
    "github.com/charmbracelet/lipgloss"
    "github.com/davidtiberius/ecs-helper/internal/aws"
    "github.com/spf13/cobra"
)

func init() {
    rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
    Use:   "status",
    Short: "Show running services and task counts",
    RunE: func(cmd *cobra.Command, args []string) error {
        titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
        keyStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("69"))
        valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
        okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)

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
            fmt.Println(titleStyle.Render("ECS Status"))
            fmt.Println("No clusters returned for request.")
            return nil
        }

        c := out.Clusters[0]
        fmt.Println(titleStyle.Render("ECS Status"))
        fmt.Printf("%s %s\n", keyStyle.Render("Cluster:"), valueStyle.Render(cluster))
        fmt.Printf("%s %s\n", keyStyle.Render("Region:"), valueStyle.Render(region))
        fmt.Printf("%s %s\n", keyStyle.Render("Status:"), okStyle.Render(valueOrEmpty(c.Status)))
        fmt.Printf("%s %d\n", keyStyle.Render("Registered container instances:"), c.RegisteredContainerInstancesCount)
        fmt.Printf("%s %d\n", keyStyle.Render("Running tasks:"), c.RunningTasksCount)
        fmt.Printf("%s %d\n", keyStyle.Render("Pending tasks:"), c.PendingTasksCount)
        fmt.Printf("%s %d\n", keyStyle.Render("Active services:"), c.ActiveServicesCount)

        if len(out.Failures) > 0 {
            warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
            fmt.Printf("\n%s\n", warnStyle.Render("Failures"))
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