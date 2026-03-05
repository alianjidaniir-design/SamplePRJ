# commands

CLI and migration scripts for maintenance jobs.

## Available scripts

- `commands/elasticsearch_reindex/main.go`
- `commands/user_migration/main.go`
- `commands/stats_update/main.go`

## Run

```bash
go run ./commands/elasticsearch_reindex --index users --batch 500
```

```bash
go run ./commands/user_migration --from v1 --to v2 --limit 1000
```

```bash
go run ./commands/stats_update --period daily --dry-run
```
