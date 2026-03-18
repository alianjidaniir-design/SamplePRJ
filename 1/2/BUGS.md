# گزارش باگ‌ها (پوشه src/Fiber/API/2)

1. scan اشتباه و ناقص در خواندن رکورد دانشجو
فایل: `src/Fiber/API/2/models/student/datasourse/mySqlDS/userDBDS.go:89`
- ترتیب و تعداد ستون‌های SELECT با Scan هم‌خوان نیست (ستون `last_name` اصلاً انتخاب نشده).
- در `Scan` برای `FirstName` و `LastName` آدرس داده نمی‌شود.
نتیجه: خطای runtime از نوع Scan (destination not pointer) یا برگشت داده ناقص.

2. امکان panic هنگام ایجاد دانشجو وقتی دیتاسورس فعال نیست
فایل: `src/Fiber/API/2/models/student/repositoryCreate.go:10`
اگر `MYSQL_DSN` خالی باشد، `repo.dbDS` nil می‌ماند و در `Create` روی nil متد صدا زده می‌شود.
نتیجه: panic در زمان اجرای درخواست Create.

3. نام متغیر محیطی اشتباه (احتمال تایپو)
فایل: `src/Fiber/API/2/models/student/datasourse/mySqlDS/config.go:26`
از `MYSQL_STUDENDSS_TABLE` خوانده می‌شود که به‌احتمال زیاد تایپو است.
نتیجه: مقدار جدول از env خوانده نمی‌شود و همیشه مقدار پیش‌فرض استفاده می‌شود.

4. مسیر import با حروف بزرگ/کوچک ناهماهنگ
فایل: `src/Fiber/API/2/services/core/route/StudentRoute.go:4`
ایمپورت `Fiber/API/2/controllers/Student` است ولی مسیر واقعی پوشه `controllers/student` است.
نتیجه: روی سیستم‌های حساس به حروف (Linux) بیلد شکست می‌خورد.

5. مسیر روت بدون اسلش ابتدایی
فایل: `src/Fiber/API/2/services/core/route/StudentRoute.go:9`
مسیر `student/create` بدون `/` ثبت شده است.
نتیجه: درخواست به `/student/create` ممکن است 404 شود.

6. تعریف دوباره ID در مدل‌های Gorm
فایل: `src/Fiber/API/2/api.go:24`
struct ها `gorm.Model` دارند و همزمان فیلد `ID` جدا تعریف شده است.
نتیجه: ایجاد ستون/کلید اصلی تکراری و رفتار غیرمنتظره در مایگریشن/کوئری‌ها.

7. ParseQuery در نبود Validator وضعیت BadRequest می‌دهد
فایل: `src/Fiber/API/2/controllers/mainController/main.go:45`
اگر req متد `Validate` نداشته باشد، تابع با `StatusBadRequest` و `err=nil` برمی‌گردد.
نتیجه: ممکن است درخواست‌های معتبر بدون دلیل رد شوند.
