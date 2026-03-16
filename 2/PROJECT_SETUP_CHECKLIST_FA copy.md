# شرح کامل پروژه پوشه 2 (Virasty-Style)

## هدف این سند
این سند برای پوشه 2 تهیه شده و هدفش ارائه توضیح دقیق درباره هدف هر پوشه، هر فایل و هر تابع است. تمرکز روی ساختار لایه‌ای Virasty و جریان کامل درخواست‌هاست.

## نمای کلی جریان درخواست
1. درخواست HTTP به Fiber می‌رسد (`services/core/main.go`).
2. مسیرها در routeها رجیستر می‌شوند (`services/core/route/*`).
3. کنترلر متناظر صدا زده می‌شود (`controllers/*`).
4. کنترلر `ParseBody` یا `ParseQuery` را اجرا می‌کند (`controllers/mainController/main.go`).
5. validation روی schema انجام می‌شود (`apiSchema/*/validate.go`).
6. repo متد را اجرا می‌کند (`models/*/repository*.go`).
7. repo از datasource (memory یا MySQL) استفاده می‌کند (`models/*/dataSources/*`).
8. پاسخ استاندارد برگردانده می‌شود (`controllers/mainController/Response`).

## ترتیب اجرای پوشه‌ها به حالت استک
بالای استک عمیق‌ترین لایه اجرا است و پایین استک نقطه ورود درخواست است. پاسخ در مسیر برگشت همین استک بالا می‌آید.
1. `2/models/*/dataSources/*`
2. `2/models/*`
3. `2/models/repositories`
4. `2/apiSchema/*`
5. `2/controllers/mainController`
6. `2/controllers/*`
7. `2/services/core/route`
8. `2/services/core`

## ترتیب اجرای فایل‌ها (از شروع برنامه تا پاسخ)
این ترتیب نشان می‌دهد برنامه از کدام فایل شروع می‌شود و سپس چه فایل‌هایی اجرا می‌شوند.
1. شروع برنامه: `2/services/core/main.go` (تابع `main`)
2. ساخت Fiber و تنظیمات اولیه: `2/services/core/main.go`
3. رجیستر مسیرها: `2/services/core/route/route.go`
4. ثبت مسیرهای task: `2/services/core/route/taskRoute.go`
5. ثبت مسیرهای user: `2/services/core/route/userRoute.go`
6. ورود درخواست به کنترلر: `2/controllers/task/*.go` یا `2/controllers/user/*.go`
7. پارس ورودی و آماده‌سازی: `2/controllers/mainController/main.go` (`ParseBody` یا `ParseQuery`)
8. اعتبارسنجی داده‌ها: `2/apiSchema/*/validate.go`
9. اجرای منطق دامنه (repo): `2/models/*/repository*.go`
10. دسترسی به دیتاسورس: `2/models/*/dataSources/*`
11. ساخت پاسخ استاندارد: `2/controllers/mainController/main.go` (`Response` یا `Error`)

## قبل از اجرای برنامه چه کار کنیم؟
بعد از تکمیل پوشه‌ها و فایل‌ها، قبل از اجرای برنامه این مراحل را انجام بده:
1. کدها را فرمت کن: `gofmt -w .`
2. وابستگی‌ها را نصب/به‌روز کن: `go mod tidy`
3. اگر از MySQL استفاده می‌کنی، envهای لازم را ست کن (مثل `MYSQL_DSN`) و مطمئن شو دسترسی DB برقرار است.
4. اگر فقط memory استفاده می‌کنی، نیازی به تنظیم DB نیست.
5. در صورت نیاز، تست‌ها را اجرا کن: `go test ./...`
6. برنامه را اجرا کن: `go run ./services/core`

## پوشه‌ها و اهداف
1. `2/apiSchema`
هدف: تعریف قراردادهای API (request/response/validate) و ساختارهای مشترک.
2. `2/apiSchema/commonSchema`
هدف: ساختارهای پایه مثل `BaseRequest` و داده‌های کمکی validate.
3. `2/apiSchema/taskSchema`
هدف: قراردادهای مربوط به دامنه task.
4. `2/apiSchema/userSchema`
هدف: قراردادهای مربوط به دامنه user.
5. `2/controllers`
هدف: هندلرهای HTTP برای هر دامنه.
6. `2/controllers/mainController`
هدف: هسته مشترک کنترلرها (پارسیگ، خطا، پاسخ).
7. `2/controllers/task`
هدف: کنترلرهای task (create/list/update/delete).
8. `2/controllers/user`
هدف: کنترلرهای user (create/info).
9. `2/models`
هدف: لایه domain و repository.
10. `2/models/repositories`
هدف: interfaceهای repository و متغیرهای global.
11. `2/models/task`
هدف: repository دامنه task و اتصال datasourceها.
12. `2/models/task/dataModel`
هدف: مدل داده Task.
13. `2/models/task/dataSources`
هدف: interfaceهای datasource و پیاده‌سازی‌ها.
14. `2/models/task/dataSources/memoryDS`
هدف: دیتاسورس حافظه (DB و cache) برای تست و حالت بدون MySQL.
15. `2/models/task/dataSources/mysqlDS`
هدف: دیتاسورس MySQL + ساخت جدول و تنظیمات اتصال.
16. `2/models/user`
هدف: repository ساده دامنه user (در حافظه).
17. `2/models/user/dataModel`
هدف: مدل داده User.
18. `2/services`
هدف: سرویس‌های قابل اجرا.
19. `2/services/core`
هدف: سرویس اصلی API.
20. `2/services/core/route`
هدف: تعریف و رجیستر مسیرها.
21. `2/statics`
هدف: ثابت‌ها و خطاهای پروژه.
22. `2/statics/constants`
هدف: تعریف کلیدهای خطا و status codeها.
23. `2/statics/constants/controllerBaseErrCode`
هدف: کد پایه خطا برای دامنه‌ها.
24. `2/statics/constants/status`
هدف: نگهداری HTTP status codeها.
25. `2/statics/customErr`
هدف: تعریف errorهای قابل استفاده از روی ثابت‌ها.
26. `2/tests`
هدف: تست‌های integration برای API.
27. `2/tests/task_tests`
هدف: تست‌های دامنه task.
28. `2/tests/user_tests`
هدف: تست‌های دامنه user.
29. `2/commands`
هدف: اسکریپت‌های CLI و migration نمونه.
30. `2/middleware`
هدف: محل نگهداری middlewareهای پروژه (فعلا خالی).
31. `2/template`
هدف: الگوهای توسعه feature جدید.
32. `2/pkg`
هدف: کتابخانه‌های اشتراکی (فعلا خالی).
33. `2/pkg/.cache`
هدف: کش تولیدی Go (فایل‌های تولیدی، جزو سورس نیست).

## فایل‌ها و توابع

### فایل‌های ریشه پوشه 2
1. `2/README.md`
هدف: معرفی پروژه، لیست endpointها، نحوه اجرا و نمونه درخواست‌ها.
2. `2/VIRASTY_CODING_STYLE_GUIDE.md`
هدف: استاندارد کدنویسی و ساختار پروژه به سبک Virasty.
3. `2/Virasty_DOCUMENTATION_FA.md`
هدف: مستندات فارسی پروژه و سبک Virasty.
4. `2/استاندارد نام_گذاری و کدنویسی در پروژه.pdf`
هدف: مرجع PDF برای استاندارد نام‌گذاری و کدنویسی.
5. `2/go.mod`
هدف: تعریف ماژول Go و وابستگی‌ها.
6. `2/go.sum`
هدف: قفل نسخه وابستگی‌ها.
7. `2/.DS_Store`
هدف: فایل سیستمی macOS (قابل حذف، جزو سورس نیست).

### apiSchema/commonSchema
1. `2/apiSchema/commonSchema/base.go`
هدف: تعریف envelope استاندارد درخواست و داده‌های اضافی validate.
ساختار `BaseRequest[T]`: بدنه درخواست و هدرها را استاندارد می‌کند.
ساختار `ValidateExtraData`: ورودی کمکی برای validateها (فعلا هدرها).
توابع: ندارد.

### apiSchema/taskSchema
1. `2/apiSchema/taskSchema/request.go`
هدف: تعریف ساختارهای درخواست task.
`CreateRequest`: عنوان و توضیحات تسک.
`ListRequest`: صفحه‌بندی لیست تسک‌ها.
`UpdateRequest`: آپدیت جزیی تسک (title/description اختیاری).
`DeleteRequest`: حذف نرم بر اساس taskID.
توابع: ندارد.

2. `2/apiSchema/taskSchema/response.go`
هدف: تعریف ساختارهای پاسخ task.
`CreateResponse`: خروجی create (Task).
`ListResponse`: خروجی list (tasks + paging info).
`UpdateResponse`: خروجی update (Task).
`DeleteResponse`: خروجی delete (Task حذف‌شده).
توابع: ندارد.

3. `2/apiSchema/taskSchema/validate.go`
هدف: اعتبارسنجی درخواست‌های task.
تابع `CreateRequest.Validate`: trim کردن title و رد title خالی.
تابع `ListRequest.Validate`: کنترل page>=1 و perPage بین 1 تا 100.
تابع `UpdateRequest.Validate`: کنترل taskID>=1 و الزام حداقل یکی از title/description.
تابع `DeleteRequest.Validate`: کنترل taskID>=1.

### apiSchema/userSchema
1. `2/apiSchema/userSchema/request.go`
هدف: تعریف ساختارهای درخواست user.
`CreateRequest`: username و email.
`InfoRequest`: userID.
توابع: ندارد.

2. `2/apiSchema/userSchema/response.go`
هدف: تعریف ساختارهای پاسخ user.
`CreateResponse`: خروجی create (User).
`InfoResponse`: خروجی info (User).
توابع: ندارد.

3. `2/apiSchema/userSchema/validate.go`
هدف: اعتبارسنجی درخواست‌های user.
تابع `CreateRequest.Validate`: trim کردن username/email و رد مقدار خالی.
تابع `InfoRequest.Validate`: کنترل userID>=1.

### controllers/mainController
1. `2/controllers/mainController/main.go`
هدف: توابع مشترک کنترلرها برای parsing, validation و پاسخ‌دهی.
تابع `InitAPI`: آماده‌سازی context برای هر درخواست (فعلا stub).
تابع `FinishAPISpan`: بستن span (فعلا stub).
تابع `ParseBody`: BodyParser + پرکردن هدرها + validate.
تابع `ParseQuery`: QueryParser + validate.
تابع `Error`: تولید پاسخ خطا با errorCode یکتا.
تابع `Response`: تولید envelope استاندارد پاسخ.
تابع `fillHeaders`: کپی کردن هدرها داخل `BaseRequest.Headers`.
تابع `validateBody`: اجرای متد Validate روی Body در `BaseRequest`.

### controllers/task
1. `2/controllers/task/create.go`
هدف: هندلر ساخت task.
تابع `Create`: ParseBody → Repo.Create → Response/Error.

2. `2/controllers/task/list.go`
هدف: هندلر لیست task.
تابع `List`: ParseQuery → Repo.List → Response/Error.

3. `2/controllers/task/update.go`
هدف: هندلر آپدیت task.
تابع `Update`: ParseBody → Repo.Update → Response/Error.

4. `2/controllers/task/delete.go`
هدف: هندلر حذف نرم task.
تابع `Delete`: ParseBody → Repo.Delete → Response/Error.

### controllers/user
1. `2/controllers/user/create.go`
هدف: هندلر ساخت user.
تابع `Create`: ParseBody → Repo.Create → Response/Error.

2. `2/controllers/user/info.go`
هدف: هندلر دریافت اطلاعات user.
تابع `Info`: ParseBody → Repo.Info → Response/Error.

### models/repositories
1. `2/models/repositories/taskRepo.go`
هدف: interface دامنه task و نگهداری متغیر global.
`TaskRepository`: قرارداد Create/List/Update/Delete.
`TaskRepo`: متغیر global برای wiring.
توابع: ندارد.

2. `2/models/repositories/userRepo.go`
هدف: interface دامنه user و نگهداری متغیر global.
`UserRepository`: قرارداد Create/Info.
`UserRepo`: متغیر global برای wiring.
توابع: ندارد.

### models/task/dataModel
1. `2/models/task/dataModel/task.go`
هدف: تعریف مدل داده task برای پاسخ‌ها و datasourceها.
`Task`: شامل ID/Title/Description/CreatedAt/UpdatedAt/DeletedAt با tagهای json/msgpack/gorm.
توابع: ندارد.

### models/user/dataModel
1. `2/models/user/dataModel/user.go`
هدف: تعریف مدل داده user.
`User`: شامل ID/Username/Email با tagهای json/msgpack.
توابع: ندارد.

### models/task/dataSources
1. `2/models/task/dataSources/taskDS.go`
هدف: تعریف interfaceهای دیتاسورس task.
`TaskDBDS`: قرارداد CRUD و list برای task.
`TaskCacheDS`: قرارداد cache لیست.
توابع: ندارد.

2. `2/models/task/dataSources/README.md`
هدف: مستند کوتاه درباره قراردادها و پیاده‌سازی‌ها.

### models/task/dataSources/memoryDS
1. `2/models/task/dataSources/memoryDS/taskDBDS.go`
هدف: دیتاسورس حافظه برای task (بدون DB).
تابع `tehranLocation`: گرفتن location تهران (fallback ثابت).
تابع `NewTaskDBDS`: ساخت نمونه دیتاسورس با startID.
تابع `CreateTask`: ساخت task جدید در حافظه با زمان تهران.
تابع `ListTasks`: لیست taskها با pagination و فیلتر soft delete.
تابع `UpdateTask`: آپدیت title/description و set کردن UpdatedAt.
تابع `SoftDeleteTask`: حذف نرم و set کردن DeletedAt/UpdatedAt.
تابع `Reset`: ریست کامل دیتاسورس (برای تست).

2. `2/models/task/dataSources/memoryDS/taskCacheDS.go`
هدف: کش لیست taskها در حافظه.
تابع `NewTaskCacheDS`: ساخت cache خالی.
تابع `GetList`: دریافت cache براساس کلید.
تابع `SetList`: ذخیره پاسخ لیست در cache.
تابع `InvalidateList`: پاک‌کردن کل cache لیست.
تابع `Reset`: ریست cache (alias برای InvalidateList).

### models/task/dataSources/mysqlDS
1. `2/models/task/dataSources/mysqlDS/config.go`
هدف: خواندن تنظیمات MySQL از env و نرمال‌سازی DSN.
تابع `LoadConfigFromEnv`: خواندن `MYSQL_DSN` و تنظیمات pool و table.
تابع `readEnvInt`: خواندن عدد از env با مقدار پیش‌فرض.
تابع `normalizeDSN`: افزودن پارامترهای parseTime/loc/time_zone/charset.

2. `2/models/task/dataSources/mysqlDS/types.go`
هدف: تعریف interface عمومی برای اجرا روی `sql.DB` یا mock.
`DBExecutor`: قرارداد ExecContext/QueryContext/QueryRowContext/Close.
توابع: ندارد.

3. `2/models/task/dataSources/mysqlDS/connection.go`
هدف: باز کردن و تنظیم connection pool به MySQL.
تابع `Open`: ساخت `sql.DB`، تنظیم pool و Ping اولیه.

4. `2/models/task/dataSources/mysqlDS/schema.go`
هدف: اعتبارسنجی نام جدول و ساخت/ارتقاء جدول task.
تابع `ValidateTableName`: بررسی ایمن بودن نام جدول.
تابع `TaskTableIdentifier`: ساخت نام جدول با backtick.
تابع `EnsureTaskTable`: ساخت جدول و اجرای migrationهای ساده (best-effort).

5. `2/models/task/dataSources/mysqlDS/schema.sql`
هدف: DDL جدول tasks برای migration دستی.

6. `2/models/task/dataSources/mysqlDS/taskDBDS.go`
هدف: پیاده‌سازی دیتاسورس task روی MySQL.
تابع `tehranLocation`: گرفتن location تهران.
تابع `isUnknownColumnErr`: تشخیص خطای ستون ناشناخته (برای backward compatibility).
تابع `NewTaskDBDSFromEnv`: ساخت datasource از env و اطمینان از وجود جدول.
تابع `CreateTask`: insert و خواندن رکورد ایجادشده.
تابع `ListTasks`: لیست با pagination و فیلتر soft delete، سازگار با جدول‌های قدیمی.
تابع `UpdateTask`: آپدیت partial با updated_at و جلوگیری از آپدیت حذف‌شده‌ها.
تابع `SoftDeleteTask`: set deleted_at/updated_at برای حذف نرم.
تابع `joinCSV`: helper برای join کردن بخش‌های SET در SQL.
تابع `readTaskByID`: خواندن task و نگاشت زمان‌ها به تهران.
تابع `TableName`: خروجی نام جدول واقعی.

### models/task
1. `2/models/task/repository.go`
هدف: singleton repo و wiring datasourceها.
تابع `GetRepo`: ساخت یک‌باره repo و cache.
تابع `init`: ثبت repo در repositories.TaskRepo.
تابع `initializeDataSources`: انتخاب MySQL اگر فعال، در غیر این صورت memory.
تابع `db`: دسترسی به DB datasource.
تابع `cache`: دسترسی به cache datasource.

2. `2/models/task/repositoryCreate.go`
هدف: منطق ساخت task در repository.
تابع `Create`: کنترل خطاهای init، ساخت در DB و invalidation cache.

3. `2/models/task/repositoryList.go`
هدف: منطق لیست task با cache.
تابع `List`: cache hit/miss و ساخت `ListResponse`.
تابع `cloneListResponse`: کپی امن پاسخ cache برای جلوگیری از mutation.

4. `2/models/task/repositoryUpdate.go`
هدف: منطق آپدیت task در repository.
تابع `Update`: اجرای UpdateTask و مدیریت TaskNotFound.

5. `2/models/task/repositoryDelete.go`
هدف: منطق حذف نرم task در repository.
تابع `Delete`: اجرای SoftDeleteTask و مدیریت TaskNotFound.

6. `2/models/task/repositoryCache_test.go`
هدف: تست یونیت رفتار cache و invalidation.
تابع `TestListCacheAndInvalidation`: بررسی cache hit و پاک شدن cache بعد از create/update/delete.

### models/user
1. `2/models/user/repository.go`
هدف: singleton repo ساده برای user (فقط حافظه).
تابع `GetRepo`: ساخت repo با شمارنده اولیه.
تابع `init`: ثبت repo در repositories.UserRepo.

2. `2/models/user/repositoryCreate.go`
هدف: منطق ساخت user در حافظه.
تابع `nextID`: تولید ID جدید با atomic.
تابع `Create`: ساخت user و ذخیره در آرایه محافظت‌شده با lock.

3. `2/models/user/repositoryInfo.go`
هدف: واکشی اطلاعات user از حافظه.
تابع `Info`: جستجوی user با ID و برگرداندن خطا در صورت نبودن.

### services/core
1. `2/services/core/main.go`
هدف: entrypoint سرویس API.
تابع `main`: ساخت Fiber، رجیستر routeها، و Listen روی :8080.

### services/core/route
1. `2/services/core/route/route.go`
هدف: تجمیع routeهای دامنه‌ها.
تابع `SetupRoutes`: merge کردن route-mapها و ثبت routeها.
تابع `mergeMaps`: ترکیب mapهای مسیر.

2. `2/services/core/route/taskRoute.go`
هدف: مسیرهای task.
متغیر `taskRoutes`: نگاشت نام route به path.
تابع `SetupTaskRoute`: رجیستر POST/GET برای task.

3. `2/services/core/route/userRoute.go`
هدف: مسیرهای user.
متغیر `userRoutes`: نگاشت نام route به path.
تابع `SetupUserRoute`: رجیستر POST برای user.

### statics/constants
1. `2/statics/constants/controllerBaseErrCode/base.go`
هدف: کد پایه خطا برای دامنه‌ها.
مقادیر: `TaskErrCode=2007` و `UserErrCode=2001`.

2. `2/statics/constants/status/status.go`
هدف: ثابت‌های HTTP status code.
مقادیر: 200, 400, 401, 403, 429, 500.

3. `2/statics/constants/errorMessage.go`
هدف: کلیدهای خطا به صورت string ثابت برای پیام خطا.

### statics/customErr
1. `2/statics/customErr/err.go`
هدف: تبدیل constantها به errorهای واقعی با `errors.New`.

### commands
1. `2/commands/README.md`
هدف: مستند اجرای اسکریپت‌های CLI.

2. `2/commands/elasticSearchReindex/main.go`
هدف: نمونه اسکریپت ری‌ایندکس Elasticsearch.
تابع `main`: پارس فلگ‌ها و اجرای نمونه (placeholder).

3. `2/commands/statsUpdate/main.go`
هدف: نمونه اسکریپت آپدیت آمار.
تابع `main`: پارس فلگ‌ها و اجرای نمونه (placeholder).

4. `2/commands/userMigration/main.go`
هدف: آماده‌سازی جدول task در MySQL.
تابع `main`: خواندن env/flag و اجرای EnsureTaskTable.

### tests/task_tests
1. `2/tests/task_tests/taskCreate_test.go`
هدف: تست endpoint ساخت task.
تابع `TestCreateTask`: ارسال POST و بررسی status.

2. `2/tests/task_tests/taskList_test.go`
هدف: تست endpoint لیست task.
تابع `TestListTask`: ساخت یک task و بررسی پاسخ لیست.

3. `2/tests/task_tests/taskUpdate_test.go`
هدف: تست endpoint آپدیت task.
تابع `TestUpdateTask`: ساخت task، گرفتن id و ارسال update.

4. `2/tests/task_tests/taskDelete_test.go`
هدف: تست endpoint حذف نرم task.
تابع `TestSoftDeleteTask`: ساخت task، حذف، و بررسی عدم نمایش در لیست.

### tests/user_tests
1. `2/tests/user_tests/userCreate_test.go`
هدف: تست endpoint ساخت user.
تابع `TestCreateUser`: ارسال POST و بررسی status.

2. `2/tests/user_tests/userInfo_test.go`
هدف: تست endpoint info user.
تابع `TestInfoUser`: ساخت user و سپس درخواست info.

### middleware
1. `2/middleware/README.md`
هدف: توضیح اینکه middlewareها در این پوشه قرار می‌گیرند.

### template
1. `2/template/README.md`
هدف: توضیح محل قرارگیری templateهای feature.

### pkg
1. `2/pkg/README.md`
هدف: توضیح محل قرارگیری packageهای اشتراکی.

## یادداشت درباره فایل‌های تولیدی
فایل‌های داخل `2/pkg/.cache` و سایر فایل‌های سیستمی مثل `.DS_Store` توسط ابزارها تولید می‌شوند و جزو سورس محسوب نمی‌شوند. این سند آن‌ها را به عنوان خروجی تولیدی معرفی می‌کند و وارد جزئیات تک‌تک فایل‌های کش نمی‌شود.
