# Project Explanation (Virasty-Style Fiber API)

## 1) Project Overview
This project is a layered Go API built with **Fiber** and designed using the Virasty coding style guide.

Current domains:
- `task`
- `user`

Current endpoints:
- `POST /task/create`
- `POST /task/list`
- `POST /user/create`
- `POST /user/info`

---

## 2) Tech Stack
- Language: Go
- Framework: Fiber (`github.com/gofiber/fiber/v2`)
- Architecture: Layered (Schema -> Controller -> Repository -> Route -> Service)
- Storage: In-memory repository state (no real DB table yet)
- Cache: In-memory list cache for `task/list`

---

## 3) Directory Structure
- `apiSchema/`: request/response/validation contracts per domain
- `controllers/`: HTTP handlers per endpoint
- `models/`: data models, repositories, and repository interfaces
- `services/core/`: main service entrypoint and route registration
- `statics/`: constants, error codes, custom errors
- `tests/`: API tests by domain

---

## 4) Runtime Flow (Request Lifecycle)
Example: `POST /task/list`

1. Route registration maps `/task/list` to `controllers/task/List`.
2. Controller initializes API context (`InitAPI`) and parses body with `ParseBody`.
3. `ParseBody` also triggers schema validation via `Validate(...)`.
4. Controller calls repository method: `repositories.TaskRepo.List(...)`.
5. Repository checks cache (cache hit/miss logic).
6. Repository returns response model.
7. Controller returns standardized JSON envelope via `Response(...)`.

Error flow:
- Parse/validation errors -> section `01`
- Repository/business errors -> section `02`
- Final error code format: `baseErrCode + section + errStr`

---

## 5) Layer Responsibilities

### 5.1 Schema Layer (`apiSchema`)
Defines domain contracts:
- `request.go`: input payload structs
- `response.go`: output payload structs
- `validate.go`: business validation rules

Base wrapper:
- `commonSchema.BaseRequest[T]`

### 5.2 Controller Layer (`controllers`)
Handles HTTP orchestration only:
- parse request
- call repository
- return standardized response/error

No direct storage logic in controllers.

### 5.3 Repository Layer (`models`)
Contains business/data handling logic:
- `models/repositories/*`: interfaces (`TaskRepository`, `UserRepository`)
- `models/task`, `models/user`: concrete singleton implementations (`GetRepo()` + `sync.Once`)

### 5.4 Route Layer (`services/core/route`)
Registers endpoints in route maps (`taskRoutes`, `userRoutes`) and binds them to handlers.

### 5.5 Service Entry (`services/core/main.go`)
Bootstraps Fiber app and starts the API server on `:8080`.

---

## 6) Caching Design
Implemented in task repository:
- Cache target: `task/list`
- Cache key: `task:list:page:<page>:perPage:<perPage>`
- Cache store: in-memory map in `models/task/repository.go`
- Invalidation: full list-cache reset after `task/create`

Safety detail:
- Returned cached response is cloned to avoid accidental mutation side effects.

---

## 7) Data Model Status
Current status:
- Data is kept in memory (`[]Task`, `[]User`) inside repositories.
- No database/table persistence is active yet.

Implication:
- Data resets on service restart.
- Suitable for learning/demo phase.

---

## 8) Testing Strategy
API-level tests are implemented for both domains:
- task create/list tests
- user create/info tests

Additional repository-level test validates:
- cache population
- cache hit behavior
- cache invalidation after create

---

## 9) Run and Verify
Run server:
```bash
go run ./services/core
```

Run tests:
```bash
go test ./...
```

---

## 10) Current Limitations
- No real DB datasource yet (PostgreSQL/MySQL not integrated)
- No migration tooling yet
- No auth middleware yet
- Cache is local-memory only (not distributed)

---

## 11) Recommended Next Steps
1. Add DB persistence (e.g., PostgreSQL + GORM) and replace in-memory repositories.
2. Add migration scripts under `commands/`.
3. Move cache to Redis for multi-instance scalability.
4. Add auth/app-scope middleware in `middleware/`.
5. Add integration tests for error-code contract and edge cases.
