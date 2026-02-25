package user

import (
	"sync"
	"sync/atomic"

	"github.com/alianjidaniir-design/SamplePRJ/models/datamodel"
	"github.com/alianjidaniir-design/SamplePRJ/models/repositories"
)

type Repository struct {
	idCounter int64
	users     []datamodel.User
	lock      sync.RWMutex
}

var (
	once    sync.Once
	repoIns *Repository
)

func GetRepo() *Repository {
	once.Do(func() {
		repoIns = &Repository{
			idCounter: 10,
			users:     []datamodel.User{},
		}
	})

	return repoIns
}

func (repo *Repository) nextID() int64 {
	return atomic.AddInt64(&repo.idCounter, 1)
}

func init() {
	repositories.UserRepo = GetRepo()
}
