package task

import (
	"context"
	"fmt"
	"testing"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/commonSchema"
	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/taskSchema"
)

func TestListCacheAndInvalidation(t *testing.T) {
	repo := GetRepo()

	dbResetter, ok := repo.dbDS.(interface{ Reset() })
	if !ok {
		t.Skip("cache unit test targets memory datasource mode")
	}
	dbResetter.Reset()

	if cacheResetter, ok := repo.cacheDS.(interface{ Reset() }); ok {
		cacheResetter.Reset()
	}

	createReq := commonSchema.BaseRequest[taskSchema.CreateRequest]{
		Body: taskSchema.CreateRequest{Title: "cache-demo", Description: "v1"},
	}
	_, _, _, err := repo.Create(context.Background(), createReq)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	listReq := commonSchema.BaseRequest[taskSchema.ListRequest]{
		Body: taskSchema.ListRequest{Page: 1, PerPage: 10},
	}

	firstRes, _, _, err := repo.List(context.Background(), listReq)
	if err != nil {
		t.Fatalf("first list failed: %v", err)
	}

	cacheKey := fmt.Sprintf("task:list:page:%d:perPage:%d", listReq.Body.Page, listReq.Body.PerPage)
	_, cacheHit := repo.cacheDS.GetList(cacheKey)
	if !cacheHit {
		t.Fatal("expected cache to be populated after first list")
	}

	secondRes, _, _, err := repo.List(context.Background(), listReq)
	if err != nil {
		t.Fatalf("second list failed: %v", err)
	}
	if len(firstRes.Tasks) != len(secondRes.Tasks) {
		t.Fatalf("cache result mismatch: first=%d second=%d", len(firstRes.Tasks), len(secondRes.Tasks))
	}

	createReq2 := commonSchema.BaseRequest[taskSchema.CreateRequest]{
		Body: taskSchema.CreateRequest{Title: "cache-demo-2", Description: "v2"},
	}
	_, _, _, err = repo.Create(context.Background(), createReq2)
	if err != nil {
		t.Fatalf("second create failed: %v", err)
	}

	_, cacheHit = repo.cacheDS.GetList(cacheKey)
	if cacheHit {
		t.Fatal("expected cache invalidation after create")
	}
}
