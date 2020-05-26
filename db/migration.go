package db

import (
	"github.com/EurosportDigital/global-transcoding-platform/lib/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Migration interface {
	UpdateTables([]interface{}, *gorm.DB)
}

type GormMigration struct {
}

// UpdateTables initializes the schema and tables on the postres RDS DB.
func (instance *GormMigration) UpdateTables(models []interface{}, db *gorm.DB) {
	transaction := db.Begin()
	for _, m := range models {
		logger.Infof("Running migration for %T model...", m)
	}
	transaction.AutoMigrate(models...)
	transaction.Commit()
}
