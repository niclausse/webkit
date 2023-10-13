package mocker

import (
	"testing"
)

func TestMySQL(t *testing.T) {
	dbMocker := InitServerAndClients(&MySQLConf{
		Server: &MySQLServerConf{
			Addr:      "localhost:3306",
			Databases: []string{"test"},
		},
		Clients: []*MySQLClientConf{
			{
				Database: "test",
				Addr:     "localhost:3306",
			},
		},
	})
	defer dbMocker.Server.Listener.Shutdown()

	if err := InitData("test.sql", dbMocker.Clients["test"]); err != nil {
		t.Errorf("%+v", err)
		return
	}

	user := struct {
		ID   int64
		Name string
		Age  int32
		Sex  uint8
	}{}
	if err := dbMocker.Clients["test"].Table("tbl_user").First(&user).Error; err != nil {
		t.Errorf("%+v", err)
		return
	}

	if user.ID != 1 {
		t.Errorf("user.id should be %d, but %d", 1, user.ID)
		return
	}
	if user.Name != "lp" {
		t.Errorf("user.id should be lp, but %s", user.Name)
		return
	}
	if user.Sex != 1 {
		t.Errorf("user.id should be %d, but %d", 1, user.Sex)
		return
	}
	if user.Age != 18 {
		t.Errorf("user.id should be %d, but %d", 18, user.Age)
		return
	}
}
