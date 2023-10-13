package mocker

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestDBMocker(t *testing.T) {
	mocker := Mock(t, DBMockerPostgresql)
	defer mocker.Close()

	user := &demoUser{
		Model: gorm.Model{},
		Name:  "Alice",
		Age:   18,
	}

	now := time.Unix(time.Now().Unix(), 0)

	_sql := `INSERT INTO "demo_user" ("created_at","updated_at","deleted_at","name","age") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`
	mocker.ExpectBegin()
	mocker.ExpectQuery(_sql).WithArgs(now, now, sql.NullTime{}, user.Name, user.Age).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1001))
	mocker.ExpectCommit()

	if err := create(mocker.GormDB(), user); err != nil {
		t.Errorf("failed to create user: %+v", err)
		return
	}

	if user.ID != 1001 {
		t.Errorf("user.id is not expected: %d", user.ID)
		return
	}
}

type demoUser struct {
	gorm.Model
	Name string
	Age  uint32
}

func create(db *gorm.DB, user *demoUser) error {
	return db.Table("demo_user").Create(user).Error
}
