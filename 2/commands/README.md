# commands

CLI and migration scripts for maintenance jobs.

## Available scripts

- `commands/elasticSearchReindex/main.go`
- `commands/statsUpdate/main.go`
- `commands/userMigration/main.go`

## Run

```bash
go run ./commands/elasticSearchReindex --index tasks --batch 500
```

```bash
go run ./commands/statsUpdate --period daily --dry-run
```

```bash
MYSQL_DSN="user:pass@tcp(127.0.0.1:3306)/sample?multiStatements=true" \
go run ./commands/userMigration --table tasks
```
