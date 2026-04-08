package aws

import (
    "context"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/ecs"
)

func NewClient(ctx context.Context, region string) (*ecs.Client, error) {
    cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
    if err != nil {
        return nil, err
    }
    return ecs.NewFromConfig(cfg), nil
}