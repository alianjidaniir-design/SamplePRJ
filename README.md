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
