package job

import (
	"fmt"

	"github.com/EurosportDigital/global-transcoding-platform/lib/errors"
	"github.com/EurosportDigital/global-transcoding-platform/model/gormmodel"
	"github.com/jinzhu/gorm"

	// Required by GORM for postgres requests
	"github.com/EurosportDigital/global-transcoding-platform/lib/logger"
	"github.com/EurosportDigital/global-transcoding-platform/lib/repository"
	"github.com/EurosportDigital/global-transcoding-platform/model"
)

type JobRepository interface {
	Get(id int) (*model.Job, error)
	Create(job *model.Job) error
	Update(job *model.Job) error
	Delete(id int) error
	All(filters *JobFilter, pagination *JobPagination) (*JobPaginationResult, error)
}

type JobFilter struct {
	Status   *model.Status
	Priority *int
}

type JobPagination struct {
	Size  int
	Page  int
	total int
}

type JobPaginationResult struct {
	Results []*model.Job
	Total   int
	Page    int
}

type gormJobRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) JobRepository {
	return &gormJobRepository{
		db: db,
	}
}

func (instance *gormJobRepository) Get(id int) (*model.Job, error) {
	logger.Infof("Getting Job with Id: %d", id)
	job := &gormmodel.Job{}
	err := instance.db.First(job, id).Error
	if gorm.IsRecordNotFoundError(err) {
		logger.Warnf("Job with id %d not found", id)
		return nil, errors.Wrapf(repository.ErrEntityNotFound, "job id %v not found", id)
	}
	if err != nil {
		logger.Errorf("An  error ocurred while trying to get job %v", err)
		return nil, errors.Wrapf(err, "unable to get job %v", id)
	}
	return gormmodel.ToJob(job), nil
}

func (instance *gormJobRepository) All(filters *JobFilter, pagination *JobPagination) (*JobPaginationResult, error) {
	logger.Infof("Listing all jobs")
	jobs := []*gormmodel.Job{}

	dbInstance := instance.addFilters(filters)

	limit := pagination.Size
	pagination.total = 0
	offset := (pagination.Page - 1) * limit

	listErr := dbInstance.Model(&jobs).Count(&pagination.total).Limit(limit).Offset(offset).Find(&jobs).Error

	if listErr != nil {
		logger.Errorf("An error occurred while trying to list jobs %v", listErr)
		return nil, errors.Wrapf(listErr, "unable to list all jobs on page %v with size %v", pagination.Page, limit)
	}
	modelJobs := make([]*model.Job, len(jobs))
	for i := range modelJobs {
		modelJobs[i] = gormmodel.ToJob(jobs[i])
	}
	return &JobPaginationResult{Results: modelJobs, Total: pagination.total, Page: pagination.Page}, nil
}

func (instance *gormJobRepository) Update(job *model.Job) error {
	logger.Infof("Updating job: %+v", job)
	gormJob := gormmodel.ToGormJob(job)
	result := instance.db.Model(&gormJob).Updates(gormJob)
	if result.Error != nil {
		logger.Error(result.Error, "Error found when trying to update record")
		return errors.Wrapf(result.Error, "updating job %v", safeGetJobID(job))
	}
	if result.RowsAffected == 0 {
		logger.Warnf("Could not find record to be updated")
		return errors.Wrapf(repository.ErrEntityNotFound, "job %v not found", safeGetJobID(job))
	}
	return nil
}

func (instance *gormJobRepository) Create(job *model.Job) error {
	logger.Infof("Creating Job %+v", job)
	gormJob := gormmodel.ToGormJob(job)
	result := instance.db.Create(gormJob)
	logger.Infof("Affected rows: %d", result.RowsAffected)
	if result.Error != nil {
		logger.Error(result.Error, "Error found when trying create job")
		return errors.Wrapf(result.Error, "creating job %v", safeGetJobID(job))
	}
	*job = *gormmodel.ToJob(gormJob)
	return nil
}

func (instance *gormJobRepository) Delete(id int) error {
	logger.Infof("Deleting job with Id: %d", id)
	result := instance.db.Delete(&gormmodel.Job{ID: id})
	if result.Error != nil {
		logger.Error(result.Error, "Error found when deleting job")
		return errors.Wrapf(result.Error, "deleting job %v", id)
	}
	if result.RowsAffected == 0 {
		logger.Warnf("Attempting to delete record that was not found")
		return errors.Wrapf(repository.ErrEntityNotFound, "deleting job %v", id)
	}
	return nil
}

func (instance *gormJobRepository) addFilters(filters *JobFilter) *gorm.DB {
	var dbInstance *gorm.DB

	dbInstance = instance.db

	if filters != nil {
		if filters.Priority != nil {
			dbInstance = dbInstance.Where("priority = ?", filters.Priority)
		}

		if filters.Status != nil {
			dbInstance = dbInstance.Where("status->>'status' = ?", filters.Status)
		}
	}
	return dbInstance
}

func safeGetJobID(job *model.Job) string {
	if job != nil {
		return fmt.Sprintf("%v", job.ID)
	}
	return "unknown"
}
