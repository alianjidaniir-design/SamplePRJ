# راهنمای سبک کدنویسی و ساختار فنی ویراستی (Virasty Coding Style Guide)

این سند برای استفاده در **پروژه‌های دیگر** تهیه شده است. با دادن این فایل به توسعه‌دهنده یا AI، انتظار می‌رود کد دقیقاً مشابه **استایل و روش تیم فنی ویراستی** نوشته شود: ساختار پوشه‌ها، نام‌گذاری، الگوهای Controller/Repository/Schema، خطاها، و قراردادهای کدنویسی.

---

## ۱. ساختار کلی پروژه (لایه‌ها)

```
apiGolang/
├── apiSchema/           # اسکیماهای درخواست/پاسخ و اعتبارسنجی (هر دامنه یک پکیج)
├── commands/            # اسکریپت‌های CLI و مایگریشن
├── controllers/         # هندلرهای HTTP (هر دامنه یک پکیج، هر endpoint یک فایل)
├── middleware/          # Auth، AppScope، و غیره
├── models/              # داده و دسترسی به داده (هر entity یک پکیج با dataModel، repository، dataSources)
├── pkg/                 # کد مشترک: pagination، cache، imageProcessor، ...
├── services/            # سرویس‌های قابل اجرا (core، cron، storage، ...) هر کدام main + route
├── statics/             # configs، constants، translate، customErr
├── template/            # قالب برای feature جدید
└── tests/               # تست‌های API (هر دامنه یک پکیج _tests)
```

- **ماژول:** نام پروژه (مثلاً `apiGolang`)؛ در پروژه دیگر با نام همان پروژه عوض می‌شود.
- **ورودی اصلی API:** `services/core/main.go`؛ سرویس‌های دیگر هر کدام `main.go` و `route.go` در پوشه خودشان.

---

## ۲. قراردادهای نام‌گذاری

### ۲.۱ پکیج‌ها

- **حروف کوچک، یک کلمه:** `user`, `post`, `notification`, `route`, `dataModel`, `dataSources`.
- نام پکیج با نام آخرین بخش مسیر دایرکتوری هماهنگ است (مثلاً `controllers/user` → `package user`).

### ۲.۲ فایل‌ها

- **snake_case برای نام فایل:**  
  `repository_info.go`, `repository_follow.go`, `update_username.go`
- در این پروژه گاهی از **camelCase** برای فایل‌های مرتبط با یک عملیات استفاده شده:  
  `repositoryInfo.go`, `repositoryFollow.go` (در پروژه جدید می‌توان هر دو را قبول کرد؛ ترجیحاً یک سبک ثابت).
- پسوند تست: `*_test.go` و نام پکیج تست: `user_tests`, `post_tests` (با underscore).

### ۲.۳ تایپ‌ها و توابع

- **PascalCase** برای: نوع‌ها، توابع عمومی، ثابت‌ها، فیلدهای export شده.
- **camelCase** برای: متغیرها، پارامترها، فیلدهای خصوصی، کلیدهای map.
- **اختصارهای رایج:**  
  `ID` (نه Id)، `URL` (نه Url)، `SID`, `DUID`, `OTP`, `API`, `HTTP`, `DB`, `Repo`, `ctx` برای context، `req` برای request، `res` برای response.

### ۲.۴ مسیرهای API و کلید route

- **مسیر:** کلمه‌های کوچک، جداسازی با `/`، سبک resource/action:  
  `/user/info`, `/user/follow`, `/post/create`, `/post/feed`, `/user/username/update`
- **کلید map مسیرها:** camelCase توصیفی:  
  `"userInfo"`, `"userFollow"`, `"postCreate"`, `"userUpdateUsername"`

### ۲.۵ کد خطا (Controller و Schema)

- **کد خطای پایه (BaseErrCode):** رشته ۴ رقمی برای هر دامنه، مثلاً `"2001"` برای user، `"2003"` برای post (در `controllerBaseErrCode`).
- **کد بخش (section):** دو رقمی در controller، مثلاً `"01"` برای خطای ParseBody، `"02"` برای خطای repository.
- **کد در Validate:** دو رقمی در schema، مثلاً `"03"`, `"06"`, `"09"` برای فیلدهای مختلف.
- **فرمت نهایی خطا:** `baseErrCode + section + errStr` (مثلاً `"20010603"`).

### ۲.۶ Repository و DataSource

- **اینترفیس repository در models/repositories:**  
  `UserRepository`, `PostRepository`؛ نام فایل: `userRepo.go`, `postRepo.go`.
- **متغیر سراسری:**  
  `var UserRepo UserRepository`, `var PostRepo PostRepository` (Repo با R بزرگ).
- **ساختار در مدل:**  
  `type Repository struct { ... }`؛ متدها روی `(repo *Repository)`؛ سینگلتون با `GetRepo() *Repository` و `sync.Once`.
- **DataSources:**  
  اینترفیسها مثل `UserDBDS`, `UserCacheDS`؛ پیاده‌سازی در زیرپوشه مثل `vitessDS`, `redisDS`.

---

## ۳. الگوی Controller (هندلر HTTP)

### ۳.۱ قالب ثابت هر endpoint

```go
package user

import (
	"apiGolang/apiSchema/commonSchema"
	"apiGolang/apiSchema/userSchema"
	"apiGolang/controllers/mainController"
	"apiGolang/models/repositories"
	"apiGolang/statics/constants/controllerBaseErrCode"
	"github.com/gofiber/fiber/v2"
)

func Info(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "06")
	defer mainController.FinishAPISpan(ctx)

	req := commonSchema.BaseRequest[userSchema.InfoRequest]{}

	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerBaseErrCode.UserErrCode, "01", errStr, code, err)
	}

	res, errStr, code, err := repositories.UserRepo.Info(spanCtx, req, mainController.GetUser(ctx))
	if err != nil {
		return mainController.Error(ctx, controllerBaseErrCode.UserErrCode, "02", errStr, code, err)
	}

	return mainController.Response(ctx, res)
}
```

- **نام تابع:** PascalCase، بدون پیشوند (مثلاً `Info`, `Create`, `Follow`).
- **ترتیب:**  
  ۱) `InitAPI(ctx, sectionErrCode)` و `defer FinishAPISpan(ctx)`  
  ۲) تعریف `req` با `BaseRequest[SchemaRequest]{}`  
  ۳) `ParseBody`؛ در صورت خطا `Error(..., baseErrCode, "01", errStr, code, err)`  
  ۴) فراخوانی repository با `spanCtx, req, GetUser(ctx)` (و در صورت نیاز `GetDevice(ctx)` و غیره)  
  ۵) در صورت خطا `Error(..., baseErrCode, "02", errStr, code, err)`  
  ۶) در نهایت `Response(ctx, res)` یا `Response(ctx, nil)`.
- **section در InitAPI:** دو رقم برای این endpoint (مثلاً `"06"` برای Info، `"02"` برای Create).
- **section در Error:**  
  `"01"` معمولاً خطای اعتبارسنجی/ParseBody، `"02"` (و بعدی‌ها) برای خطای لایه repository/ business.

### ۳.۲ زمانی که پاسخ خالی است

```go
return mainController.Response(ctx, nil)
```

### ۳.۳ استفاده از Device یا پارامتر اضافه

```go
res, errStr, code, err := repositories.SomeRepo.Action(spanCtx, req, mainController.GetUser(ctx), mainController.GetDevice(ctx))
```

---

## ۴. الگوی Route (ثبت endpoint)

- **پکیج:** `route` (در `services/core/route/` یا مشابه).
- **یک map برای هر دامنه:**  
  `var userRoutes = map[string]string{ ... }`  
  کلید: نام توصیفی camelCase، مقدار: مسیر با `/`.
- **تابع setup:**  
  `func SetupUserRoute(app *fiber.App) map[string]string`  
  داخلش فقط `app.Post(userRoutes["key"], HandlerName)` و در پایان `return userRoutes`.
- **هندلرها:** بدون پرانتز؛ نام تابع دقیقاً همانی که در controller تعریف شده (مثلاً `Info`, `Create`).
- **import dot برای controller:**  
  `. "apiGolang/controllers/user"` تا بتوان مستقیم `Info`, `Create` را استفاده کرد.

مثال (خلاصه):

```go
var userRoutes = map[string]string{
	"userInfo":   "/user/info",
	"userFollow": "/user/follow",
}

func SetupUserRoute(app *fiber.App) map[string]string {
	app.Post(userRoutes["userInfo"], Info)
	app.Post(userRoutes["userFollow"], Follow)
	return userRoutes
}
```

- در `route.go` اصلی سرویس، همه `SetupXxxRoute(app)` صدا زده می‌شوند و نتیجه در یک map ادغام می‌شود (مثلاً با `MergeMaps`).

---

## ۵. الگوی ApiSchema (درخواست / پاسخ / اعتبارسنجی)

### ۵.۱ ساختار پکیج برای هر دامنه

- **request.go:** ساختارهای درخواست با تگ `msgpack` و در صورت نیاز `json`, `validate`.
- **response.go:** ساختارهای پاسخ (همیشه با `msgpack` و `json`).
- **validate.go:** متد `Validate(validateExtraData commonSchema.ValidateExtraData) (string, int, error)` برای هر نوع درخواست.
- **checkAccess.go** (در صورت نیاز): متد `CheckAccess(ctx *fiber.Ctx, userPermissionsLimitation permission.PermissionsLimitation) (string, int, error)` برای چک دسترسی.
- **receivers.go** (در صورت نیاز): متدهای کمکی روی نوع‌های همان schema (مثل `PrepareScore`, `GetUniqueKey`).
- **extraData.go** یا پوشه **extraData/** در صورت نیاز.

### ۵.۲ Request

- نام نوع: `XxxRequest` (مثلاً `InfoRequest`, `FollowRequest`, `CreateRequest`).
- فیلدها: تگ اجباری `msgpack:"fieldName"`؛ در صورت استفاده در JSON هم `json:"fieldName"`.
- اعتبارسنجی: از تگ `validate:"..."` استفاده می‌شود (مثلاً `required`, `max=32`, `oneof=follow unfollow`).
- برای IDهای عددی که از کلاینت به صورت رشته می‌آیند: از `commonSchema.ID` یا رشته و تبدیل در Validate با `cast.ToInt64`.

مثال:

```go
type FollowRequest struct {
	FollowingUserID string `msgpack:"followingUserID" validate:"required"`
	Operation      string `msgpack:"operation" validate:"oneof=follow unfollow"`
}
```

### ۵.۳ Validate

- امضای ثابت:  
  `func (req *XxxRequest) Validate(validateExtraData commonSchema.ValidateExtraData) (string, int, error)`  
- در صورت خطا:  
  `return "03", status.StatusBadRequest, err` (یا `customErr.SomeErr`).  
  کدهای دو رقمی (`"03"`, `"06"`, `"09"`, ...) برای تفکیک فیلد/نوع خطا.
- در صورت موفقیت:  
  `return "", status.StatusOK, nil`.
- برای خطای اعتبارسنجی ساختاری از کتابخانه:  
  `errMsg, err := validate.Struct(req)` و سپس `switch errMsg { case "FieldName,required": return "03", ... }`.
- برای چک‌های سفارشی از `validateExtraData.Headers` استفاده می‌شود (مثلاً `AppVersion`, `VerifiedToken`).

### ۵.۴ Response

- نام نوع: `XxxResponse` یا `XxxListResponse` و غیره.
- همه فیلدها با `msgpack` و `json`؛ نام فیلدها camelCase در تگ.

### ۵.۵ CheckAccess (در صورت نیاز)

- فقط برای endpointهایی که محدودیت دسترسی (permission) دارند.
- امضا:  
  `func (req *XxxRequest) CheckAccess(ctx *fiber.Ctx, userPermissionsLimitation permission.PermissionsLimitation) (string, int, error)`  
- در صورت عدم دسترسی:  
  `return "50", status.StatusForbidden, customErr.SomeForbiddenErr`  
- در صورت مجاز:  
  `return "", status.StatusOK, nil`.

---

## ۶. الگوی Repository (مدل و دسترسی به داده)

### ۶.۱ ساختار پکیج مدل (مثلاً user)

- **dataModel/user.go:** ساختارهای مربوط به جدول/کش (با تگ `gorm` و `msgpack`). فیلدهای مجازی با `gorm:"-"`.
- **repository.go:** تعریف `type Repository struct { cacheDS, dbDS, httpDS ... }`، `GetRepo()`، و متدهای کمکی مثل `db()`, `cache()`.
- **repositoryXxx.go:** هر عملیات اصلی در یک فایل (مثلاً `repositoryInfo.go`, `repositoryFollow.go`).

### ۶.۲ امضای متدهای repository

- پارامتر اول: `ctx context.Context` (یا `spanCtx`).
- بعد از آن: `req commonSchema.BaseRequest[xxxSchema.XxxRequest]` و در صورت نیاز `user dataModel.User`, `device dataModel.Device`.
- خروجی:  
  `(res xxxSchema.XxxResponse, errStr string, code int, err error)`  
  یا برای عملیات بدون پاسخ:  
  `(errStr string, code int, err error)`.
- در صورت خطا: مقدار خالی برای `res`، کد دو رقمی برای `errStr`، `status.StatusBadRequest` یا `status.StatusInternalServerError` برای `code`، و یکی از `customErr.*` یا `err` برای خطا.

### ۶.۳ دادهٔ خام (dataModel)

- نام فیلد: PascalCase؛ تگ `gorm:"column:columnName"` برای ستون‌های DB.
- فیلدهای مجازی (غیر DB): با `gorm:"-"`.

### ۶.۴ اینترفیس در models/repositories

- فایل جدا برای هر entity، مثلاً `userRepo.go`: تعریف `UserRepository` و `var UserRepo UserRepository`.
- در init (مثلاً از طریق import پکیج مدل) مقداردهی:  
  `repositories.UserRepo = user.GetRepo()` (یا مشابه، بسته به طراحی پروژه).

---

## ۷. خطاها و ثابت‌ها

### ۷.۱ customErr

- در `statics/customErr/err.go`:  
  `var Xxx = errors.New(constants.Xxx)`  
  متن خطا از `statics/constants` می‌آید تا یکجا قابل ترجمه باشد.

### ۷.۲ constants

- ثابت‌های متنی خطا (برای کاربر و لاگ): در `statics/constants`؛ نام‌ها PascalCase یا با حروف درهم (مثلاً `InvalidUsername`, `YouAreBlocked`).
- کدهای HTTP در `statics/constants/status`:  
  `StatusOK`, `StatusBadRequest`, `StatusForbidden`, `StatusUnauthorized`, `StatusInternalServerError`, `TooManyRequests`.

### ۷.۳ controllerBaseErrCode

- یک ثابت ۴ رقمی رشته برای هر دامنه:  
  `UserErrCode = "2001"`, `PostErrCode = "2003"`, و غیره.

---

## ۸. قراردادهای کد (سبک عمومی)

### ۸.۱ import

- گروه‌بندی: اول پکیج‌های داخلی پروژه (با پیشوند ماژول)، بعد پکیج‌های خارجی (مثلاً fiber، context، ...).
- ترتیب حروف الفبا در هر گروه معمولاً رعایت می‌شود.

### ۸.۲ نام متغیرهای رایج

- `ctx` / `spanCtx`: context.
- `req`: درخواست (ساختار schema).
- `res`: پاسخ.
- `errStr`: کد دو رقمی خطا در لایه business.
- `code`: کد HTTP (int).
- `err`: error.
- `logger`: از `log.GetLoggerFromContext(ctx)`.

### ۸.۳ لاگ

- از `code.ts.co.ir/gaplib/golib/logger` (یا معادل):  
  `log.GetLoggerFromContext(ctx)` و سپس `logger.Info`, `logger.Error`, `logger.Debug` با فیلدهای zap در صورت نیاز.

### ۸.۴ تبدیل نوع

- از کتابخانه داخلی cast استفاده می‌شود:  
  `cast.ForceToInt64`, `cast.ForceToString`, `cast.ForceToBool`.

### ۸.۵ پشتیبانی از MessagePack و JSON

- در request/response همیشه تگ `msgpack` وجود دارد؛ در صورت نیاز `json` هم اضافه می‌شود.
- نوع درخواست با `commonSchema.BaseRequest[T]` و پر شدن هدرها با `FillHeader(ctx)` در mainController.

---

## ۹. تست API

- پکیج تست: `user_tests`, `post_tests` (با underscore).
- نام فایل: `userInfo_test.go`, `postCreate_test.go`.
- استفاده از یک testUtil مشترک که درخواست را با token و DUID و headerهای لازم می‌فرستد.
- نام تابع تست: `TestXxx` (مثلاً `TestInfo`, `TestCreate`).

مثال:

```go
func TestInfo(t *testing.T) {
	req := userSchema.InfoRequest{ Usernames: []string{} }
	testUtils.SendTestReq(req, testUtils.ApiTestDetail{
		ApiUrl:  "/user/info",
		Headers: map[string]string{ "app-version": "2.8", "os": "web" },
		Token:   "...",
		DUID:    "...",
	})
}
```

---

## ۱۰. چک‌لیست برای feature جدید (مشابه ویراستی)

1. **کد خطای پایه:** در `controllerBaseErrCode` یک ثابت ۴ رقمی برای دامنه (در صورت دامنه جدید).
2. **apiSchema:**  
   - `XxxRequest` در request.go  
   - `Validate` در validate.go  
   - در صورت نیاز `CheckAccess` در checkAccess.go  
   - نوع پاسخ در response.go
3. **Controller:** یک فایل جدید در `controllers/domain/action.go` با همان قالب InitAPI → ParseBody → Repo → Response/Error.
4. **Repository:** متد جدید در مدل مربوطه با امضای (ctx, req, user, ...) و خروجی (res, errStr, code, err).
5. **Route:** یک کلید در map مسیرها و یک خط `app.Post(routes["key"], HandlerName)` در SetupXxxRoute.
6. **Route اصلی:** در `route.go` سرویس، `SetupXxxRoute(app)` را صدا بزن و نتیجه را در map کلی ادغام کن.
7. در صورت نیاز: ثابت در constants، خطا در customErr، و ترجمه در statics/translate.

---

## ۱۱. خلاصه برای استفاده در پروژه دیگر

- این فایل را در آن پروژه قرار بده (یا محتوای آن را به قالب آن پروژه تبدیل کن).
- هنگام نوشتن کد جدید یا بازنویسی، به این موارد پایبند باش:  
  **ساختار پوشه‌ها، نام‌گذاری پکیج/فایل/نوع/تابع/route، قالب Controller (InitAPI → ParseBody → Repo → Response/Error)، قالب Request/Response/Validate/CheckAccess، قالب Repository و امضای متدها، قرارداد کد خطا (۴ رقم + ۲ رقم + ۲ رقم)، و استفاده از constants/customErr.**
- در پروژه جدید اگر کتابخانه داخلی (مثل gaplib) وجود ندارد، می‌توان با پکیج‌های استاندارد (مثلاً برای cast، logger) جایگزین کرد؛ ولی **ساختار و نام‌گذاری و الگوها** طبق همین راهنما حفظ شود تا خروجی شبیه کد تیم ویراستی باشد.

با رعایت این راهنما، کد نوشته‌شده از نظر سبک و ساختار با کدبیس ویراستی همخوان خواهد بود.

---

## ۱۲. استفاده در Cursor یا پروژه دیگر

### ۱۲.۱ در پروژه مقصد

1. این فایل (`VIRASTY_CODING_STYLE_GUIDE.md`) را در ریشه یا در `docs/` پروژه مقصد کپی کن.
2. در Cursor: یک **Rule** بساز که اشاره به این فایل کند (مثلاً در `.cursor/rules/virasty-style.mdc`):

   ```markdown
   # Virasty-style backend

   When writing or modifying Go API code (controllers, routes, repositories, schemas), follow the conventions in `docs/VIRASTY_CODING_STYLE_GUIDE.md` (or `VIRASTY_CODING_STYLE_GUIDE.md` at project root). Match:
   - Controller pattern: InitAPI → ParseBody → Repo call → Response/Error
   - Naming: PascalCase for types/funcs, camelCase for routes map keys, snake_case or camelCase for files as in the guide
   - Error codes: base (4 digits) + section (2 digits) + detail (2 digits)
   - Schema: request.go, response.go, validate.go, checkAccess.go when needed
   - Repository: (ctx, req, user, ...) → (res, errStr, code, err)
   ```

3. در اولین پیام به AI بگو: «در این پروژه طبق راهنمای Virasty (فایل VIRASTY_CODING_STYLE_GUIDE.md) کد بزن.»

### ۱۲.۲ جایگزینی ماژول و مسیرها

در پروژه دیگر نام ماژول و مسیر importها عوض می‌شود. در راهنما هر جا `apiGolang` آمده، با نام ماژول پروژه مقصد عوض کن (مثلاً `myapp`). بقیه مسیرها (مثل `controllers/user`, `apiSchema/userSchema`, `models/repositories`) را متناسب با ساختار همان پروژه نگه دار یا طبق همان اصول تطبیق بده.
