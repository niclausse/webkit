package mocker

import (
	"fmt"
	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/niclausse/webkit/resource"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	dbDefaultUser     = "root"
	dbDefaultPassword = "antiy123"
	dbDefaultCharset  = "utf8mb4"
)

type MySQLMocker struct {
	Server  *server.Server
	Clients map[string]*gorm.DB
}

type MySQLConf struct {
	Server  *MySQLServerConf
	Clients []*MySQLClientConf
}

type MySQLServerConf struct {
	Addr      string
	Databases []string
}

type MySQLClientConf struct {
	Database string
	Addr     string
}

// InitServerAndClients start a mysql server and init mysql connection clients according to configs
func InitServerAndClients(configs *MySQLConf) *MySQLMocker {
	if configs == nil {
		panic(errors.New("empty configs"))
	}

	m := new(MySQLMocker)
	m.initMySQLServerAndDatabases(configs.Server)
	m.Clients = initMySQLClients(configs.Clients)
	return m
}

type TableOption func(db *memory.Database)

func (m *MySQLMocker) StartMySQLServer(addr string, databases ...sql.Database) {
	_ = sql.NewEmptyContext()
	engine := sqle.NewDefault(
		memory.NewDBProvider(
			databases...,
		),
	)

	_addr := strings.Split(addr, ":")
	if len(_addr) != 2 {
		log.Fatalf("invalid addr: %s", addr)
	}

	ed := engine.Analyzer.Catalog.MySQLDb.Editor()
	engine.Analyzer.Catalog.MySQLDb.AddSuperUser(ed, dbDefaultUser, _addr[0], dbDefaultPassword)
	ed.Close()

	config := server.Config{
		Protocol: "tcp",
		Address:  addr,
		Version:  "5.7.24-log",
	}
	var err error
	m.Server, err = server.NewDefaultServer(config, engine)
	if err != nil {
		panic(err)
	}

	if err = m.Server.Start(); err != nil {
		panic(err)
	}
}

func initMySQLClients(configs []*MySQLClientConf) (clients map[string]*gorm.DB) {
	clients = make(map[string]*gorm.DB, len(configs))

	for _, v := range configs {
		db, err := resource.InitMySQL(resource.MysqlConf{
			Service:  v.Database,
			Addr:     v.Addr,
			Database: v.Database,
			User:     dbDefaultUser,
			Password: dbDefaultPassword,
			Charset:  dbDefaultCharset,
		})
		if err != nil {
			panic(err)
		}
		clients[v.Database] = db
	}

	return
}

func (m *MySQLMocker) initMySQLServerAndDatabases(configs *MySQLServerConf) {
	memoryDBs := make([]sql.Database, 0, len(configs.Databases))
	for _, v := range configs.Databases {
		memoryDBs = append(memoryDBs, registerDB(v))
	}

	go func() {
		m.StartMySQLServer(configs.Addr, memoryDBs...)
	}()

	conn := make(chan struct{}, 1)
	ready := time.Now()
	tick := time.NewTicker(time.Millisecond * 100)
	defer tick.Stop()

	for {
		isPortOn(configs.Addr, conn)

		select {
		case now := <-tick.C:
			if sub := now.Sub(ready); sub > time.Second*120 {
				panic(fmt.Sprintf("failed to start test db server, cost%fs", sub.Seconds()))
			}
		case <-conn:
			return
		}
	}
}

func InitData(sqlFile string, db *gorm.DB) error {
	b, err := os.ReadFile(sqlFile)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, _sql := range strings.Split(strings.Trim(string(b), "\n"), ";") {
		if err = db.Exec(_sql).Error; err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func isPortOn(addr string, ch chan struct{}) {
	// 检测端口
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil || conn == nil {
		return
	}

	_ = conn.Close()
	ch <- struct{}{}
}

func registerDB(dbName string) *memory.Database {
	db := memory.NewDatabase(dbName)
	db.EnablePrimaryKeyIndexes()
	_ = db.SetCollation(sql.NewEmptyContext(), sql.Collation_utf8mb4_general_ci)
	return db
}
