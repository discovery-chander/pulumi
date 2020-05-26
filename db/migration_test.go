package db

import (
	"testing"

	"github.com/EurosportDigital/global-transcoding-platform/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/suite"
)

type UpdateTableSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *UpdateTableSuite) SetupTest() {
	mocket.Catcher.Register() // Safe register. Allowed multiple calls to save
	mocket.Catcher.Logging = true
	// GORM

	db, err := gorm.Open(mocket.DriverName, "connectionString") // Can be any connection string
	db.LogMode(true)
	if err != nil {
		panic(err)
	}
	suite.db = db
}

func (suite *UpdateTableSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *UpdateTableSuite) TestUpdateTable() {
	suite.Run("Create audio tracks table schema and tables with gorm automigration", func() {
		mocket.Catcher.Reset()
		createAudioTrackTableMock := mocket.Catcher.NewMock().WithQuery(`SELECT count(*) FROM INFORMATION_SCHEMA.TABLES WHERE table_schema =  AND table_name = audio_tracks`)
		createEncoderTableMock := mocket.Catcher.NewMock().WithQuery(`SELECT count(*) FROM INFORMATION_SCHEMA.TABLES WHERE table_schema =  AND table_name = encoders`)
		createEncoderConfigsTableMock := mocket.Catcher.NewMock().WithQuery(`SELECT count(*) FROM INFORMATION_SCHEMA.TABLES WHERE table_schema =  AND table_name = encoder_configs`)
		createJobRequestTableMock := mocket.Catcher.NewMock().WithQuery(`SELECT count(*) FROM INFORMATION_SCHEMA.TABLES WHERE table_schema =  AND table_name = jobs`)
		createOutputsTableMock := mocket.Catcher.NewMock().WithQuery(`SELECT count(*) FROM INFORMATION_SCHEMA.TABLES WHERE table_schema =  AND table_name = outputs`)
		createProfileTableMock := mocket.Catcher.NewMock().WithQuery(`SELECT count(*) FROM INFORMATION_SCHEMA.TABLES WHERE table_schema =  AND table_name = profiles`)
		createSubtitlesTableMock := mocket.Catcher.NewMock().WithQuery(`SELECT count(*) FROM INFORMATION_SCHEMA.TABLES WHERE table_schema =  AND table_name = subtitles`)
		createTargetsTableMock := mocket.Catcher.NewMock().WithQuery(`SELECT count(*) FROM INFORMATION_SCHEMA.TABLES WHERE table_schema =  AND table_name = targets`)
		var models []interface{}
		models = append(models, &model.AudioTrack{}, &model.Encoder{}, &model.EncoderConfig{}, &model.Job{}, &model.Output{}, &model.Profile{}, &model.Subtitle{}, &model.Target{})
		migration := &GormMigration{}
		migration.UpdateTables(models, suite.db)

		suite.Assert().True(createAudioTrackTableMock.Triggered, "Create audio tracks table query was called correctly")
		suite.Assert().True(createEncoderTableMock.Triggered, "Create encoder table query was called correctly")
		suite.Assert().True(createEncoderConfigsTableMock.Triggered, "Create encoder configs table query was called correctly")
		suite.Assert().True(createJobRequestTableMock.Triggered, "Create job request table query was called correctly")
		suite.Assert().True(createOutputsTableMock.Triggered, "Create outputs table query was called correctly")
		suite.Assert().True(createProfileTableMock.Triggered, "Create profile table query was called correctly")
		suite.Assert().True(createSubtitlesTableMock.Triggered, "Create subtitles table query was called correctly")
		suite.Assert().True(createTargetsTableMock.Triggered, "Create targets table query was called correctly")
	})
}

func TestUpdateTablesSuite(t *testing.T) {
	suite.Run(t, new(UpdateTableSuite))
}
