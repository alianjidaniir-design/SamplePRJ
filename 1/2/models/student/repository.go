package student

import (
	"Fiber/API/2/models/repositories"
	studentDataSourses "Fiber/API/2/models/student/datasourse"
	mysqlDataSourse "Fiber/API/2/models/student/datasourse/mySqlDS"
	"log"
	"sync"
)

type Repository struct {
	dbDS    studentDataSourses.StudentDBDS
	initErr error
}

var (
	once    sync.Once
	repoIns *Repository
)

func GetRepo() *Repository {
	once.Do(func() {
		repoIns = &Repository{}
		repoIns.initializeDataSources()
	})
	return repoIns
}

func init() {
	repositories.UserRepo = GetRepo()
}

func (repo *Repository) initializeDataSources() {
	mysqlDS, enabled, err := mysqlDataSourse.NewTaskDBDSFromEnv()
	if err != nil {
		repo.initErr = err
		return
	}
	if enabled {
		repo.dbDS = mysqlDS
		log.Printf("[task-repository] mysql datasource enabled table=%s", mysqlDS.TableName())
	}
}

func (repo *Repository) db() studentDataSourses.StudentDBDS {
	return repo.dbDS
}
