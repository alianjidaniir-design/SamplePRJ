# چک‌لیست پیاده‌سازی پروژه به سبک Virasty (Task بدون MySQL)

این چک‌لیست برای ساخت همان معماری روی پروژه جدید است، فقط با datasource حافظه (memory).

## مرحله 1: ساخت اسکلت پروژه
1. پوشه‌های اصلی را بساز:
```bash
mkdir -p apiSchema/commonSchema apiSchema/taskSchema \
controllers/mainController controllers/task \
models/repositories models/task/dataModel models/task/dataSources/memory \
services/core/route statics/constants/controllerBaseErrCode \
statics/constants/status statics/customErr tests/task_tests commands
```
2. فایل `go.mod` را ایجاد کن.
3. Fiber را نصب کن:
```bash
go get github.com/gofiber/fiber/v2
```

## مرحله 2: تعریف مدل دامنه
1. `models/task/dataModel/task.go` را بساز.
2. struct `Task` را تعریف کن.

## مرحله 3: تعریف قراردادهای API
1. `apiSchema/commonSchema/base.go`:
- `BaseRequest[T]`
- `ValidateExtraData`
2. `apiSchema/taskSchema/request.go`:
- `CreateRequest`
- `ListRequest`
3. `apiSchema/taskSchema/response.go`:
- `CreateResponse`
- `ListResponse`
4. `apiSchema/taskSchema/validate.go`:
- title خالی نباشد
- page/perPage معتبر باشند

## مرحله 4: repository interface
1. `models/repositories/taskRepo.go`:
- `Create(ctx, req)`
- `List(ctx, req)`
- `var TaskRepo TaskRepository`

## مرحله 5: dataSources memory
1. `models/task/dataSources/taskDS.go`:
- interfaceهای `TaskDBDS` و `TaskCacheDS`
2. `models/task/dataSources/memory/taskDBDS.go`:
- نگهداری taskها داخل حافظه
- create/list با lock
3. `models/task/dataSources/memory/taskCacheDS.go`:
- کش لیست با key صفحه‌بندی

## مرحله 6: repository task
1. `models/task/repository.go`:
- singleton repo با `sync.Once`
- assign شدن `repositories.TaskRepo` در `init()`
- datasource حافظه به‌صورت پیش‌فرض
2. `models/task/repositoryCreate.go`:
- create task و invalidation کش
3. `models/task/repositoryList.go`:
- cache hit/miss + pagination

## مرحله 7: controllerهای مشترک
1. `controllers/mainController/main.go`:
- `InitAPI`
- `ParseBody`
- `ParseQuery`
- `Response`
- `Error`

## مرحله 8: controllerهای task
1. `controllers/task/create.go`
2. `controllers/task/list.go`

## مرحله 9: routeها
1. `services/core/route/taskRoute.go`
- `POST /task/create`
- `GET /task/list`
2. `services/core/route/route.go`

## مرحله 10: entrypoint
1. `services/core/main.go`
- ساخت app
- رجیستر route
- `app.Listen(":8080")`

## مرحله 11: constants و errors
1. `statics/constants/controllerBaseErrCode/base.go`
2. `statics/constants/status/status.go`
3. `statics/constants/errorMessage.go`
4. `statics/customErr/err.go`

## مرحله 12: تست و اجرا
1. `go test ./...`
2. `go run ./services/core`
3. تست دستی API:
```bash
curl -X POST http://localhost:8080/task/create \
  -H 'Content-Type: application/json' \
  -d '{"body":{"title":"Write docs","description":"draft"}}'
```
```bash
curl "http://localhost:8080/task/list?page=1&perPage=10"
```
