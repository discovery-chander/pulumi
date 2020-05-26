package db

import (
	"github.com/EurosportDigital/global-transcoding-platform/lib/errors"
	"github.com/EurosportDigital/global-transcoding-platform/lib/logger"
	"github.com/jinzhu/gorm"
)

// OpenDBConnection open the connection to the RDS instance using GORM.
func OpenDBConnection(connectionString string, connectionDriver string) (*gorm.DB, error) {
	if len(connectionDriver) == 0 {
		connectionDriver = "postgres"
	}
	if len(connectionString) == 0 {
		return nil, errors.New("Missing arguments to open DB connection")
	}
	db, dbOpenError := gorm.Open(connectionDriver, connectionString)
	if dbOpenError != nil {
		logger.Error(dbOpenError, "Error when opening DB connection")
		return nil, errors.Wrap(dbOpenError, "opening DB connection")
	}
	return db, nil
}
