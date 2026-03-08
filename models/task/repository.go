package task

import (
	"log"
	"sync"

	"github.com/alianjidaniir-design/SamplePRJ/models/repositories"
	taskDataSources "github.com/alianjidaniir-design/SamplePRJ/models/task/dataSources"
	memoryDataSource "github.com/alianjidaniir-design/SamplePRJ/models/task/dataSources/memory"
	mysqlDataSource "github.com/alianjidaniir-design/SamplePRJ/models/task/dataSources/mysql"
)

type Repository struct {
	cacheDS taskDataSources.TaskCacheDS
	dbDS    taskDataSources.TaskDBDS
	initErr error
}

var (
	once    sync.Once
	repoIns *Repository
)

func GetRepo() *Repository {
	once.Do(func() {
		repoIns = &Repository{
			cacheDS: memoryDataSource.NewTaskCacheDS(),
		}
		repoIns.initializeDataSources()
	})

	return repoIns
}

func init() {
	repositories.TaskRepo = GetRepo()
}

func (repo *Repository) initializeDataSources() {
	mysqlDS, enabled, err := mysqlDataSource.NewTaskDBDSFromEnv()
	if err != nil {
		repo.initErr = err
		return
	}

	if enabled {
		repo.dbDS = mysqlDS
		log.Printf("[task-repository] mysql datasource enabled table=%s", mysqlDS.TableName())
		return
	}

	repo.dbDS = memoryDataSource.NewTaskDBDS(100)
	log.Println("[task-repository] MYSQL_DSN is empty, using memory datasource")
}

func (repo *Repository) db() taskDataSources.TaskDBDS {
	return repo.dbDS
}

func (repo *Repository) cache() taskDataSources.TaskCacheDS {
	return repo.cacheDS
}
