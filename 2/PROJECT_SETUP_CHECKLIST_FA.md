# چک‌لیست ساخت پروژه Virasty (Task + MySQL + Memory Fallback)

این فایل برای ساخت/بازسازی همین پروژه داخل پوشه `2` نوشته شده است: دامنه `task` با Fiber، کش in-memory برای `task/list`، و دیتاسورس MySQL (در صورت تنظیم `MYSQL_DSN`) با fallback حافظه.

## هدف پوشه‌ها

- `apiSchema/`: قرارداد API (request/response/validate) و مدل‌های ورودی/خروجی.
- `controllers/`: هندلرهای HTTP؛ فقط orchestration (parse/validate -> repo -> response/error).
- `models/`: هسته business/data.
- `models/repositories/`: interface repoها و متغیرهای global برای inject (مثل `TaskRepo`).
- `models/task/`: پیاده‌سازی use-caseهای دامنه task (Create/List) و singleton repo.
- `models/task/dataSources/`: contract دیتاسورس‌ها + پیاده‌سازی‌ها.
- `models/task/dataSources/memoryDS/`: دیتابیس و کش حافظه (برای تست/لوکال).
- `models/task/dataSources/mysqlDS/`: اتصال MySQL، schema helper، و پیاده‌سازی datasource.
- `services/`: entrypoint سرویس و route registration.
- `statics/`: کدهای status، base error code، پیام خطا و custom errorها.
- `tests/`: تست‌های API-level برای endpointها.
- `commands/`: اسکریپت‌های CLI (migration/reindex/stats).
- `middleware/`: جای middlewareهای سراسری (فعلاً reserved).
- `pkg/`: پکیج‌های مشترک (فعلاً reserved).
- `template/`: templateهای feature (فعلاً reserved).

## هدف فایل‌های کلیدی (Task + MySQL)

- `services/core/main.go`: ورودی اصلی API؛ ساخت app Fiber، رجیستر routeها، اجرا روی `:8080`.
- `services/core/route/route.go`: تجمیع route mapها (در این پروژه فقط task).
- `services/core/route/taskRoute.go`: تعریف و رجیستر endpointها:
  - `POST /task/create`
  - `GET /task/list?page=1&perPage=10`
- `controllers/mainController/main.go`: ابزارهای مشترک controller:
  - `InitAPI` و `FinishAPISpan` (فعلاً ساده)
  - `ParseBody` برای POST body
  - `ParseQuery` برای GET query
  - `Response` و `Error` برای پاسخ استاندارد
- `controllers/task/create.go`: جریان `POST /task/create` (ParseBody -> Validate -> Repo.Create -> Response/Error).
- `controllers/task/list.go`: جریان `GET /task/list` (ParseQuery -> Validate -> Repo.List -> Response/Error).
- `apiSchema/commonSchema/base.go`: wrapper `BaseRequest[T]` و `ValidateExtraData`.
- `apiSchema/taskSchema/request.go`: ورودی‌های `CreateRequest` و `ListRequest`.
- `apiSchema/taskSchema/validate.go`: قوانین اعتبارسنجی (title/page/perPage).
- `apiSchema/taskSchema/response.go`: خروجی‌های `CreateResponse` و `ListResponse`.
- `models/repositories/taskRepo.go`: interface `TaskRepository` و `var TaskRepo TaskRepository`.
- `models/task/repository.go`: singleton repository + انتخاب دیتاسورس:
  - اگر `MYSQL_DSN` ست باشد: `mysqlDS.NewTaskDBDSFromEnv()`
  - اگر خالی باشد: `memoryDS.NewTaskDBDS(...)`
  - کش list همیشه in-memory است (`memoryDS.NewTaskCacheDS()`).
- `models/task/repositoryCreate.go`: use-case create:
  - ساخت task از طریق db datasource
  - `InvalidateList()` روی cache datasource
- `models/task/repositoryList.go`: use-case list:
  - cache hit/miss
  - خواندن list از db datasource
  - clone response برای جلوگیری از mutation
- `models/task/dataSources/taskDS.go`: interfaceهای `TaskDBDS` و `TaskCacheDS`.
- `models/task/dataSources/memoryDS/taskDBDS.go`: پیاده‌سازی DB حافظه (Create/List + Reset برای تست).
- `models/task/dataSources/memoryDS/taskCacheDS.go`: کش حافظه (Get/Set/Invalidate + Reset).
- `models/task/dataSources/mysqlDS/config.go`: خواندن env و نرمال‌سازی DSN.
- `models/task/dataSources/mysqlDS/connection.go`: ساخت و تنظیم pool `database/sql` + import driver MySQL.
- `models/task/dataSources/mysqlDS/schema.go` و `schema.sql`: validate نام جدول + helper ساخت جدول.
- `models/task/dataSources/mysqlDS/taskDBDS.go`: پیاده‌سازی MySQL (INSERT/SELECT + LIMIT/OFFSET + COUNT).
- `commands/userMigration/main.go`: ساخت/اطمینان از وجود جدول task در MySQL (برای migration اولیه).
- `tests/task_tests/taskCreate_test.go`: تست `POST /task/create`.
- `tests/task_tests/taskList_test.go`: تست `GET /task/list`.

## چک‌لیست اجرا (روی پروژه جدید)

1. ساخت پوشه‌ها مطابق ساختار بالا.
2. تعریف مدل `Task` در `models/task/dataModel/task.go`.
3. تعریف schemaها و validateها در `apiSchema/taskSchema/`.
4. تعریف `TaskRepository` در `models/repositories/taskRepo.go`.
5. پیاده‌سازی دیتاسورس‌های memory و mysql در `models/task/dataSources/`.
6. پیاده‌سازی repository singleton و use-caseها در `models/task/`.
7. پیاده‌سازی controllerهای مشترک و controllerهای task.
8. رجیستر routeها در `services/core/route/`.
9. اجرای `gofmt` بعد از هر تغییر کد Go:
   - بررسی: `find . -name '*.go' -not -path './pkg/.cache/*' -print0 | xargs -0 gofmt -l`
   - اعمال: `find . -name '*.go' -not -path './pkg/.cache/*' -print0 | xargs -0 gofmt -w`
10. اجرای `go mod tidy` و سپس `go test ./...`.

## اجرا (Memory fallback)

```bash
go run ./services/core
```

## اجرا (MySQL mode)

1. تنظیم env:
```bash
export MYSQL_DSN="user:pass@tcp(127.0.0.1:3306)/sample?multiStatements=true"
export MYSQL_TASK_TABLE="tasks"
```
2. ساخت جدول:
```bash
go run ./commands/userMigration --table tasks
```
3. اجرای API:
```bash
go run ./services/core
```

## تست دستی endpointها

```bash
curl -X POST http://localhost:8080/task/create \
  -H 'Content-Type: application/json' \
  -d '{"body":{"title":"Write docs","description":"draft"}}'
```

```bash
curl "http://localhost:8080/task/list?page=1&perPage=10"
```
