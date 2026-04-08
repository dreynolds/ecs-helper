# ecs-helper

A CLI for inspecting and monitoring Amazon ECS clusters and services.

## Installation

```bash
go install github.com/davidtiberius/ecs-helper@latest
```

Or build from source:

```bash
git clone https://github.com/davidtiberius/ecs-helper
cd ecs-helper
go build -o ecs-helper .
```

## Authentication

ecs-helper uses the standard AWS credential chain (environment variables, `~/.aws/credentials`, IAM roles, etc.). Make sure your credentials are configured before use.

## Global flags

| Flag        | Default      | Description                      |
|-------------|--------------|----------------------------------|
| `--cluster` |              | ECS cluster name                 |
| `--region`  | `eu-west-1`  | AWS region                       |
| `--env`     | `prod`       | Environment (`prod` / `staging`) |

## Commands

### `status`

Show the cluster's task and service counts.

```bash
ecs-helper status --cluster my-cluster
```

### `list-service`

List all services in a cluster with their desired, running, and pending task counts.

```bash
ecs-helper list-service --cluster my-cluster
```

### `watch`

Poll a service's deployment until it completes or fails, then exit with a non-zero status on failure. Useful in CI pipelines after triggering a deploy.

```bash
ecs-helper watch --cluster my-cluster --service my-service
```

| Flag         | Default | Description                          |
|--------------|---------|--------------------------------------|
| `--service`  |         | ECS service name (required)          |
| `--interval` | `5s`    | How often to poll                    |
| `--timeout`  | `15m`   | Maximum time to wait before giving up|

## CI

[![CI](https://github.com/davidtiberius/ecs-helper/actions/workflows/ci.yml/badge.svg)](https://github.com/davidtiberius/ecs-helper/actions/workflows/ci.yml)

Tests run automatically on every push and pull request via GitHub Actions.

```bash
go test ./...
```
