package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EnrollmentStatus string

const (
	StatusCanceled EnrollmentStatus = "canceled"
	StatusEnrolled EnrollmentStatus = "enrolled"
)

type Students struct {
	gorm.Model
	StudentCode string `gorm:"size:32;unique;not null"`
	FirstName   string `gorm:"size:64;not null"`
	LastName    string `gorm:"size:64;not null"`
	ID          uint   `gorm:"primaryKey;autoIncrement;not null"`
}

type Courses struct {
	gorm.Model
	CourseCode    string `gorm:"size:32;unique;not null"`
	Title         string `gorm:"size:128;not null"`
	Capacity      int    `gorm:"not null"`
	EnrolledCount int    `gorm:"default:0;not null"`
	IsActive      bool   `gorm:"not null;default:true"`
	ID            uint   `gorm:"primaryKey;autoIncrement;not null"`
}

type Enrollments struct {
	gorm.Model
	Status     EnrollmentStatus `gorm:"size : 10 ;not null"`
	CanceledAt time.Time
	EnrolledAt time.Time `gorm:"autoCreateTime:milli"`
	StudentId  uint      `gorm:"not null"`
	CourseId   uint      `gorm:"not null"`
	ID         uint      `gorm:"primaryKey;autoIncrement"`
}

func database() *gorm.DB {
	dsn := "root:123456@tcp(127.0.0.1:3306)/ali-db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Students{}, &Courses{}, &Enrollments{})
	return db
}

func CreateUser(c fiber.Ctx) error {
	students := new(Students)

	db2 := database()

	if err := c.Bind().JSON(students); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := db2.Create(students).Error; err != nil {
		fmt.Println("err", err, 12*22)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	c.Status(fiber.StatusOK)
	return c.Status(201).JSON(students)

}

func listUsers(c fiber.Ctx) error {
	var students []Students

	db2 := database()

	if err := db2.Find(&students, c.Params("id")).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	q := c.Query("q")
	if q == "" {
		return errors.New("query is empty")
	}
	search := "%" + strings.ToLower(q) + "%"
	if err := db2.Model(&Students{}).Where("LOWER(first_name) LIKE ? ", search).Find(&students).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(students)
}

func GetUsers(c fiber.Ctx) error {
	var students Students

	db2 := database()
	if err := db2.Find(&students, c.Params("id")).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(students)
}

func CreateCourse(c fiber.Ctx) error {
	var course Courses

	db2 := database()

	if err := c.Bind().JSON(&course); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := db2.Create(&course).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(course)
}

func UpdateUser(c fiber.Ctx) error {

	updateData := Students{
		StudentCode: "40506070809",
		FirstName:   "John",
		LastName:    "Smith",
	}
	var students Students

	db2 := database()

	if err := db2.Model(&students).Where("id = ?", c.Params("id")).Updates(updateData).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(students)
}

func UpdateUser2(c fiber.Ctx) error {
	var students Students

	db2 := database()

	d := c.Params("id")
	f, _ := strconv.Atoi(d)
	updateUser := Students{
		ID:          uint(f),
		StudentCode: "0928428941",
		FirstName:   "Ali",
		LastName:    "Ali",
	}

	if err := db2.Model(&Students{ID: updateUser.ID}).Updates(&updateUser).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(students)

}

func DeleteUser(c fiber.Ctx) error {
	var students Students

	db2 := database()

	if err := db2.Delete(&students, c.Params("id")).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(students)

}

func ListCourses(c fiber.Ctx) error {
	var courses []Courses

	db2 := database()

	if err := db2.Find(&courses).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error1": err.Error()})
	}

	q := c.Query("q")
	if q == "" {
		return errors.New("query is empty")
	}

	search := "%" + strings.ToLower(q) + "%"
	if err := db2.Model(&Courses{}).Where("LOWER(title) LIKE ? ", search).Find(&courses).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error1": err.Error()})
	}
	return c.Status(200).JSON(courses)

}

func Recourses(c fiber.Ctx) error {
	var courses Courses

	db2 := database()

	if err := db2.Find(&courses, c.Params("id")).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error2": err.Error()})
	}
	return c.JSON(courses)
}

func DeleteCourse(c fiber.Ctx) error {
	var courses Courses

	db2 := database()

	if err := db2.Delete(&courses, c.Params("id")).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(courses)
}

func disactiveCourse(c string, f string, db2 *gorm.DB) error {

	var courses Courses

	if err := db2.Model(&courses).Select("id , is_active").Where(" id = ? AND is_active = ?", f, true).Find(&courses).Error; err != nil {
		return err
	}
	courses.IsActive = false
	return db2.Model(&courses).Select("is_active").Updates(&courses).Error
}

func handledis(c fiber.Ctx) error {
	db2 := database()
	tt := c.Params(c.Query("deactivate"))
	ff := c.Params("id")
	if err := disactiveCourse(tt, ff, db2); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"success": true})
}

func UpdateCourse2(c fiber.Ctx) error {
	var courses Courses

	db2 := database()

	d := c.Params("id")
	f, _ := strconv.Atoi(d)
	asd := Courses{
		ID:         uint(f),
		CourseCode: c.Params("course_code"),
		Title:      c.Params("title"),
	}
	if err := db2.Model(&Courses{ID: asd.ID}).Updates(&asd).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(courses)
}

func CreateEnrollment(c fiber.Ctx, tx *gorm.DB) error {
	var enrollment Enrollments
	var course Courses
	var student Students

	if err := c.Bind().JSON(&enrollment); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if err := tx.First(&student, enrollment.StudentId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"code": 404, "message": "student not found"})
		}
	}

	if err := tx.Clauses(clause.Locking{Strength: "Update"}).First(&course, "id = ?", enrollment.CourseId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"code": 404, "message": "course not found"})
		}
	}

	if course.Capacity <= (course).EnrolledCount {
		return c.Status(409).JSON(fiber.Map{"code": 409, "massage": "capacity is completed"})
	}

	enrollment.Status = StatusEnrolled

	if err := tx.Create(&enrollment).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	course.EnrolledCount++

	if err := tx.Save(&course).Error; err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"err": err.Error()})
	}

	return c.Status(200).JSON(enrollment)
}

func Cancel(c fiber.Ctx, tx *gorm.DB) error {
	var enrollment Enrollments
	var course Courses

	if err := tx.First(&enrollment, "id = ?", c.Params("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"code": 404, "message": "student not found in enrollments"})
		}
	}
	if err := tx.First(&enrollment, enrollment.Status).Error; err != nil {
		if enrollment.Status == StatusCanceled {
			return c.Status(409).JSON(fiber.Map{"code": 409, "message": "enrollment is canceled"})
		}
	}

	f := c.Params("id")
	d, err := strconv.Atoi(f)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	cor := Enrollments{
		ID:        uint(d),
		StudentId: enrollment.StudentId,
		CourseId:  enrollment.CourseId,
	}

	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&course, cor.CourseId).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error() + "Bye"})
	}

	if course.Capacity < 0 {
		return c.Status(409).JSON(fiber.Map{"code": 409, "message": "student capacity can not be less than 0"})
	}

	enrollment.Status = StatusCanceled
	enrollment.CanceledAt = time.Now()

	if err := tx.Save(&enrollment).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	course.EnrolledCount--

	if err := tx.Save(&course).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(enrollment)
}

func ErrorHandler(c fiber.Ctx) error {
	db2 := database()

	err := db2.Transaction(func(tx *gorm.DB) error {
		return CreateEnrollment(c, db2)
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err})
	}
	return nil
}

func chandler(c fiber.Ctx) error {
	db2 := database()
	err := db2.Transaction(func(tx *gorm.DB) error {
		return Cancel(c, db2)
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err})
	}
	return nil
}

func paginating(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page < 1 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 10
		} else if pageSize > 100 {
			pageSize = 100
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func Get23(db2 *gorm.DB, page, pageSize int) ([]Enrollments, int64, error) {
	var enrollments []Enrollments
	var i int64
	if err := db2.Model(&enrollments).Count(&i).Error; err != nil {
		return nil, 0, err
	}
	if err := db2.Scopes(paginating(page, pageSize)).First(&enrollments).Error; err != nil {
		return nil, 0, err
	}
	return enrollments, i, nil

}

func GetHandler(c fiber.Ctx) error {
	db2 := database()
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))
	users, total, err := Get23(db2, page, pageSize)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err})
	}
	return c.JSON(fiber.Map{
		"data":     users,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

type sssss struct {
	StudentId uint
	Status    string
}

func std(s string, f string, db2 *gorm.DB) ([]sssss, error) {
	var enrollment []Enrollments
	var sss []sssss
	if err := db2.Model(enrollment).Select("student_id").Where("course_id = ? ", f).Find(&sss).Error; err != nil {
		return nil, err
	}
	if s == "enrolled" {
		if err := db2.Model(enrollment).Select("student_id", "status").Where("course_id = ? AND status = ? ", f, s).Find(&sss).Error; err != nil {
			return nil, err
		}
		return sss, nil
	}

	return sss, nil
}

func StatusHandler(c fiber.Ctx) error {

	db2 := database()
	id := c.Params("course_id")
	st := c.Query("status")

	f1, err := std(st, id, db2)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err})
	}
	return c.JSON(fiber.Map{
		"data": f1,
	})

}

type ssssss struct {
	CourseId uint
	Status   string
}

func coursecreator(s string, f string, db *gorm.DB) ([]ssssss, error) {
	var enrollment []Enrollments
	var ssss []ssssss
	if ere := db.Model(enrollment).Select("course_id").Where("student_id = ? ", f).Find(&ssss).Error; ere != nil {
		return nil, ere
	}
	if s == "enrolled" {
		if err := db.Model(enrollment).Select("course_id", "status").Where("student_id = ? AND status = ? ", f, s).Find(&ssss).Error; err != nil {
			return nil, err
		}
		return ssss, nil
	}
	return ssss, nil
}

func handlecs(c fiber.Ctx) error {
	db2 := database()
	id := c.Params("student_id")
	status := c.Query("status")
	f5, err := coursecreator(status, id, db2)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err})
	}
	return c.JSON(fiber.Map{"data": f5})

}

func main() {
	app := fiber.New()

	api := app.Group("/api")
	api.Post("/v1/students", CreateUser)
	api.Get("/v2/students", listUsers)
	api.Get("/v1/students/:id", GetUsers)
	api.Patch("/v1/students/:id", UpdateUser)
	api.Delete("/v1/students/:id", DeleteUser)
	api.Post("/v1/courses", CreateCourse)
	api.Get("/v1/courses", ListCourses)
	api.Get("/v1/courses/:id", Recourses)
	api.Delete("/v1/courses/:id", DeleteCourse)
	api.Patch("/v1/courses/:id/deactivate", handledis)
	api.Put("/v1/students/:id", UpdateUser2)
	api.Put("/v1/courses/:id", UpdateCourse2)
	api.Post("/v1/enrollment", ErrorHandler)
	api.Post("v1/enrollment/:id/cancel", chandler)
	api.Get("/v1/enrollment", GetHandler)
	api.Get("/v2/courses/:course_id/students", StatusHandler)
	api.Get("/v2/students/:student_id/courses", handlecs)

	log.Fatal(app.Listen(":3000"))

}
