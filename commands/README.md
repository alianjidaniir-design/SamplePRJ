# commands

CLI and migration scripts for maintenance jobs.

## Available scripts

- `commands/elasticsearch_reindex/main.go`
- `commands/stats_update/main.go`

## Run

```bash
go run ./commands/elasticsearch_reindex --index tasks --batch 500
```

```bash
go run ./commands/stats_update --period daily --dry-run
```
