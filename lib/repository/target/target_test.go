package target

import (
	stderrors "errors"
	"testing"

	"github.com/EurosportDigital/global-transcoding-platform/lib/repository"
	"github.com/EurosportDigital/global-transcoding-platform/model"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type targetTestSuite struct {
	suite.Suite
	db         *gorm.DB
	targetRepo Repository
}

func TestTargetTestSuite(t *testing.T) {
	mocket.Catcher.Register()
	mocket.Catcher.Logging = false

	db, err := gorm.Open(mocket.DriverName, "")
	require.NoError(t, err)
	pSuite := &targetTestSuite{
		db:         db,
		targetRepo: New(db),
	}
	suite.Run(t, pSuite)
}

func (pts *targetTestSuite) TearDownSuite() {
	pts.db.Close()
}

func (pts *targetTestSuite) SetupTest() {
	mocket.Catcher.Reset()
}

func (pts *targetTestSuite) TestGormTargetGet() {
	getQuery := `SELECT * FROM "targets"  WHERE ("targets"."id" = 1) ORDER BY "targets"."id" ASC LIMIT 1`
	pts.Run("Should return expected result", func() {
		pts.SetupTest()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  getQuery,
				Response: []map[string]interface{}{{"id": 1}},
				Once:     true,
			},
		})

		target, err := pts.targetRepo.Get(1)
		pts.Require().NoError(err)
		pts.Require().EqualValues(1, target.ID)
	})
	pts.Run("Should return repository.ErrEntityNotFound if target is not found", func() {
		pts.SetupTest()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  getQuery,
				Response: []map[string]interface{}{},
				Once:     true,
			},
		})

		_, err := pts.targetRepo.Get(1)
		pts.Require().Error(err)
		pts.Require().EqualError(errors.Cause(err), repository.ErrEntityNotFound.Error())
	})
	pts.Run("Should bubble up any unhandled error", func() {
		pts.SetupTest()
		expectedError := stderrors.New("my error")
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern: getQuery,
				Once:    true,
				Error:   expectedError,
			},
		})

		_, err := pts.targetRepo.Get(1)
		pts.Require().Error(err)
		pts.Require().EqualError(errors.Cause(err), expectedError.Error())
	})
}

func (pts *targetTestSuite) TestGormTargetGetMany() {
	pts.Run("Should return expected result", func() {
		pts.SetupTest()
		expectedIDs := []int{1, 2}
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  `SELECT * FROM "targets"  WHERE (id IN (1,2))`,
				Args:     []interface{}{int64(expectedIDs[0]), int64(expectedIDs[1])},
				Response: []map[string]interface{}{{"id": expectedIDs[0]}, {"id": expectedIDs[1]}},
				Once:     true,
			},
		})

		targets, err := pts.targetRepo.GetMany(expectedIDs)
		pts.Require().NoError(err)
		pts.Require().Len(targets, len(expectedIDs))
		for i := range targets {
			pts.Require().EqualValues(expectedIDs[i], targets[i].ID)
		}
	})
	pts.Run("Should bubble up any unhandled error", func() {
		pts.SetupTest()
		expectedError := stderrors.New("my error")
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern: `SELECT * FROM "targets"  WHERE (id IN (1))`,
				Args:    []interface{}{int64(1)},
				Once:    true,
				Error:   expectedError,
			},
		})

		_, err := pts.targetRepo.GetMany([]int{1})
		pts.Require().Error(err)
		pts.Require().EqualError(errors.Cause(err), expectedError.Error())
	})
}

func (pts *targetTestSuite) TestGormTargetCreate() {
	pts.Run("Should create new target in database", func() {
		pts.SetupTest()
		newTarget := &model.Target{
			ID:         10,
			AuthKey:    "authkey",
			Path:       "path",
			TargetType: "type",
		}

		targetInsert := &mocket.FakeResponse{
			Pattern: `INSERT INTO "targets" ("id","target_type","path","auth_key") VALUES (?,?,?,?)`,
			Args: []interface{}{
				int64(newTarget.ID),
				newTarget.TargetType,
				newTarget.Path,
				newTarget.AuthKey,
			},
			Once: true,
		}
		mocket.Catcher.Attach([]*mocket.FakeResponse{targetInsert})

		err := pts.targetRepo.Create(newTarget)
		pts.Require().NoError(err)
		pts.Require().True(targetInsert.Triggered, "target insert reference must be triggered")
	})

	pts.Run("Should bubble up any unhandled error", func() {
		pts.SetupTest()
		expectedError := stderrors.New("my error")
		newTarget := model.Target{}
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern: `INSERT INTO "targets" ("target_type","path","auth_key") VALUES (?,?,?)`,
				Args: []interface{}{
					newTarget.TargetType,
					newTarget.Path,
					newTarget.AuthKey,
				},
				Once:  true,
				Error: expectedError,
			},
		})

		err := pts.targetRepo.Create(&newTarget)
		pts.Require().Error(err)
		pts.Require().EqualError(errors.Cause(err), expectedError.Error())
	})
}

func (pts *targetTestSuite) TestGormTargetUpdate() {
	pts.Run("Should update an existing target in database", func() {
		pts.SetupTest()
		newTarget := &model.Target{
			ID:         10,
			AuthKey:    "my auth key",
			Path:       "my path",
			TargetType: "my type",
		}

		err := pts.targetRepo.Create(newTarget)
		pts.Require().NoError(err)
		newTarget.AuthKey = "updated auth key"

		targetUpdate := &mocket.FakeResponse{
			Pattern: `UPDATE "targets" SET "auth_key" = ?, "id" = ?, "path" = ?, "target_type" = ?  WHERE "targets"."id" = ?`,
			Args: []interface{}{
				newTarget.AuthKey,
				int64(newTarget.ID),
				newTarget.Path,
				newTarget.TargetType,
				int64(newTarget.ID),
			},
			Once: true,
		}

		mocket.Catcher.Attach([]*mocket.FakeResponse{targetUpdate})

		err = pts.targetRepo.Update(newTarget)
		pts.Require().NoError(err)
		pts.Require().True(targetUpdate.Triggered, "target update reference must be triggered")
	})

	pts.Run("Should bubble up any unhandled error", func() {
		pts.SetupTest()
		expectedError := stderrors.New("my error")
		newTarget := model.Target{ID: 10}
		err := pts.targetRepo.Create(&newTarget)
		pts.Require().NoError(err)
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern: `UPDATE "targets" SET "id" = ?  WHERE "targets"."id" = ?`,
				Args: []interface{}{
					int64(newTarget.ID),
					int64(newTarget.ID),
				},
				Once:  true,
				Error: expectedError,
			},
		})

		err = pts.targetRepo.Update(&newTarget)
		pts.Require().Error(err)
		pts.Require().EqualError(errors.Cause(err), expectedError.Error())
	})
}

func (pts *targetTestSuite) TestGormTargetDelete() {
	pts.Run("Should delete an existing target in database", func() {
		pts.SetupTest()

		targetDelete := &mocket.FakeResponse{
			Pattern: `DELETE FROM "targets"  WHERE "targets"."id" = ?`,
			Args: []interface{}{
				int64(1),
			},
			RowsAffected: 1,
			Once:         true,
		}

		mocket.Catcher.Attach([]*mocket.FakeResponse{targetDelete})

		err := pts.targetRepo.Delete(1)
		pts.Require().NoError(err)
		pts.Require().True(targetDelete.Triggered, "target delete reference must be triggered")
	})

	pts.Run("Should return EntityNotFound error when target is not found.", func() {
		pts.SetupTest()

		targetDelete := &mocket.FakeResponse{
			Pattern: `DELETE FROM "targets"  WHERE "targets"."id" = ?`,
			Args: []interface{}{
				int64(1),
			},
			RowsAffected: 0,
			Once:         true,
		}

		mocket.Catcher.Attach([]*mocket.FakeResponse{targetDelete})

		err := pts.targetRepo.Delete(1)
		pts.Require().Error(err)
		pts.Require().True(targetDelete.Triggered, "target delete reference must be triggered")
		pts.Require().EqualError(errors.Cause(err), repository.ErrEntityNotFound.Error())
	})

	pts.Run("Should bubble up any unhandled error", func() {
		pts.SetupTest()
		expectedError := stderrors.New("my error")
		newTarget := model.Target{ID: 10}
		err := pts.targetRepo.Create(&newTarget)
		pts.Require().NoError(err)
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern: `DELETE FROM "targets"  WHERE "targets"."id" = ?`,
				Args: []interface{}{
					int64(1),
				},
				Once:  true,
				Error: expectedError,
			},
		})

		err = pts.targetRepo.Delete(1)
		pts.Require().Error(err)
		pts.Require().EqualError(errors.Cause(err), expectedError.Error())
	})
}

func (pts *targetTestSuite) TestGormTargetAll() {
	pts.Run("Should return all targets", func() {
		pts.SetupTest()
		expectedTargets := []*model.Target{
			{
				ID: 123,
			},
			{
				ID: 321,
			},
		}
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  `SELECT * FROM "targets"`,
				Response: []map[string]interface{}{{"id": expectedTargets[0].ID}, {"id": expectedTargets[1].ID}},
				Once:     true,
			},
		})

		targets, err := pts.targetRepo.All()
		pts.Require().NoError(err)
		for i := range expectedTargets {
			pts.Assert().EqualValues(expectedTargets[i].ID, targets[i].ID)
		}
	})
	pts.Run("Should bubble up any unhandled error", func() {
		pts.SetupTest()
		expectedError := stderrors.New("my error")
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern: `SELECT * FROM "targets"`,
				Once:    true,
				Error:   expectedError,
			},
		})

		_, err := pts.targetRepo.All()
		pts.Require().Error(err)
		pts.Require().EqualError(errors.Cause(err), expectedError.Error())
	})
}
