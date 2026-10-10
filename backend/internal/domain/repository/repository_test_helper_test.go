package repository

import (
	"testing"

	logger "backend/internal/log"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// setupMockDB returns a GORM DB backed by sqlmock, so repository tests run
// the real GORM MySQL query building without a database.
func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	logger.InitForTest()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}

	dialector := mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm DB: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return gormDB, mock, cleanup
}
