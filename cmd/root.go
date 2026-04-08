package cmd

import (
    "github.com/spf13/cobra"
    // "github.com/davidtiberius/ecs-helper/internal/config"
)

var (
    cluster string
    region  string
    env     string
)

var rootCmd = &cobra.Command{
    Use:   "ecs-helper",
    Short: "ECS management CLI",
}

func Execute() {
    rootCmd.Execute()
}

func init() {
    rootCmd.PersistentFlags().StringVar(&cluster, "cluster", "", "ECS cluster")
    rootCmd.PersistentFlags().StringVar(&region, "region", "eu-west-1", "AWS region")
    rootCmd.PersistentFlags().StringVar(&env, "env", "prod", "Environment (prod/staging)")
}