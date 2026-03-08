# مستندات فنی پروژه شبکه اجتماعی (ویراستی / Microblog API-Golang)

این سند حاصل بررسی فایل‌به‌فایل پروژه است و شامل ساختار پروژه، بهبودهای پیشنهادی و باگ‌ها/مشکلات شناسایی‌شده می‌باشد.

---

## ۱. ساختار پروژه

### ۱.۱ نمای کلی

پروژه یک **API بک‌اند** برای شبکه اجتماعی (میکروبلاگ) است که با **Go 1.24** و فریم‌ورک **Fiber v2** نوشته شده و از معماری لایه‌ای (Controllers → Services/Repositories → Models) پیروی می‌کند.

- **ماژول:** `apiGolang`
- **ورودی اصلی سرویس Core:** `services/core/main.go` روی پورت **7575**
- **وابستگی داخلی:** کتابخانه خصوصی `code.ts.co.ir/gaplib/golib`

### ۱.۲ ساختار دایرکتوری‌ها

| مسیر | توضیح |
|------|--------|
| **apiSchema/** | اسکیماهای درخواست/پاسخ و اعتبارسنجی برای هر دامنه (user، post، conversation و غیره) |
| **commands/** | اسکریپت‌های CLI و مایگریشن (ElasticSearch، مایگریشن کاربران، آپدیت آمار و غیره) |
| **controllers/** | هندلرهای HTTP به ازای هر دامنه (user، post، notification، community و غیره) |
| **middleware/** | احراز هویت، AppScope، دسترسی community و OpenAPI |
| **models/** | مدل‌های داده، repositoryها و لایه دسترسی به داده (DB، Redis، Vitess، HTTP) |
| **pkg/** | کد مشترک: pagination، cache، imageProcessor، sharding، url utils |
| **services/** | سرویس‌های قابل اجرا: **core** (API اصلی)، **cron**، **storage** (TUS)، **openAPI**، **cloudMeeting**، **sitemap** |
| **statics/** | تنظیمات (appConfig برای هر اپ)، ثابت‌ها، ترجمه، خطاها، فونت و لوگو |
| **template/** | قالب برای ساخت سرویس/کنترلر/تست جدید |
| **tests/** | تست‌های API (integration/functional) |

### ۱.۳ جریان درخواست (Core API)

1. **Fiber App** در `services/core/main.go` ساخته می‌شود و `setupRoute(app)` صدا زده می‌شود.
2. **Middlewareها به ترتیب:**
   - `CheckServiceAvaiblaity` — در دسترس بودن سرویس
   - `AppScope` — اعتبار scope و tenant (و در صورت نیاز دسترسی community)
   - `AuthMiddleware` — احراز هویت با توکن/DUID و تعیین user/device
   - `CheckServiceMaintenanceMode` — حالت تعمیرات
   - تنظیم زبان از روی `Accept-Language`
3. **Routeها** از پکیج‌های داخل `services/core/route/` (مثلاً `userRoute`، `postRoute`) ثبت می‌شوند و هر endpoint به یک تابع در **controllers** متصل می‌شود.
4. در **Controller:**  
   `InitAPI` → `ParseBody` (دیکد بدنه، اعتبارسنجی، بررسی دسترسی) → فراخوانی **Repository** → `Response` یا `Error`؛ در پایان `FinishAPISpan`.

### ۱.۴ الگوی Controller

- هر endpoint معمولاً: `InitAPI` با کد خطای پایه، `ParseBody` برای request، فراخوانی یک یا چند repo، و در نهایت `Response` یا `Error`.
- کدهای خطا در `statics/constants/controllerBaseErrCode` تعریف شده‌اند (مثلاً `UserErrCode = "2001"`).
- پشتیبانی از **JSON** و **MessagePack** و در صورت نیاز **GCM encryption** با نمک گرفته‌شده از DUID.

### ۱.۵ مدل‌ها و Repositoryها

- هر دامنه در `models/` معمولاً شامل: **dataModel** (ساختارها)، **repository** (لاجیک دسترسی)، و **dataSources** (DB، Redis، Vitess، HTTP).
- Repositoryهای سراسری در `models/repositories/` تعریف و در init با استفاده از `initialRepositories` مقداردهی می‌شوند.
- از **GORM**، **Redis**، **Elasticsearch**، **MySQL/Vitess** و در مواردی MinIO/Ceph برای ذخیره فایل استفاده شده است.

### ۱.۶ سرویس‌های جداگانه

- **core** — API اصلی (پورت 7575)
- **cron** — جاب‌های زمان‌بندی‌شده
- **storage** — آپلود فایل (TUS)
- **openAPI** — API برای مصرف خارجی (مثلاً Farsnext)
- **cloudMeeting** — جلسات ویدیویی/صوتی
- **sitemap** — تولید sitemap

هر کدام `main.go` و `route.go` (و در صورت نیاز زیرمسیرهای route) دارند.

### ۱.۷ پیکربندی و اپ‌ها

- **appConfig** در `statics/configs/appConfig` برای هر اپ (virasty، fars، wplus، iSamad و غیره) تعریف شده است.
- **commonConfig** شامل تنظیمات مشترک (device، post، user، notification، storage و غیره) است.
- ثابت‌ها در `statics/constants` (وضعیت‌ها، هدرها، زبان‌ها، کد خطاها).

### ۱.۸ امنیت و احراز هویت

- **Auth:** توکن دستگاه (با/بدون فرمت جدید امضا)، لیست سفید مسیرهای بدون لاگین، مسیرهای مخصوص پنل ادمین با DUID و توکن ادمین، مسیرهای استاتیک API Key.
- **Rate limiting:** در مسیرهای حساس (مثلاً لیست فالوور/فالوینگ) و **RequestLimiter** سراسری (با whitelist برای مسیرهای فقط-خواندی).
- **GCM:** برای رمزنگاری بدنه در صورت نیاز؛ نمک از ۱۶ کاراکتر اول DUID گرفته می‌شود.

---

## ۲. بهبودهای پیشنهادی

### ۲.۱ معماری و ساختار

1. **تفکیک بهتر لایه سرویس:** در حال حاضر کنترلرها مستقیماً repository صدا می‌زنند. اضافه کردن یک لایه **Service** (Business Logic) بین Controller و Repository خوانایی و تست‌پذیری را بالا می‌برد.
2. **Dependency Injection:** استفاده از DI (مثلاً با wire یا fx) به جای init و متغیرهای سراسری repository باعث تست‌پذیری و مدیریت وابستگی‌ها می‌شود.
3. **یکپارچگی مسیرها:** ثبت routeها در یک جا (مثلاً از روی تگ یا متادیتا) و کاهش تکرار در `route.go` و فایل‌های route دامنه‌ها.
4. **مستندسازی API:** استفاده از Swagger/OpenAPI به صورت متمرکز و به‌روز تا هم تیم و هم کلاینت‌ها از قرارداد API مطلع باشند.

### ۲.۲ کیفیت کد

1. **نام‌گذاری یکسان:** یکسان‌سازی نام فایل‌ها و پکیج‌ها (مثلاً حذف تایپوها و یکسان‌سازی پسوند Route).
2. **خطاها:** استفاده از **typed errors** و **error wrapping** (`fmt.Errorf("%w", err)`) به جای رشته ثابت تا لاگ و مانیتورینگ بهتر شود.
3. **Context:** اطمینان از پاس دادن `context` در همه لایه‌ها و استفاده از timeout/cancel در فراخوانی‌های خارجی.
4. **لاگ:** ساختار یکسان (مثلاً با zap) و حذف لاگ‌های دیباگ غیرضروری در production.

### ۲.۳ امنیت و عملیات

1. **حذف هاردکد:** هیچ شناسه کاربر، توکن یا کلید در کد نباشد؛ همه از config/env بیایند.
2. **محدودیت DUID:** قبل از استفاده از `DUID[:16]` برای GCM، طول DUID چک شود تا از panic جلوگیری شود.
3. **Rate limit و Backpressure:** برای endpointهای سنگین (فید، جستجو) محدودیت و backpressure روشن و قابل پیکربندی باشد.
4. **مستندسازی README:** به‌روزرسانی README با مسیر صحیح اجرا (`services/core` به جای `services/core` در متن فعلی) و اضافه کردن بخش env و اجرای تست.

### ۲.۴ تست و CI

1. **تست واحد:** برای لاجیک حساس (validation، permission، pagination) واحد تست اضافه شود.
2. **تست یکپارچگی:** تست‌های موجود در `tests/` با یک runner یکسان و گزارش پوشش اجرا شوند.
3. **Lint و فرمت:** استفاده از `golangci-lint` و `go fmt` در CI تا سبک کد یکدست بماند.

### ۲.۵ پایگاه داده و کش

1. **Connection pooling و timeout:** برای MySQL/Vitess و Redis تنظیمات pool و timeout به صورت صریح و مناسب محیط.
2. **حساسیت به replication lag:** در خواندن‌های بعد از نوشتن، در صورت نیاز خواندن از primary تا از consistency اطمینان حاصل شود.
3. **کلیدهای کش:** نام‌گذاری و TTL یکسان برای کلیدهای Redis تا از تداخل و نشت حافظه جلوگیری شود.

---

## ۳. باگ‌ها و مشکلات شناسایی‌شده

### ۳.۱ باگ‌های منطقی / احتمالی

| محل | شرح | پیشنهاد رفع |
|-----|------|-------------|
| **controllers/mainController/mainController.go** (حدود خط ۵۴۱) | در `CheckServiceMaintenanceMode` از `apiNameByRoute[ctx.OriginalURL()]` استفاده شده؛ در حالی که کلیدهای `apiNameByRoute` با **مسیر بدون query** (`ctx.Route().Path`) پر می‌شوند. با وجود query string، نام API خالی می‌ماند و ممکن است در حالت تعمیرات رفتار نادرست داشته باشد. | استفاده از `ctx.Route().Path` (یا مشابه آن) به جای `ctx.OriginalURL()` برای lookup. |
| **apiSchema/commonSchema/common.go** تابع `(ids *IDs) Set(values []int64)` (حدود خط ۱۰۵–۱۱۳) | اسلایس با `make([]ID, len(values))` ساخته می‌شود (با مقدار پیش‌فرض صفر) و بعد در حلقه فقط `append` می‌شود؛ در نتیجه خروجی شامل `len(values)` تا مقدار صفر در ابتدا و سپس مقادیر تکراری (به خاطر استفاده از همان متغیر `id` در حلقه) است. | ساخت اسلایس با طول صفر و ظرفیت `len(values)` و در حلقه مقداردهی درست هر عنصر؛ یا استفاده از اندیس و اختصاص مستقیم. |
| **controllers/mainController/mainController.go** تابع `GetGCMUserSalt` (حدود خط ۴۰۴–۴۰۶) | `return []byte(ctx.Get(headers.DUID)[:16])` در صورت خالی بودن یا کوتاه بودن DUID باعث **panic** (slice bounds out of range) می‌شود. | قبل از برش، بررسی `len(duid) >= 16` و در غیر این صورت برگرداندن خطا یا مقدار امن و عدم فعال کردن encryption. |
| **middleware/auth.go** (حدود خط ۴۱۹–۴۲۱) | شرط `auth.user.ID == 1757676064482597066 && routePath == "/log"` یک **کاربر خاص** را از Request Limiter معاف می‌کند. این یک backdoor/دیباگ است و در production نامناسب است. | حذف این شرط یا انتقال آن به config (مثلاً لیست userIDهای معاف) و غیرفعال بودن پیش‌فرض در production. |

### ۳.۲ تایپوها و نام‌گذاری (امکان باگ یا گیج‌کننده)

| محل | وضعیت فعلی | پیشنهاد |
|-----|------------|----------|
| **services/core/route.go** و **mainController** و **cloudMeeting/route.go** | `CheckServiceAvaiblaity` | `CheckServiceAvailability` |
| **statics/configs/appConfig/commonAppSchema/scopes.go** و **middleware/appScope.go** و **models/user/repository.go** | `NeedChekScopePermissionInstanceOfGlobalPermission` | `NeedCheckScopePermissionInstanceOfGlobalPermission` |
| **apiSchema/showcaseSchema/** | فایل `requset.go` | تغییر نام به `request.go` |
| **services/core/route/** | فایل `showcaseRuote.go` | تغییر نام به `showcaseRoute.go` |
| **statics/constants/status/statusCodes.go** | `UnAviableServiceError` | `UnavailableServiceError` |
| **statics/constants/constants.go** | `DefultThreadViewMode`, `DefultMediaViewMode` | `DefaultThreadViewMode`, `DefaultMediaViewMode` |
| **statics/constants/constants.go** | `ParticipantsSeprator` | `ParticipantsSeparator` |
| **statics/constants/constants.go** و ترجمه/خطاها | `ThisUserDoseNotHaveThisEduGroupLevel` | `ThisUserDoesNotHaveThisEduGroupLevel` |
| **statics/configs/commonConfig/otp/global.go** | `DefultTemplate` | `DefaultTemplate` |
| **models/showcaseContent/** و **showcaseBlockTab** | `SetDefultAutomaticConfigWhenConfigIsEmpty`, `DefultTab`, `HasDefultTab` | استفاده از `Default` به جای `Defult` |
| **pkg/pagination/pagination.go** | کلید ثابت `paginationEncryptionKey`: `"this!P!a!g!ition"` | در صورت هدف «pagination» اصلاح به `"this!P!a!g!ination"` (و ترجیحاً انتقال به config). |

### ۳.۳ نام فایل‌ها و پکیج‌ها (تایپو)

- **commands/fillPostsReplyDeleted/** — فایل `fiiPostsReplyDeleted.go` → پیشنهاد: `fillPostsReplyDeleted.go`
- **commands/fiilFeedPostsTable/** — پوشه و فایل با `fiil` و `fii` → پیشنهاد: `fillFeedPostsTable` و نام فایل متناسب
- **models/repositories/** — `posFeedRepo.go` → احتمالاً `postsFeedRepo.go`؛ `postsReletedRepo.go` → `postsRelatedRepo.go`؛ `postCampaignFollowupUserhistoryRepo.go` → `postCampaignFollowupUserHistoryRepo.go`؛ `roleSectionVisibiltyRepo.go` → `roleSectionVisibilityRepo.go`
- **models/deleteAndEditNewsVirastarDailyStat/** — `Virastar` → در صورت یکسان‌سازی با نام برند: `Virasty`

### ۳.۴ مشکلات کوچک دیگر

| محل | شرح |
|-----|------|
| **controllers/mainController/mainController.go** (حدود خط ۳۳۴–۳۳۵) | در لاگ «api started» دو بار `zap.String("scope", ...)` با `AppScope` و `AppScopeTenant` ثبت شده؛ کلید دوم بهتر است مثلاً `"scopeTenant"` باشد تا در لاگ خوانا باشد. |
| **README.md** | اشاره به `go1.19` و مسیر `cd /services/core` با اسلش مطلق؛ بهتر است مسیر نسبی `services/core` و نسخه Go مطابق `go.mod` (مثلاً 1.24) ذکر شود. |

---

## ۴. جمع‌بندی

- **ساختار:** پروژه با لایه‌های Controller، Model/Repository و Config/Constants مشخص است و چند سرویس جدا (core، cron، storage، openAPI، cloudMeeting، sitemap) دارد.
- **بهبودها:** تمرکز روی لایه سرویس، DI، حذف هاردکد، امنیت (GCM/DUID و معافیت کاربر خاص)، و تقویت تست و مستندسازی پیشنهاد می‌شود.
- **باگ‌ها:** چهار مورد باگ منطقی/امنیتی (Maintenance mode lookup، `IDs.Set`، `GetGCMUserSalt`، و معافیت کاربر از rate limit) و چندین تایپو در نام تابع/فایل/ثابت و یک باگ در `IDs.Set` در این سند فهرست شده‌اند.

با رفع موارد بالا، پایداری، امنیت و نگه‌داری پروژه بهبود خواهد یافت.
