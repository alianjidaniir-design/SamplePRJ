package task

import (
	"sync"

	"github.com/alianjidaniir-design/SamplePRJ/models/repositories"
	taskDataSources "github.com/alianjidaniir-design/SamplePRJ/models/task/dataSources"
	memoryDataSource "github.com/alianjidaniir-design/SamplePRJ/models/task/dataSources/memory"
)

type Repository struct {
	cacheDS taskDataSources.TaskCacheDS
	dbDS    taskDataSources.TaskDBDS
}

var (
	once    sync.Once
	repoIns *Repository
)

func GetRepo() *Repository {
	once.Do(func() {
		repoIns = &Repository{
			cacheDS: memoryDataSource.NewTaskCacheDS(),
			dbDS:    memoryDataSource.NewTaskDBDS(100),
		}
	})

	return repoIns
}

func init() {
	repositories.TaskRepo = GetRepo()
}

func (repo *Repository) db() taskDataSources.TaskDBDS {
	return repo.dbDS
}

func (repo *Repository) cache() taskDataSources.TaskCacheDS {
	return repo.cacheDS
}
