# شرح دقیق فایل‌ها و توابع پروژه

این سند، فایل‌هایی را که در این پروژه ایجاد/ویرایش شده‌اند با جزئیات داخلی توضیح می‌دهد؛ یعنی فقط نام فایل نیست، بلکه type/function/const و رفتار دقیق هرکدام هم آمده است.

## services/core/main.go
- `main()`:
- یک `fiber.New()` می‌سازد.
- با `route.SetupRoutes(app)` تمام routeها را رجیستر می‌کند و نقشه route را می‌گیرد.
- با `fmt.Println` و `fmt.Printf` پیام اجرا و routeها را چاپ می‌کند.
- با `app.Listen(":8080")` سرور را بالا می‌آورد.
- در صورت خطا، با `log.Fatal` برنامه را متوقف می‌کند.
- importهای blank زیر فقط برای اجرای `init()` پکیج‌ها هستند:
- `_ ".../models/task"`
- `_ ".../models/user"`
- هدف این است که `repositories.TaskRepo` و `repositories.UserRepo` قبل از هندل درخواست مقداردهی شوند.

## services/core/route/route.go
- `SetupRoutes(app *fiber.App) map[string]string`:
- routeهای task و user را با هم merge می‌کند.
- خروجی، یک map تجمیعی از نام route به path است.
- `mergeMaps(maps ...map[string]string) map[string]string`:
- mapهای ورودی را به یک map جدید تبدیل می‌کند.
- اگر key تکراری باشد، مقدار map آخر overwrite می‌شود.

## services/core/route/task_route.go
- `taskRoutes`:
- map نام‌گذاری routeهای task.
- شامل `taskCreate: /task/create` و `taskList: /task/list`.
- `SetupTaskRoute(app *fiber.App) map[string]string`:
- `Create` را روی `POST /task/create` ثبت می‌کند.
- `List` را روی `POST /task/list` ثبت می‌کند.
- همان `taskRoutes` را برمی‌گرداند.

## services/core/route/user_route.go
- `userRoutes`:
- map نام‌گذاری routeهای user.
- شامل `userCreate: /user/create` و `userInfo: /user/info`.
- `SetupUserRoute(app *fiber.App) map[string]string`:
- `Create` را روی `POST /user/create` ثبت می‌کند.
- `Info` را روی `POST /user/info` ثبت می‌کند.
- همان `userRoutes` را برمی‌گرداند.

## controllers/mainController/main.go
- `type errorResponse`:
- قالب پاسخ خطا: `errorCode` و `message`.
- `type responseEnvelope`:
- قالب پاسخ موفق: `data`.
- `InitAPI(ctx *fiber.Ctx, sectionErrCode string) context.Context`:
- فعلاً tracing/context واقعی ندارد.
- ورودی‌ها را فقط discard می‌کند و `context.Background()` می‌دهد.
- `FinishAPISpan(ctx *fiber.Ctx)`:
- placeholder برای پایان span.
- فعلاً no-op است.
- `ParseBody(ctx *fiber.Ctx, req any) (string, int, error)`:
- body را با `ctx.BodyParser(req)` parse می‌کند.
- اگر parse خطا بدهد: `("01", 400, err)`.
- سپس `fillHeaders` را صدا می‌زند تا headers به فیلد `Headers` در request تزریق شود.
- سپس `validateBody` را اجرا می‌کند.
- موفقیت: `( "", 200, nil )`.
- `Error(ctx *fiber.Ctx, baseErrCode, section, errStr string, code int, err error) error`:
- کد خطای نهایی را با فرمت `baseErrCode + section + errStr` می‌سازد.
- پاسخ JSON خطا را با status ارسالی برمی‌گرداند.
- `Response(ctx *fiber.Ctx, res any) error`:
- پاسخ موفق را با status=200 و envelope استاندارد برمی‌گرداند.
- `GetUser(ctx *fiber.Ctx) userDataModel.User`:
- فعلاً کاربر دمو برمی‌گرداند: `ID=11, Username=demo-user`.
- `fillHeaders(ctx *fiber.Ctx, req any)`:
- با reflection اگر req pointer به struct باشد و فیلد `Headers` از نوع map قابل set باشد، headerهای HTTP را در آن قرار می‌دهد.
- `validateBody(req any) (string, int, error)`:
- با reflection دنبال فیلد `Body` می‌گردد.
- اگر `Body` متدی با امضای `Validate(validateExtraData commonSchema.ValidateExtraData)` داشته باشد آن را اجرا می‌کند.
- headerهای استخراج‌شده را داخل `ValidateExtraData` می‌فرستد.

## controllers/task/create.go
- `Create(ctx *fiber.Ctx) error`:
- `InitAPI(..., "11")` را صدا می‌زند.
- request را به شکل `BaseRequest[taskSchema.CreateRequest]` می‌سازد.
- `ParseBody` برای parse+validate اجرا می‌شود.
- خطای parse/validate با `TaskErrCode + section 01` برمی‌گردد.
- `repositories.TaskRepo.Create(...)` اجرا می‌شود.
- خطای repository با `TaskErrCode + section 02` برمی‌گردد.
- موفقیت با `mainController.Response`.

## controllers/task/list.go
- `List(ctx *fiber.Ctx) error`:
- `InitAPI(..., "12")`.
- request: `BaseRequest[taskSchema.ListRequest]`.
- parse/validate مانند create.
- repository call: `repositories.TaskRepo.List(...)`.
- الگوی خطا/موفقیت دقیقاً مشابه create.

## controllers/user/create.go
- `Create(ctx *fiber.Ctx) error`:
- `InitAPI(..., "21")`.
- request: `BaseRequest[userSchema.CreateRequest]`.
- parse/validate و سپس `repositories.UserRepo.Create(...)`.
- خطاها با base `UserErrCode` و sectionهای `01/02`.

## controllers/user/info.go
- `Info(ctx *fiber.Ctx) error`:
- `InitAPI(..., "22")`.
- request: `BaseRequest[userSchema.InfoRequest]`.
- parse/validate و سپس `repositories.UserRepo.Info(...)`.
- خطاها با base `UserErrCode` و sectionهای `01/02`.

## apiSchema/commonSchema/base.go
- `type BaseRequest[T any]`:
- فیلد `Body T` برای payload اصلی.
- فیلد `Headers map[string]string` برای headerهای تزریق‌شده توسط کنترلر.
- `type ValidateExtraData`:
- داده کمکی برای validateها.
- فعلاً فقط `Headers` دارد.

## apiSchema/taskSchema/request.go
- `type CreateRequest`:
- `Title`, `Description`.
- tagهای `json/msgpack/validate` تنظیم شده‌اند.
- `type ListRequest`:
- `Page`, `PerPage`.

## apiSchema/taskSchema/validate.go
- `(req *CreateRequest) Validate(...)`:
- `Title` را trim می‌کند.
- اگر خالی باشد: `errStr="03"`, `code=400`, `err=customErr.InvalidTitle`.
- در غیر این صورت موفق.
- `(req *ListRequest) Validate(...)`:
- اگر `Page < 1`: `errStr="06"`, `400`, `InvalidPage`.
- اگر `PerPage < 1 || PerPage > 100`: `errStr="09"`, `400`, `InvalidPerPage`.
- در غیر این صورت موفق.

## apiSchema/taskSchema/response.go
- `type CreateResponse`:
- خروجی create شامل `Task`.
- `type ListResponse`:
- خروجی list شامل `Tasks`, `Page`, `PerPage`, `Total`.

## apiSchema/userSchema/request.go
- `type CreateRequest`:
- `Username`, `Email`.
- `type InfoRequest`:
- `UserID`.

## apiSchema/userSchema/validate.go
- `(req *CreateRequest) Validate(...)`:
- `Username` و `Email` را trim می‌کند.
- username خالی: `errStr="03"`, `400`, `InvalidUsername`.
- email خالی: `errStr="06"`, `400`, `InvalidEmail`.
- `(req *InfoRequest) Validate(...)`:
- اگر `UserID < 1`: `errStr="09"`, `400`, `InvalidUserID`.

## apiSchema/userSchema/response.go
- `type CreateResponse`:
- خروجی create شامل `User`.
- `type InfoResponse`:
- خروجی info شامل `User`.

## models/repositories/taskRepo.go
- `type TaskRepository interface`:
- قرارداد متدهای `Create` و `List` برای task.
- خروجی هر متد: `(res, errStr, code, err)`.
- `var TaskRepo TaskRepository`:
- متغیر سراسری برای inject پیاده‌سازی.
- در `init()` پکیج task مقداردهی می‌شود.

## models/repositories/userRepo.go
- `type UserRepository interface`:
- قرارداد متدهای `Create` و `Info` برای user.
- `var UserRepo UserRepository`:
- مرجع سراسری پیاده‌سازی user repository.

## models/task/datamodel/task.go
- `type Task`:
- فیلدها: `ID`, `Title`, `Description`, `CreatedAt`.
- مدل پایه ذخیره‌سازی/پاسخ برای task.

## models/task/repository_create.go
- `type Repository`:
- `idCounter int64`: شمارنده id.
- `tasks []Task`: ذخیره in-memory.
- `lock sync.RWMutex`: قفل روی tasks.
- `listCache map[string]taskSchema.ListResponse`: کش list.
- `cacheLock sync.RWMutex`: قفل روی کش.
- `var once, repoIns`:
- الگوی singleton با `sync.Once`.
- `GetRepo() *Repository`:
- نمونه singleton را می‌سازد/برمی‌گرداند.
- مقدار اولیه: `idCounter=100`, `tasks=[]`, `listCache={}`.
- `init()`:
- `repositories.TaskRepo = GetRepo()`.
- `nextID() int64`:
- با `atomic.AddInt64` شناسه یکتا می‌سازد.
- `Create(...)`:
- task جدید از request می‌سازد.
- `CreatedAt` با `time.Now().UTC().Format(time.RFC3339)`.
- داخل lock به `tasks` append می‌کند.
- کل `listCache` را invalidate می‌کند (reset به map خالی).
- `taskSchema.CreateResponse` برمی‌گرداند.

## models/task/repository_list.go
- `List(...)`:
- cacheKey از page/perPage می‌سازد: `task:list:page:%d:perPage:%d`.
- ابتدا کش را با read lock چک می‌کند.
- در cache-hit: clone پاسخ کش‌شده را برمی‌گرداند.
- در cache-miss: از `tasks` کپی می‌گیرد.
- pagination با start/end محاسبه می‌شود.
- `ListResponse` ساخته می‌شود.
- پاسخ در کش ذخیره می‌شود (با clone).
- خروجی استاندارد `(res, "", 200, nil)`.
- `cloneListResponse(source)`:
- کپی امن از struct و slice `Tasks` می‌سازد تا mutation بیرونی روی cache اثر نگذارد.

## models/task/repository_cache_test.go
- `TestListCacheAndInvalidation`:
- state repository و cache را ریست می‌کند.
- یک task create می‌کند.
- اولین `List` را می‌زند و انتظار دارد cache پر شود.
- دومین `List` را می‌زند و تطابق تعداد taskها را چک می‌کند.
- create دوم انجام می‌دهد.
- انتظار دارد cache کامل invalidate شده باشد (`len(listCache)==0`).

## models/user/datamodel/user.go
- `type User`:
- فیلدها: `ID`, `Username`, `Email`.
- نکته فنی: نام پکیج این فایل `dataModel` است (حرف M بزرگ).

## models/user/repository_create.go
- `type Repository`:
- `idCounter`, `users`, `lock`.
- singleton با `once` و `repoIns`.
- `GetRepo()`:
- نمونه singleton user repo با مقدار اولیه `idCounter=10` و users خالی.
- `init()`:
- `repositories.UserRepo = GetRepo()`.
- `nextID()`:
- تولید id جدید با atomic.
- `Create(...)`:
- user جدید می‌سازد.
- با write lock به `users` اضافه می‌کند.
- `CreateResponse` را با status OK برمی‌گرداند.

## models/user/repository_info.go
- `Info(...)`:
- با read lock لیست users را می‌خواند.
- با `req.Body.UserID` کاربر را جستجو می‌کند.
- اگر پیدا شود `InfoResponse{User: ...}` برمی‌گرداند.
- اگر پیدا نشود: `errStr="12"`, `400`, `customErr.UserNotFound`.

## statics/constants/controllerBaseErrCode/base.go
- `UserErrCode = "2001"`.
- `TaskErrCode = "2007"`.
- این‌ها base کد خطا در controller هستند.

## statics/constants/error_message.go
- string constantهای پیام خطا:
- `InvalidTitle`, `InvalidPage`, `InvalidPerPage`, `InvalidUsername`, `InvalidEmail`, `InvalidUserID`, `UserNotFound`.

## statics/constants/status/status.go
- کدهای status:
- `200`, `400`, `401`, `403`, `500`, `429`.

## statics/customErr/err.go
- متغیرهای error با `errors.New(...)` بر اساس constantها:
- `InvalidTitle`, `InvalidPage`, `InvalidPerPage`, `InvalidUsername`, `InvalidEmail`, `InvalidUserID`, `UserNotFound`.

## tests/task_tests/task_create_test.go
- `TestCreateTask`:
- app Fiber تستی بالا می‌آورد.
- routeها را رجیستر می‌کند.
- payload create task می‌سازد.
- request `POST /task/create` می‌زند.
- status code را 200 انتظار دارد.

## tests/task_tests/task_list_test.go
- `TestListTask`:
- ابتدا یک task ایجاد می‌کند.
- سپس `POST /task/list` با `page=1,perPage=10` می‌زند.
- status code را 200 انتظار دارد.

## tests/user_tests/user_create_test.go
- `TestCreateUser`:
- payload create user می‌سازد.
- `POST /user/create` می‌زند.
- status code را 200 انتظار دارد.

## tests/user_tests/user_info_test.go
- `TestInfoUser`:
- اول user می‌سازد (`POST /user/create`).
- body پاسخ create را parse می‌کند.
- `data.user.id` را استخراج می‌کند.
- با همان id درخواست `POST /user/info` می‌زند.
- status code را 200 انتظار دارد.

## commands/README.md
- کاربرد پوشه `commands` را توضیح می‌دهد.
- سه اسکریپت CLI موجود را لیست می‌کند.
- نمونه اجرای هر اسکریپت را نشان می‌دهد.

## commands/elasticsearch_reindex/main.go
- `main()`:
- flagها: `--index`, `--batch`.
- لاگ start/done چاپ می‌کند.
- `TODO` دارد و فعلاً فقط `time.Sleep(100ms)` برای شبیه‌سازی اجرا.

## commands/user_migration/main.go
- `main()`:
- flagها: `--from`, `--to`, `--limit`.
- لاگ start/done چاپ می‌کند.
- منطق مهاجرت واقعی هنوز `TODO` است.

## commands/stats_update/main.go
- `main()`:
- flagها: `--period`, `--dry-run`.
- لاگ start/done چاپ می‌کند.
- منطق واقعی آپدیت آمار هنوز `TODO` است.

## models/task/dataSources/README.md
- نقش پوشه data source برای task را مشخص می‌کند.
- اعلام می‌کند پیاده‌سازی datasourceهای task باید اینجا قرار بگیرد.

## models/user/dataSources/README.md
- مشابه task datasource؛ برای user.

## middleware/README.md
- پوشه middleware برای auth و middlewareهای سراسری رزرو شده است.

## template/README.md
- پوشه template برای الگوهای feature تعریف شده است.

## pkg/README.md
- پوشه pkg برای پکیج‌های مشترک مثل pagination/cache/helper تعریف شده است.

## PROJECT_EXPLANATION.md
- مستند معماری و جریان پروژه به زبان انگلیسی.
- شامل: overview، tech stack، lifecycle، caching، testing، limitations و next steps.
- در این فایل، توضیح سطح معماری آمده؛ اما سند فعلی (`PROJECT_FILES_DETAILED_FA.md`) سطح فایل/تابع را پوشش می‌دهد.

## وضعیت فعلی اجرا
- تست‌ها با `go test ./...` پاس شدند.
- endpoint `task/list` در این پروژه با `POST` پیاده‌سازی شده است.
- ذخیره‌سازی فعلی in-memory است و data import به جدول/DB واقعی هنوز انجام نشده است.
