# commands

CLI and migration scripts for maintenance jobs.

## Available scripts

- `commands/elasticsearch_reindex/main.go`
- `commands/stats_update/main.go`
- `commands/user_migration/main.go`

## Run

```bash
go run ./commands/elasticsearch_reindex --index tasks --batch 500
```

```bash
go run ./commands/stats_update --period daily --dry-run
```

```bash
MYSQL_DSN="user:pass@tcp(127.0.0.1:3306)/sample?multiStatements=true" \
go run ./commands/user_migration --table tasks
```
