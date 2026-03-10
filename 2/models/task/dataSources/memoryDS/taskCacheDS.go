package memoryDS

import (
	"sync"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/taskSchema"
)

type TaskCacheDS struct {
	lock      sync.RWMutex
	listCache map[string]taskSchema.ListResponse
}

func NewTaskCacheDS() *TaskCacheDS {
	return &TaskCacheDS{
		listCache: map[string]taskSchema.ListResponse{},
	}
}

func (ds *TaskCacheDS) GetList(cacheKey string) (taskSchema.ListResponse, bool) {
	ds.lock.RLock()
	res, ok := ds.listCache[cacheKey]
	ds.lock.RUnlock()
	return res, ok
}

func (ds *TaskCacheDS) SetList(cacheKey string, res taskSchema.ListResponse) {
	ds.lock.Lock()
	ds.listCache[cacheKey] = res
	ds.lock.Unlock()
}

func (ds *TaskCacheDS) InvalidateList() {
	ds.lock.Lock()
	ds.listCache = map[string]taskSchema.ListResponse{}
	ds.lock.Unlock()
}

func (ds *TaskCacheDS) Reset() {
	ds.InvalidateList()
}
