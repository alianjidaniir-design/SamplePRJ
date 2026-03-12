package user

import (
	"sync"

	"github.com/alianjidaniir-design/SamplePRJ/models/repositories"
	userDataModel "github.com/alianjidaniir-design/SamplePRJ/models/user/dataModel"
)

type Repository struct {
	lock      sync.RWMutex
	idCounter int64
	users     []userDataModel.User
}

var (
	once    sync.Once
	repoIns *Repository
)

func GetRepo() *Repository {
	once.Do(func() {
		repoIns = &Repository{
			idCounter: 10,
			users:     []userDataModel.User{},
		}
	})

	return repoIns
}

func init() {
	repositories.UserRepo = GetRepo()
}
