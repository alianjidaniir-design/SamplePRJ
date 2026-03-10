# چک‌لیست پیاده‌سازی پروژه به سبک Virasty (Task + MySQL)

این چک‌لیست مخصوص پروژه‌ای است که دامنه `task` دارد و همزمان از `MySQL` + `Memory fallback` استفاده می‌کند.

## مرحله 1: ساخت اسکلت پروژه
1. پوشه‌های اصلی را بساز:
```bash
mkdir -p apiSchema/commonSchema apiSchema/taskSchema \
controllers/mainController controllers/task \
models/repositories models/task/dataModel \
models/task/dataSources/memoryDS models/task/dataSources/mysqlDS \
services/core/route statics/constants/controllerBaseErrCode \
statics/constants/status statics/customErr tests/task_tests commands
```
2. فایل `go.mod` را ایجاد کن.
3. وابستگی‌ها را نصب کن:
```bash
go get github.com/gofiber/fiber/v2
go get github.com/go-sql-driver/mysql
```

## مرحله 2: تعریف مدل دامنه
1. فایل `models/task/dataModel/task.go` را بساز.
2. struct `Task` را تعریف کن:
- `ID int64`
- `Title string`
- `Description string`
- `CreatedAt string`

## مرحله 3: تعریف قراردادهای API
1. فایل `apiSchema/commonSchema/base.go`:
- `BaseRequest[T]`
- `ValidateExtraData`
2. فایل `apiSchema/taskSchema/request.go`:
- `CreateRequest`
- `ListRequest`
3. فایل `apiSchema/taskSchema/response.go`:
- `CreateResponse`
- `ListResponse`
4. فایل `apiSchema/taskSchema/validate.go`:
- `title` خالی نباشد
- `page >= 1`
- `1 <= perPage <= 100`

## مرحله 4: تعریف Repository Interface
1. فایل `models/repositories/taskRepo.go`:
- `Create(ctx, req)`
- `List(ctx, req)`
- `var TaskRepo TaskRepository`

## مرحله 5: تعریف DataSource Contract
1. فایل `models/task/dataSources/taskDS.go`:
- interface `TaskDBDS`
- interface `TaskCacheDS`

## مرحله 6: پیاده‌سازی DataSourceهای Memory
1. فایل `models/task/dataSources/memoryDS/taskDBDS.go`:
- نگهداری taskها در حافظه
- `CreateTask`
- `ListTasks`
- `Reset`
2. فایل `models/task/dataSources/memoryDS/taskCacheDS.go`:
- `GetList`
- `SetList`
- `InvalidateList`
- `Reset`

## مرحله 7: پیاده‌سازی DataSourceهای MySQL
1. فایل `models/task/dataSources/mysqlDS/config.go`:
- خواندن envها: `MYSQL_DSN`, `MYSQL_TASK_TABLE`
- تنظیمات pool: `MYSQL_MAX_OPEN_CONNS`, `MYSQL_MAX_IDLE_CONNS`, `MYSQL_CONN_MAX_LIFETIME_SECONDS`
2. فایل `models/task/dataSources/mysqlDS/connection.go`:
- اتصال با `database/sql`
- driver: `go-sql-driver/mysql`
3. فایل `models/task/dataSources/mysqlDS/schema.go`:
- validate نام جدول
- `EnsureTaskTable(...)`
4. فایل `models/task/dataSources/mysqlDS/schema.sql`:
- SQL schema جدول task
5. فایل `models/task/dataSources/mysqlDS/taskDBDS.go`:
- `CreateTask` با INSERT + SELECT
- `ListTasks` با LIMIT/OFFSET + COUNT

## مرحله 8: پیاده‌سازی Repository Task
1. فایل `models/task/repository.go`:
- singleton با `sync.Once`
- مقداردهی `repositories.TaskRepo` در `init()`
- انتخاب datasource:
  - اگر `MYSQL_DSN` ست باشد: mysql datasource
  - اگر خالی باشد: memory datasource
- اتصال cache datasource حافظه
2. فایل `models/task/repositoryCreate.go`:
- create task از طریق db datasource
- invalidation کش لیست
3. فایل `models/task/repositoryList.go`:
- cache key بر اساس `page/perPage`
- cache hit/miss
- خواندن list از db datasource
- `cloneListResponse`
4. فایل `models/task/repositoryCache_test.go`:
- تست رفتار cache و invalidation

## مرحله 9: Controller مشترک
1. فایل `controllers/mainController/main.go`:
- `InitAPI`
- `FinishAPISpan`
- `ParseBody`
- `ParseQuery`
- `Response`
- `Error`
- `fillHeaders`
- `validateBody`

## مرحله 10: Controllerهای Task
1. فایل `controllers/task/create.go`:
- ParseBody
- call `TaskRepo.Create`
- return Response/Error
2. فایل `controllers/task/list.go`:
- ParseQuery
- call `TaskRepo.List`
- return Response/Error

## مرحله 11: Routeها
1. فایل `services/core/route/taskRoute.go`:
- `POST /task/create`
- `GET /task/list`
2. فایل `services/core/route/route.go`:
- `SetupRoutes`
- merge route maps

## مرحله 12: Entrypoint سرویس
1. فایل `services/core/main.go`:
- ساخت Fiber app
- `SetupRoutes`
- `app.Listen(":8080")`
- blank import برای `models/task`

## مرحله 13: Constants و Errorها
1. `statics/constants/controllerBaseErrCode/base.go`
2. `statics/constants/status/status.go`
3. `statics/constants/errorMessage.go`
4. `statics/customErr/err.go`

## مرحله 14: Migration command
1. فایل `commands/userMigration/main.go`:
- خواندن DSN از env یا فلگ `--dsn`
- خواندن table از env یا فلگ `--table`
- `EnsureTaskTable(...)` برای ساخت/اعتبارسنجی جدول
2. اجرا:
```bash
MYSQL_DSN="user:pass@tcp(127.0.0.1:3306)/sample?multiStatements=true" \
go run ./commands/userMigration --table tasks
```

## مرحله 15: تست
1. `tests/task_tests/taskCreate_test.go`
2. `tests/task_tests/taskList_test.go`
3. اجرای تست:
```bash
go test ./...
```

## مرحله 16: اجرای پروژه
1. اجرای API با fallback حافظه:
```bash
go run ./services/core
```
2. اجرای API با MySQL:
```bash
export MYSQL_DSN="user:pass@tcp(127.0.0.1:3306)/sample"
export MYSQL_TASK_TABLE="tasks"
go run ./services/core
```
3. تست دستی API:
```bash
curl -X POST http://localhost:8080/task/create \
  -H 'Content-Type: application/json' \
  -d '{"body":{"title":"Write docs","description":"draft"}}'
```
```bash
curl "http://localhost:8080/task/list?page=1&perPage=10"
```

## مرحله 17: چک نهایی
1. `gofmt` اجرا شده باشد.
2. `go mod tidy` اجرا شده باشد.
3. `go test ./...` پاس باشد.
4. پاسخ موفق envelope استاندارد `{data: ...}` داشته باشد.
5. خطاها با فرمت `base + section + errStr` برگردند.
6. هر دو mode (MySQL و Memory fallback) سالم کار کنند.
