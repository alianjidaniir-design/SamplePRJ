# Virasty-Style Example Project

This project is rebuilt from scratch based on `VIRASTY_CODING_STYLE_GUIDE.md`.

## Project idea

A simple API with two domains:
- Task domain:
  - `POST /task/create`
  - `POST /task/list`
- User domain:
  - `POST /user/create`
  - `POST /user/info`

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
curl -X POST http://localhost:8080/task/list \
  -H 'Content-Type: application/json' \
  -d '{"body":{"page":1,"perPage":10}}'
```

```bash
curl -X POST http://localhost:8080/user/create \
  -H 'Content-Type: application/json' \
  -d '{"body":{"username":"virasty","email":"virasty@example.com"}}'
```

```bash
curl -X POST http://localhost:8080/user/info \
  -H 'Content-Type: application/json' \
  -d '{"body":{"userID":11}}'
```

## Test

```bash
go test ./...
```
