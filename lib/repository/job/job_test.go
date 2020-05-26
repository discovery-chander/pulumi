package job

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/EurosportDigital/global-transcoding-platform/lib/errors"
	"github.com/EurosportDigital/global-transcoding-platform/lib/logger"
	repositories "github.com/EurosportDigital/global-transcoding-platform/lib/repository"
	"github.com/EurosportDigital/global-transcoding-platform/model"
	"github.com/EurosportDigital/global-transcoding-platform/model/gormmodel"
	"github.com/jinzhu/gorm"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/suite"
)

var mockError = fmt.Errorf("Mock error")

const readyStatus = "ready"

type JobTestSuite struct {
	suite.Suite
	database   *gorm.DB
	repository JobRepository
	tableName  string
}

func (suite *JobTestSuite) SetupTest() {
	mocket.Catcher.Register() // Safe register. Allowed multiple calls to save
	mocket.Catcher.Logging = true

	// GORM
	db, err := gorm.Open(mocket.DriverName, "connection_string") // Can be any connection string
	db.LogMode(true)
	if err != nil {
		panic(err)
	}
	suite.database = db
	suite.repository = New(db)
	suite.tableName = "jobs"
}

func (suite *JobTestSuite) TearDownTest() {
	suite.database.Close()
}

func buildJobPayload(id int, priority int, status model.Status) map[string]interface{} {
	outputs, _ := json.Marshal([]*gormmodel.Output{})
	rawStatus, _ := json.Marshal(gormmodel.JobStatus{Status: status})
	return map[string]interface{}{
		"id":            id,
		"priority":      priority,
		"status":        rawStatus,
		"source_path":   "s3://mock-bucket/mock",
		"preroll_path":  "s3://mock-bucket/mockpreroll",
		"postroll_path": "s3://mock-bucket/mockpostroll",
		"outputs":       outputs,
	}
}

func buildCountPayload(count int) map[string]interface{} {
	return map[string]interface{}{
		"count": count,
	}
}

func newMockJob(id int) *model.Job {
	return &model.Job{
		ID:           id,
		Priority:     3,
		Status:       model.JobStatus{Status: model.StatusFailed},
		SourcePath:   "s3://mock-bucket/mock",
		PrerollPath:  "s3://mock-bucket/mockpreroll",
		PostrollPath: "s3://mock-bucket/mockpostroll",
	}
}

func (suite *JobTestSuite) TestGet() {
	query := fmt.Sprintf(`SELECT * FROM "%[1]s"  WHERE ("%[1]s"."id" = 1) ORDER BY "%[1]s"."id" ASC LIMIT 1`, suite.tableName)
	suite.Run("Should return get", func() {
		mocket.Catcher.Reset()
		logger.Infof("Query  %s", query)
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  query,
				Response: []map[string]interface{}{buildJobPayload(1, 1, readyStatus)},
				Once:     true,
			},
		})
		result, err := suite.repository.Get(1)
		logger.Infof("Result %v", result)
		suite.Require().NoError(err, "Invoking method should not produce an error")
		suite.Require().Equal(1, result.ID, "Result should not have a different result than expected")
	})
	suite.Run("Should fail if no records are obtained", func() {
		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern: query,
				Once:    true,
				Error:   gorm.ErrRecordNotFound,
			},
		})
		result, err := suite.repository.Get(1)
		var nilJob *model.Job = nil
		logger.Infof("Result %v", result)
		suite.Require().Equal(nilJob, result, "Result shouldn't be anything but empty struct")
		suite.Require().EqualError(errors.Cause(err), repositories.ErrEntityNotFound.Error(), "Error shouldn't be different than expected")
	})
	suite.Run("Should return errors if any", func() {
		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern: query,
				Once:    true,
				Error:   mockError,
			},
		})
		var nilJob *model.Job = nil
		result, err := suite.repository.Get(1)
		logger.Infof("Result %v", result)
		suite.Require().Equal(nilJob, result, "Result shouldn't be anything but empty struct")
		suite.Require().EqualError(errors.Cause(err), mockError.Error(), "Error shouldn't be different than expected")
	})
}
func (suite *JobTestSuite) TestList() {
	countQuery := fmt.Sprintf(`SELECT count(*) FROM "%s"`, suite.tableName)
	countQueryWithFilter := fmt.Sprintf(`SELECT count(*) FROM "%s"  WHERE`, suite.tableName)
	listQuery := fmt.Sprintf(`SELECT * FROM "%s"   LIMIT 3 OFFSET 0`, suite.tableName)
	listQueryWithFilter := fmt.Sprintf(`SELECT * FROM "%s"  WHERE`, suite.tableName)

	suite.Run("Should return list with pagination", func() {
		mocket.Catcher.Reset()
		paginationMock := &JobPagination{Size: 3, Page: 1}
		logger.Infof("Count Query  %s", countQuery)
		logger.Infof("List Query  %s", listQuery)
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  countQuery,
				Response: []map[string]interface{}{buildCountPayload(2)},
				Once:     true,
			},
			{
				Pattern:  listQuery,
				Response: []map[string]interface{}{buildJobPayload(1, 1, readyStatus), buildJobPayload(2, 1, readyStatus)},
				Once:     true,
			},
		})
		result, err := suite.repository.All(nil, paginationMock)
		suite.Require().NoError(err, "Invoking method should not produce an error")
		suite.Require().Equal(2, len(result.Results), "Calling get should return a result")
	})

	suite.Run("Should return list with pagination and filtered by priority and status", func() {
		mocket.Catcher.Reset()
		paginationMock := &JobPagination{Size: 3, Page: 1}
		mockedStatus := model.StatusProcessing
		mockedPriority := 2
		filterMock := &JobFilter{Status: &mockedStatus, Priority: &mockedPriority}
		logger.Infof("Count Query  %s", countQueryWithFilter)
		logger.Infof("List Query  %s", listQueryWithFilter)
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  countQueryWithFilter,
				Response: []map[string]interface{}{buildCountPayload(5)},
				Once:     true,
			},
			{
				Pattern:  listQueryWithFilter,
				Response: []map[string]interface{}{buildJobPayload(1, mockedPriority, mockedStatus), buildJobPayload(2, mockedPriority, mockedStatus), buildJobPayload(3, mockedPriority, mockedStatus)},
				Once:     true,
			},
		})
		result, err := suite.repository.All(filterMock, paginationMock)
		suite.Require().NoError(err, "Invoking method should not produce an error")
		suite.Require().Equal(3, len(result.Results), "Calling get should return a result")
	})

	suite.Run("Should return empty list", func() {
		mocket.Catcher.Reset()
		paginationMock := &JobPagination{Size: 3, Page: 1}
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  countQuery,
				Response: []map[string]interface{}{buildCountPayload(0)},
				Once:     true,
			},
		})
		result, err := suite.repository.All(nil, paginationMock)
		suite.Require().NoError(err, "Invoking method should not produce an error")
		suite.Require().Equal(0, len(result.Results), "Calling get should return empty list")
	})

	suite.Run("Should return error if any", func() {
		mocket.Catcher.Reset()
		paginationMock := &JobPagination{Size: 3, Page: 1}
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern: countQuery,
				Once:    true,
				Error:   mockError,
			},
			{
				Pattern: listQuery,
				Once:    true,
				Error:   mockError,
			},
		})
		result, err := suite.repository.All(nil, paginationMock)
		suite.Require().EqualError(errors.Cause(err), mockError.Error(), "Invoking method should not produce an error")
		suite.Require().Nil(result, "Should return empty job pagination result")
	})
}

func (suite *JobTestSuite) TestDelete() {
	query := fmt.Sprintf(`DELETE FROM "%[1]s"  WHERE "%[1]s"."id" = ?`, suite.tableName)
	logger.Infof("Query  %s", query)
	suite.Run("Should return delete", func() {
		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:      query,
				RowsAffected: 1,
				Once:         true,
			},
		})
		err := suite.repository.Delete(1)
		suite.Require().NoError(err, "Invoking method should not produce an error")
	})
	suite.Run("Should fail if no records are deleted", func() {
		mocket.Catcher.Reset()
		err := suite.repository.Delete(2)
		suite.Require().EqualError(errors.Cause(err), repositories.ErrEntityNotFound.Error(), "Error shouldn't be different than expected")
	})
	suite.Run("Should return errors if any", func() {
		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern: query,
				Once:    true,
				Error:   mockError,
			},
		})
		err := suite.repository.Delete(1)
		suite.Require().EqualError(errors.Cause(err), mockError.Error(), "Invoking method should not produce an error")
	})
}

func (suite *JobTestSuite) TestUpdate() {
	job := &model.Job{
		ID:           2,
		Priority:     3,
		Status:       model.JobStatus{Status: model.StatusFailed},
		SourcePath:   "s3://mock-bucket/mock",
		PrerollPath:  "s3://mock-bucket/mockpreroll",
		PostrollPath: "s3://mock-bucket/mockpostroll",
	}
	query := fmt.Sprintf(`UPDATE "%[1]s"`, suite.tableName)
	logger.Infof("Query  %s", query)
	suite.Run("Should return update", func() {
		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:      query,
				RowsAffected: 1,
				Once:         true,
			},
		})
		err := suite.repository.Update(job)
		suite.Require().NoError(err, "Invoking method should not produce an error")
	})
	suite.Run("Should fail if no records are updated", func() {
		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:      query,
				RowsAffected: 0,
				Once:         true,
			},
		})
		err := suite.repository.Update(job)
		suite.Require().EqualError(errors.Cause(err), repositories.ErrEntityNotFound.Error(), "Error shouldn't be different than expected")
	})
	suite.Run("Should return errors if any", func() {
		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern: query,
				Once:    true,
				Error:   mockError,
			},
		})
		err := suite.repository.Update(job)
		suite.Require().EqualError(errors.Cause(err), mockError.Error(), "Invoking method should not produce an error")
	})
}

func (suite *JobTestSuite) TestInsert() {
	job := newMockJob(2)
	query := fmt.Sprintf(`INSERT INTO "%s"`, suite.tableName)
	logger.Infof("Query  %s", query)
	suite.Run("Should return update", func() {
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern: query,
				Once:    true,
			},
		})
		err := suite.repository.Create(job)
		suite.Require().NoError(err, "Invoking method should not produce an error")
	})
	suite.Run("Should fail if no records are insert", func() {
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern: query,
				Once:    true,
				Error:   fmt.Errorf("Mock error"),
			},
		})
		mocket.Catcher.NewMock().IsMatch(query, []driver.NamedValue{})
		err := suite.repository.Create(job)
		suite.Require().EqualError(errors.Cause(err), "Mock error", "Should return errors if the insert was not successful")
	})
	suite.Run("Should return errors if any", func() {
		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern: query,
				Once:    true,
				Error:   mockError,
			},
		})
		err := suite.repository.Create(job)
		suite.Require().EqualError(errors.Cause(err), mockError.Error(), "Invoking method should not produce an error")
	})
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestJobTestSuite(t *testing.T) {
	suite.Run(t, new(JobTestSuite))
}
