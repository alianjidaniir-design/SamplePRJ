package user

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/commonSchema"
	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/userSchema"
	"github.com/alianjidaniir-design/SamplePRJ/models/repositories"
	userDataModel "github.com/alianjidaniir-design/SamplePRJ/models/user/datamodel"
	"github.com/alianjidaniir-design/SamplePRJ/statics/constants/status"
)

type Repository struct {
	idCounter int64
	users     []userDataModel.User
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
			users:     []userDataModel.User{},
		}
	})

	return repoIns
}

func init() {
	repositories.UserRepo = GetRepo()
}

func (repo *Repository) nextID() int64 {
	return atomic.AddInt64(&repo.idCounter, 1)
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[userSchema.CreateRequest], user userDataModel.User) (res userSchema.CreateResponse, errStr string, code int, err error) {
	_ = ctx
	_ = user

	newUser := userDataModel.User{
		ID:       repo.nextID(),
		Username: req.Body.Username,
		Email:    req.Body.Email,
	}

	repo.lock.Lock()
	repo.users = append(repo.users, newUser)
	repo.lock.Unlock()

	res = userSchema.CreateResponse{User: newUser}
	return res, "", status.StatusOK, nil
}
