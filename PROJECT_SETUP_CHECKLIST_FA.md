# چک‌لیست پیاده‌سازی پروژه به سبک Virasty (Task + MySQL)

این چک‌لیست برای ساخت همان معماری روی پروژه جدید است، با پشتیبانی کامل از MySQL.

## مرحله 1: ساخت اسکلت پروژه
1. پوشه‌های اصلی را بساز:
```bash
mkdir -p apiSchema/commonSchema apiSchema/taskSchema \
controllers/mainController controllers/task \
models/repositories models/task/datamodel models/task/dataSources/mysql \
services/core/route statics/constants/controllerBaseErrCode \
statics/constants/status statics/customErr tests/task_tests \
commands/user_migration
```
2. فایل `go.mod` را ایجاد کن و ماژول پروژه را تنظیم کن.
3. وابستگی‌ها را نصب کن:
```bash
go get github.com/gofiber/fiber/v2
go get github.com/go-sql-driver/mysql
```

## مرحله 2: تعریف مدل دامنه
1. فایل `models/task/datamodel/task.go` را بساز.
2. struct `Task` را با فیلدهای زیر تعریف کن:
- `ID int64`
- `Title string`
- `Description string`
- `CreatedAt string`

## مرحله 3: تعریف قراردادهای API
1. در `apiSchema/commonSchema/base.go` بنویس:
- `BaseRequest[T]` با `Body` و `Headers`
- `ValidateExtraData`
2. در `apiSchema/taskSchema/request.go` بساز:
- `CreateRequest`
- `ListRequest`
3. در `apiSchema/taskSchema/response.go` بساز:
- `CreateResponse`
- `ListResponse`
4. در `apiSchema/taskSchema/validate.go` اعتبارسنجی را اضافه کن:
- `title` خالی نباشد
- `page >= 1`
- `1 <= perPage <= 100`

## مرحله 4: تعریف repository interface
1. فایل `models/repositories/taskRepo.go` را بساز.
2. interface `TaskRepository` را تعریف کن:
- `Create(ctx, req)`
- `List(ctx, req)`
3. متغیر سراسری `var TaskRepo TaskRepository` را بگذار.

## مرحله 5: پیاده‌سازی datasource MySQL
1. `models/task/dataSources/mysql/config.go`:
- خواندن envها: `MYSQL_DSN`, `MYSQL_TASK_TABLE`
- تنظیمات pool: `MYSQL_MAX_OPEN_CONNS`, `MYSQL_MAX_IDLE_CONNS`, `MYSQL_CONN_MAX_LIFETIME_SECONDS`
2. `models/task/dataSources/mysql/connection.go`:
- `sql.Open` با driver `mysql`
- `Ping` برای تست اتصال
3. `models/task/dataSources/mysql/schema.go`:
- اعتبارسنجی نام جدول
- `EnsureTaskTable(...)` برای ساخت table
4. `models/task/dataSources/mysql/schema.sql`:
- اسکریپت SQL رسمی جدول task

## مرحله 6: پیاده‌سازی repository task
1. `models/task/repository_create.go`:
- struct `Repository` با پشتیبانی `db *sql.DB` و cache
- `GetRepo()` با `sync.Once`
- `init()` برای assign کردن `repositories.TaskRepo = GetRepo()`
- `Create(...)`:
  - در حالت MySQL: `INSERT` و سپس `SELECT` رکورد ایجادشده
  - در حالت بدون DSN: fallback به in-memory
  - invalidation کش لیست
2. `models/task/repository_list.go`:
- `List(...)` با cache key
- در حالت MySQL:
  - `SELECT ... LIMIT/OFFSET`
  - `COUNT(*)` برای `Total`
- در حالت بدون DSN: لیست از حافظه
- `cloneListResponse(...)`

## مرحله 7: توابع مشترک controller
1. فایل `controllers/mainController/main.go`:
- `InitAPI`
- `FinishAPISpan`
- `ParseBody`
- `ParseQuery`
- `Response`
- `Error`
- `fillHeaders`
- `validateBody`

## مرحله 8: کنترلرهای task
1. `controllers/task/create.go`:
- ParseBody
- call `TaskRepo.Create`
- return Response/Error
2. `controllers/task/list.go`:
- ParseQuery
- call `TaskRepo.List`
- return Response/Error

## مرحله 9: routeها
1. `services/core/route/task_route.go`:
- `POST /task/create`
- `GET /task/list`
2. `services/core/route/route.go`:
- `SetupRoutes` و merge کردن route maps

## مرحله 10: entrypoint سرویس
1. `services/core/main.go`:
- ساخت app با Fiber
- `SetupRoutes`
- اجرای `app.Listen(":8080")`
- blank import برای `models/task`

## مرحله 11: migration command
1. فایل `commands/user_migration/main.go`:
- خواندن `MYSQL_DSN` و `MYSQL_TASK_TABLE`
- امکان override با فلگ:
  - `--dsn`
  - `--table`
- اجرای `EnsureTaskTable(...)`
2. اجرای migration:
```bash
MYSQL_DSN="user:pass@tcp(127.0.0.1:3306)/sample" \
go run ./commands/user_migration
```

## مرحله 12: ثابت‌ها و خطاها
1. `statics/constants/controllerBaseErrCode/base.go`:
- `TaskErrCode`
2. `statics/constants/status/status.go`:
- status codeها
3. `statics/constants/error_message.go`:
- پیام خطاهای task
4. `statics/customErr/err.go`:
- `errors.New(...)` برای خطاهای دامنه task

## مرحله 13: تست
1. `tests/task_tests/task_create_test.go`
2. `tests/task_tests/task_list_test.go`
3. اجرای تست:
```bash
go test ./...
```

## مرحله 14: اجرای پروژه
1. اجرای سرویس با MySQL:
```bash
export MYSQL_DSN="user:pass@tcp(127.0.0.1:3306)/sample"
export MYSQL_TASK_TABLE="tasks"
go run ./services/core
```
2. تست دستی API:
```bash
curl -X POST http://localhost:8080/task/create \
  -H 'Content-Type: application/json' \
  -d '{"body":{"title":"Write docs","description":"draft"}}'
```
```bash
curl "http://localhost:8080/task/list?page=1&perPage=10"
```

## مرحله 15: چک نهایی
1. `gofmt` روی فایل‌ها اجرا شده باشد.
2. `go mod tidy` اجرا شده باشد.
3. `go test ./...` پاس باشد.
4. پاسخ‌ها envelope استاندارد `{data: ...}` داشته باشند.
5. خطاها کد `base+section+errStr` داشته باشند.
6. اگر `MYSQL_DSN` خالی بود، fallback in-memory سالم کار کند.
