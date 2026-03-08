package task

import (
	"context"
	"fmt"
	"time"

	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/commonSchema"
	"github.com/alianjidaniir-design/SamplePRJ/apiSchema/taskSchema"
	taskDataModel "github.com/alianjidaniir-design/SamplePRJ/models/task/datamodel"
	"github.com/alianjidaniir-design/SamplePRJ/statics/constants/status"
)

func (repo *Repository) List(ctx context.Context, req commonSchema.BaseRequest[taskSchema.ListRequest]) (res taskSchema.ListResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return taskSchema.ListResponse{}, "03", status.StatusInternalServerError, repo.initErr
	}

	if repo.db != nil {
		return repo.listFromMySQL(ctx, req)
	}

	return repo.listFromMemory(req)
}

func (repo *Repository) listFromMemory(req commonSchema.BaseRequest[taskSchema.ListRequest]) (res taskSchema.ListResponse, errStr string, code int, err error) {
	cacheKey := fmt.Sprintf("task:list:page:%d:perPage:%d", req.Body.Page, req.Body.PerPage)

	repo.cacheLock.RLock()
	cachedRes, cacheHit := repo.listCache[cacheKey]
	repo.cacheLock.RUnlock()
	if cacheHit {
		return cloneListResponse(cachedRes), "", status.StatusOK, nil
	}

	repo.lock.RLock()
	clonedTasks := make([]taskDataModel.Task, len(repo.tasks))
	copy(clonedTasks, repo.tasks)
	repo.lock.RUnlock()

	start := (req.Body.Page - 1) * req.Body.PerPage
	if start > len(clonedTasks) {
		start = len(clonedTasks)
	}

	end := start + req.Body.PerPage
	if end > len(clonedTasks) {
		end = len(clonedTasks)
	}

	res = taskSchema.ListResponse{
		Tasks:   clonedTasks[start:end],
		Page:    req.Body.Page,
		PerPage: req.Body.PerPage,
		Total:   len(clonedTasks),
	}

	repo.cacheLock.Lock()
	repo.listCache[cacheKey] = cloneListResponse(res)
	repo.cacheLock.Unlock()

	return res, "", status.StatusOK, nil
}

func (repo *Repository) listFromMySQL(ctx context.Context, req commonSchema.BaseRequest[taskSchema.ListRequest]) (res taskSchema.ListResponse, errStr string, code int, err error) {
	cacheKey := fmt.Sprintf("task:list:page:%d:perPage:%d", req.Body.Page, req.Body.PerPage)

	repo.cacheLock.RLock()
	cachedRes, cacheHit := repo.listCache[cacheKey]
	repo.cacheLock.RUnlock()
	if cacheHit {
		return cloneListResponse(cachedRes), "", status.StatusOK, nil
	}

	offset := (req.Body.Page - 1) * req.Body.PerPage

	rowsQuery := fmt.Sprintf(
		"SELECT id, title, description, created_at FROM %s ORDER BY id ASC LIMIT ? OFFSET ?",
		repo.tableSQL,
	)
	rows, err := repo.db.QueryContext(ctx, rowsQuery, req.Body.PerPage, offset)
	if err != nil {
		return taskSchema.ListResponse{}, "04", status.StatusInternalServerError, err
	}
	defer rows.Close()

	tasks := make([]taskDataModel.Task, 0, req.Body.PerPage)
	for rows.Next() {
		var task taskDataModel.Task
		var createdAt time.Time
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &createdAt); err != nil {
			return taskSchema.ListResponse{}, "05", status.StatusInternalServerError, err
		}
		task.CreatedAt = createdAt.UTC().Format(time.RFC3339)
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return taskSchema.ListResponse{}, "06", status.StatusInternalServerError, err
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", repo.tableSQL)
	total := 0
	if err := repo.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return taskSchema.ListResponse{}, "07", status.StatusInternalServerError, err
	}

	res = taskSchema.ListResponse{
		Tasks:   tasks,
		Page:    req.Body.Page,
		PerPage: req.Body.PerPage,
		Total:   total,
	}

	repo.cacheLock.Lock()
	repo.listCache[cacheKey] = cloneListResponse(res)
	repo.cacheLock.Unlock()

	return res, "", status.StatusOK, nil
}

func cloneListResponse(source taskSchema.ListResponse) taskSchema.ListResponse {
	cloned := source
	cloned.Tasks = make([]taskDataModel.Task, len(source.Tasks))
	copy(cloned.Tasks, source.Tasks)
	return cloned
}
