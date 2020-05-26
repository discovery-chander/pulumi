package profile

import (
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/EurosportDigital/global-transcoding-platform/lib/repository"
	"github.com/EurosportDigital/global-transcoding-platform/model"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type profileTestSuite struct {
	suite.Suite
	db          *gorm.DB
	profileRepo Repository
}

func TestProfileTestSuite(t *testing.T) {
	mocket.Catcher.Register()
	mocket.Catcher.Logging = false

	db, err := gorm.Open(mocket.DriverName, "")
	require.NoError(t, err)
	pSuite := &profileTestSuite{
		db:          db,
		profileRepo: New(db),
	}
	suite.Run(t, pSuite)
}

func (pts *profileTestSuite) TearDownSuite() {
	pts.db.Close()
}

func (pts *profileTestSuite) SetupTest() {
	mocket.Catcher.Reset()
}

func (pts *profileTestSuite) TestGormProfileGet() {
	getQuery := `SELECT * FROM "profiles"  WHERE ("profiles"."id" = 1) ORDER BY "profiles"."id" ASC LIMIT 1`
	pts.Run("Should return expected result", func() {
		pts.SetupTest()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  getQuery,
				Response: []map[string]interface{}{{"id": 1}},
				Once:     true,
			},
		})

		profile, err := pts.profileRepo.Get(1)
		pts.Require().NoError(err)
		pts.Require().EqualValues(1, profile.ID)
	})
	pts.Run("Should return repository.ErrEntityNotFound if profile is not found", func() {
		pts.SetupTest()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  getQuery,
				Response: []map[string]interface{}{},
				Once:     true,
			},
		})

		_, err := pts.profileRepo.Get(1)
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

		_, err := pts.profileRepo.Get(1)
		pts.Require().Error(err)
		pts.Require().EqualError(errors.Cause(err), expectedError.Error())
	})
}

func (pts *profileTestSuite) TestGormProfileGetMany() {
	pts.Run("Should return expected result", func() {
		pts.SetupTest()
		expectedIDs := []int{1, 2}
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  `SELECT * FROM "profiles"  WHERE (id IN (1,2))`,
				Args:     []interface{}{int64(expectedIDs[0]), int64(expectedIDs[1])},
				Response: []map[string]interface{}{{"id": expectedIDs[0]}, {"id": expectedIDs[1]}},
				Once:     true,
			},
		})

		profiles, err := pts.profileRepo.GetMany(expectedIDs)
		pts.Require().NoError(err)
		pts.Require().Len(profiles, len(expectedIDs))
		for i := range profiles {
			pts.Require().EqualValues(expectedIDs[i], profiles[i].ID)
		}
	})
	pts.Run("Should bubble up any unhandled error", func() {
		pts.SetupTest()
		expectedError := stderrors.New("my error")
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern: `SELECT * FROM "profiles"  WHERE (id IN (1))`,
				Args:    []interface{}{int64(1)},
				Once:    true,
				Error:   expectedError,
			},
		})

		_, err := pts.profileRepo.GetMany([]int{1})
		pts.Require().Error(err)
		pts.Require().EqualError(errors.Cause(err), expectedError.Error())
	})
}

func (pts *profileTestSuite) TestGormProfileGetManyByName() {
	pts.Run("Should return expected result", func() {
		pts.SetupTest()
		expectedNames := []string{"name 1", "name 2"}
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  `SELECT * FROM "profiles"  WHERE (name IN (name 1,name 2))`,
				Args:     []interface{}{expectedNames[0], expectedNames[1]},
				Response: []map[string]interface{}{{"name": expectedNames[0]}, {"name": expectedNames[1]}},
				Once:     true,
			},
		})

		profiles, err := pts.profileRepo.GetManyByName(expectedNames)
		pts.Require().NoError(err)
		pts.Require().Len(profiles, len(expectedNames))
		for i := range profiles {
			pts.Require().EqualValues(expectedNames[i], profiles[i].Name)
		}
	})
	pts.Run("Should bubble up any unhandled error", func() {
		pts.SetupTest()
		expectedError := stderrors.New("my error")
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern: `SELECT * FROM "profiles"  WHERE (name IN (name))`,
				Args:    []interface{}{"name"},
				Once:    true,
				Error:   expectedError,
			},
		})

		_, err := pts.profileRepo.GetManyByName([]string{"name"})
		pts.Require().Error(err)
		pts.Require().EqualError(errors.Cause(err), expectedError.Error())
	})
}

func (pts *profileTestSuite) TestGormProfileGetByName() {
	expectedName := "my profile name"
	getQuery := fmt.Sprintf(`SELECT * FROM "profiles"  WHERE (name = %[1]s) ORDER BY "profiles"."id" ASC LIMIT 1`, expectedName)
	pts.Run("Should return expected result", func() {
		pts.SetupTest()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  getQuery,
				Response: []map[string]interface{}{{"name": expectedName}},
				Once:     true,
			},
		})

		profile, err := pts.profileRepo.GetByName(expectedName)
		pts.Require().NoError(err)
		pts.Require().EqualValues(expectedName, profile.Name)
	})
	pts.Run("Should return repository.ErrEntityNotFound if profile is not found", func() {
		pts.SetupTest()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  getQuery,
				Response: []map[string]interface{}{},
				Once:     true,
			},
		})

		_, err := pts.profileRepo.GetByName(expectedName)
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

		_, err := pts.profileRepo.GetByName(expectedName)
		pts.Require().Error(err)
		pts.Require().EqualError(errors.Cause(err), expectedError.Error())
	})
}

func (pts *profileTestSuite) TestGormProfileCreate() {
	pts.Run("Should create new profile in database", func() {
		pts.SetupTest()
		newProfile := &model.Profile{
			Name:          "my name",
			Codec:         "my codec",
			PackageFormat: "format",
			EncConfig: model.EncoderConfig{
				Name:   "my encoder config",
				Config: "what is a config",
				Encoder: model.Encoder{
					Name:        "my encoder",
					ApiEndpoint: "google.com",
					InfoUrl:     "stillgoogle.com",
				},
			},
		}

		profileInsert := &mocket.FakeResponse{
			Pattern: `INSERT INTO "profiles" ("name","codec","package_format","encoder_config_id") VALUES (?,?,?,?)`,
			Args: []interface{}{
				newProfile.Name,
				newProfile.Codec,
				newProfile.PackageFormat,
				int64(1),
			},
			LastInsertID: int64(1),
			Once:         true,
		}
		encoderConfigInsert := &mocket.FakeResponse{
			Pattern: `INSERT INTO "encoder_configs" ("name","config","encoder_id") VALUES (?,?,?)`,
			Args: []interface{}{
				newProfile.EncConfig.Name,
				newProfile.EncConfig.Config,
				int64(1),
			},
			LastInsertID: int64(1),
			Once:         true,
		}
		encoderInsert := &mocket.FakeResponse{
			Pattern: `INSERT INTO "encoders" ("name","api_endpoint","info_url") VALUES (?,?,?)`,
			Args: []interface{}{
				newProfile.EncConfig.Encoder.Name,
				newProfile.EncConfig.Encoder.ApiEndpoint,
				newProfile.EncConfig.Encoder.InfoUrl,
			},
			LastInsertID: int64(1),
			Once:         true,
		}
		mocket.Catcher.Attach([]*mocket.FakeResponse{encoderInsert, encoderConfigInsert, profileInsert})

		err := pts.profileRepo.Create(newProfile)
		pts.Require().NoError(err)
		pts.Require().True(profileInsert.Triggered, "profile insert reference must be triggered")
		pts.Require().True(encoderConfigInsert.Triggered, "encoder config insert reference must be triggered")
		pts.Require().True(encoderInsert.Triggered, "encoder insert reference must be triggered")
	})

	pts.Run("Should bubble up any unhandled error", func() {
		pts.SetupTest()
		expectedError := stderrors.New("my error")
		newProfile := model.Profile{EncConfig: model.EncoderConfig{ID: 1}}
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern: `INSERT INTO "profiles" ("name","codec","package_format","encoder_config_id") VALUES (?,?,?,?)`,
				Args: []interface{}{
					newProfile.Name,
					newProfile.Codec,
					newProfile.PackageFormat,
					int64(newProfile.EncConfig.ID),
				},
				Once:  true,
				Error: expectedError,
			},
		})

		err := pts.profileRepo.Create(&newProfile)
		pts.Require().Error(err)
		pts.Require().EqualError(errors.Cause(err), expectedError.Error())
	})
}

func (pts *profileTestSuite) TestGormProfileUpdate() {
	pts.Run("Should update an existing profile in database", func() {
		pts.SetupTest()
		newProfile := &model.Profile{
			Name:          "my name",
			Codec:         "my codec",
			PackageFormat: "format",
			EncConfig: model.EncoderConfig{
				Name:   "my encoder config",
				Config: "what is a config",
				Encoder: model.Encoder{
					Name:        "my encoder",
					ApiEndpoint: "google.com",
					InfoUrl:     "stillgoogle.com",
				},
			},
		}

		err := pts.profileRepo.Create(newProfile)
		pts.Require().NoError(err)
		newProfile.Name = "updated name"

		profileUpdate := &mocket.FakeResponse{
			Pattern: `UPDATE "profiles" SET "codec" = ?, "encoder_config_id" = ?, "id" = ?, "name" = ?, "package_format" = ?  WHERE "profiles"."id" = ?`,
			Args: []interface{}{
				newProfile.Codec,
				int64(newProfile.EncConfig.ID),
				int64(newProfile.ID),
				newProfile.Name,
				newProfile.PackageFormat,
				int64(newProfile.ID),
			},
			Once: true,
		}

		mocket.Catcher.Attach([]*mocket.FakeResponse{profileUpdate})

		err = pts.profileRepo.Update(newProfile)
		pts.Require().NoError(err)
		pts.Require().True(profileUpdate.Triggered, "profile update reference must be triggered")
	})

	pts.Run("Should bubble up any unhandled error", func() {
		pts.SetupTest()
		expectedError := stderrors.New("my error")
		newProfile := model.Profile{EncConfig: model.EncoderConfig{}}
		err := pts.profileRepo.Create(&newProfile)
		pts.Require().NoError(err)
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern: `UPDATE "profiles" SET "encoder_config_id" = ?, "id" = ?  WHERE "profiles"."id" = ?`,
				Args: []interface{}{
					int64(newProfile.EncConfig.ID),
					int64(newProfile.ID),
					int64(newProfile.ID),
				},
				Once:  true,
				Error: expectedError,
			},
		})

		err = pts.profileRepo.Update(&newProfile)
		pts.Require().Error(err)
		pts.Require().EqualError(errors.Cause(err), expectedError.Error())
	})
}

func (pts *profileTestSuite) TestGormProfileDelete() {
	pts.Run("Should delete an existing profile in database", func() {
		pts.SetupTest()

		profileDelete := &mocket.FakeResponse{
			Pattern: `DELETE FROM "profiles"  WHERE "profiles"."id" = ?`,
			Args: []interface{}{
				int64(1),
			},
			RowsAffected: 1,
			Once:         true,
		}

		mocket.Catcher.Attach([]*mocket.FakeResponse{profileDelete})

		err := pts.profileRepo.Delete(1)
		pts.Require().NoError(err)
		pts.Require().True(profileDelete.Triggered, "profile delete reference must be triggered")
	})

	pts.Run("Should return EntityNotFound error when profile is not found.", func() {
		pts.SetupTest()

		profileDelete := &mocket.FakeResponse{
			Pattern: `DELETE FROM "profiles"  WHERE "profiles"."id" = ?`,
			Args: []interface{}{
				int64(1),
			},
			RowsAffected: 0,
			Once:         true,
		}

		mocket.Catcher.Attach([]*mocket.FakeResponse{profileDelete})

		err := pts.profileRepo.Delete(1)
		pts.Require().Error(err)
		pts.Require().True(profileDelete.Triggered, "profile delete reference must be triggered")
		pts.Require().EqualError(errors.Cause(err), repository.ErrEntityNotFound.Error())
	})

	pts.Run("Should bubble up any unhandled error", func() {
		pts.SetupTest()
		expectedError := stderrors.New("my error")
		newProfile := model.Profile{ID: 10, EncConfig: model.EncoderConfig{ID: 1}}
		err := pts.profileRepo.Create(&newProfile)
		pts.Require().NoError(err)
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern: `DELETE FROM "profiles"  WHERE "profiles"."id" = ?`,
				Args: []interface{}{
					int64(1),
				},
				Once:  true,
				Error: expectedError,
			},
		})

		err = pts.profileRepo.Delete(1)
		pts.Require().Error(err)
		pts.Require().EqualError(errors.Cause(err), expectedError.Error())
	})
}

func (pts *profileTestSuite) TestGormProfileAll() {
	pts.Run("Should return all profiles", func() {
		pts.SetupTest()
		expectedProfiles := []*model.Profile{
			{
				ID: 123,
			},
			{
				ID: 321,
			},
		}
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  `SELECT * FROM "profiles"`,
				Response: []map[string]interface{}{{"id": expectedProfiles[0].ID}, {"id": expectedProfiles[1].ID}},
				Once:     true,
			},
		})

		profiles, err := pts.profileRepo.All()
		pts.Require().NoError(err)
		for i := range expectedProfiles {
			pts.Assert().EqualValues(expectedProfiles[i].ID, profiles[i].ID)
		}
	})
	pts.Run("Should bubble up any unhandled error", func() {
		pts.SetupTest()
		expectedError := stderrors.New("my error")
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern: `SELECT * FROM "profiles"`,
				Once:    true,
				Error:   expectedError,
			},
		})

		_, err := pts.profileRepo.All()
		pts.Require().Error(err)
		pts.Require().EqualError(errors.Cause(err), expectedError.Error())
	})
}
