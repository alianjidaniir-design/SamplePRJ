package task

import (
	"context"
	"testing"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/commonSchema"
	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/taskSchema"
	"github.com/alianjidaniir-design/SamplePRJ/models/datamodel"
)

func TestListCacheAndInvalidation(t *testing.T) {
	repo := GetRepo()

	repo.lock.Lock()
	repo.tasks = []datamodel.Task{}
	repo.lock.Unlock()

	repo.cacheLock.Lock()
	repo.listCache = map[string]taskSchema.ListResponse{}
	repo.cacheLock.Unlock()

	createReq := commonSchema.BaseRequest[taskSchema.CreateRequest]{
		Body: taskSchema.CreateRequest{Title: "cache-demo", Description: "v1"},
	}
	_, _, _, err := repo.Create(context.Background(), createReq, datamodel.User{})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	listReq := commonSchema.BaseRequest[taskSchema.ListRequest]{
		Body: taskSchema.ListRequest{Page: 1, PerPage: 10},
	}

	firstRes, _, _, err := repo.List(context.Background(), listReq, datamodel.User{})
	if err != nil {
		t.Fatalf("first list failed: %v", err)
	}

	repo.cacheLock.RLock()
	cacheCount := len(repo.listCache)
	repo.cacheLock.RUnlock()
	if cacheCount == 0 {
		t.Fatal("expected cache to be populated after first list")
	}

	secondRes, _, _, err := repo.List(context.Background(), listReq, datamodel.User{})
	if err != nil {
		t.Fatalf("second list failed: %v", err)
	}
	if len(firstRes.Tasks) != len(secondRes.Tasks) {
		t.Fatalf("cache result mismatch: first=%d second=%d", len(firstRes.Tasks), len(secondRes.Tasks))
	}

	createReq2 := commonSchema.BaseRequest[taskSchema.CreateRequest]{
		Body: taskSchema.CreateRequest{Title: "cache-demo-2", Description: "v2"},
	}
	_, _, _, err = repo.Create(context.Background(), createReq2, datamodel.User{})
	if err != nil {
		t.Fatalf("second create failed: %v", err)
	}

	repo.cacheLock.RLock()
	cacheCount = len(repo.listCache)
	repo.cacheLock.RUnlock()
	if cacheCount != 0 {
		t.Fatalf("expected cache invalidation after create, got entries=%d", cacheCount)
	}
}
