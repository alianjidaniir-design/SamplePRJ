# Virasty-Style Example Project

This project is rebuilt from scratch based on `VIRASTY_CODING_STYLE_GUIDE.md`.

## API Endpoints

Task:
- `POST /task/create`
- `GET /task/list?page=1&perPage=10`

## Run

```bash
go run ./services/core
```

### Run with MySQL table

```bash
export MYSQL_DSN="user:pass@tcp(127.0.0.1:3306)/sample"
export MYSQL_TASK_TABLE="tasks"
```

Create table:

```bash
go run ./commands/userMigration
```

Start API:

```bash
go run ./services/core
```

Notes:
- If `MYSQL_DSN` is empty, repository uses in-memory mode (useful for local tests).
- If `MYSQL_DSN` is set, API uses MySQL and automatically ensures the task table on startup.

## Example requests

```bash
curl -X POST http://localhost:8080/task/create \
  -H 'Content-Type: application/json' \
  -d '{"body":{"title":"Write docs","description":"draft guide"}}'
```

```bash
curl "http://localhost:8080/task/list?page=1&perPage=10"
```

## Test

```bash
go test ./...
```
