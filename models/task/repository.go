package task

import (
	"sync"
	"sync/atomic"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/taskSchema"
	"github.com/alianjidaniir-design/SamplePRJ/models/datamodel"
	"github.com/alianjidaniir-design/SamplePRJ/models/repositories"
)

type Repository struct {
	idCounter int64
	tasks     []datamodel.Task
	lock      sync.RWMutex
	listCache map[string]taskSchema.ListResponse
	cacheLock sync.RWMutex
}

var (
	once    sync.Once
	repoIns *Repository
)

func GetRepo() *Repository {
	once.Do(func() {
		repoIns = &Repository{
			idCounter: 100,
			tasks:     []datamodel.Task{},
			listCache: map[string]taskSchema.ListResponse{},
		}
	})
	return repoIns
}

func (repo *Repository) nextID() int64 {
	return atomic.AddInt64(&repo.idCounter, 1)
}

func init() {
	repositories.TaskRepo = GetRepo()
}
