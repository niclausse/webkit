package mocker

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
	"time"
)

type DBMocker struct {
	db     *sql.DB
	gormDB *gorm.DB
	sqlmock.Sqlmock
}

func (m *DBMocker) GormDB() *gorm.DB {
	return m.gormDB
}

type DBMockerType = int8

const (
	DBMockerMySQL      DBMockerType = 0
	DBMockerPostgresql DBMockerType = 1
)

func Mock(t *testing.T, tp DBMockerType) *DBMocker {
	if tp != DBMockerMySQL && tp != DBMockerPostgresql {
		tp = DBMockerMySQL
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal(err)
	}

	var dial gorm.Dialector

	if tp == DBMockerPostgresql {
		dial = postgres.New(postgres.Config{
			DriverName:           "",
			DSN:                  "",
			PreferSimpleProtocol: false,
			WithoutReturning:     false,
			Conn:                 db,
		})
	} else {
		dial = mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		})
	}

	gormDB, err := gorm.Open(dial, &gorm.Config{
		SkipDefaultTransaction: false,
		NamingStrategy:         nil,
		FullSaveAssociations:   false,
		Logger:                 nil,
		NowFunc: func() time.Time {
			return time.Unix(time.Now().Local().Unix(), 0)
		},
		DryRun:                                   false,
		PrepareStmt:                              false,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		DisableNestedTransaction:                 false,
		AllowGlobalUpdate:                        false,
		QueryFields:                              false,
		CreateBatchSize:                          0,
		ClauseBuilders:                           nil,
		ConnPool:                                 db,
		Dialector:                                nil,
		Plugins:                                  nil,
	})
	if err != nil {
		t.Fatal(err)
	}

	return &DBMocker{
		db:      db,
		gormDB:  gormDB,
		Sqlmock: mock,
	}
}

func (m *DBMocker) Close() {
	_ = m.db.Close()
}
