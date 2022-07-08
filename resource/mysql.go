package resource

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type MysqlConf struct {
	Addr            string        `yaml:"addr"`
	Database        string        `yaml:"database"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Charset         string        `yaml:"charset"`
	MaxIdleConns    int           `yaml:"maxidleconns"`
	MaxOpenConns    int           `yaml:"maxopenconns"`
	ConnMaxIdlTime  time.Duration `yaml:"maxIdleTime"`
	ConnMaxLifeTime time.Duration `yaml:"connMaxLifeTime"`
	ConnTimeOut     time.Duration `yaml:"connTimeOut"`
	WriteTimeOut    time.Duration `yaml:"writeTimeOut"`
	ReadTimeOut     time.Duration `yaml:"readTimeOut"`
}

func (conf *MysqlConf) checkConf() {
	if conf.MaxIdleConns == 0 {
		conf.MaxIdleConns = 50
	}
	if conf.MaxOpenConns == 0 {
		conf.MaxOpenConns = 50
	}
	if conf.ConnMaxIdlTime == 0 {
		conf.ConnMaxLifeTime = 5 * time.Minute
	}
	if conf.ConnMaxLifeTime == 0 {
		conf.ConnMaxLifeTime = 10 * time.Minute
	}
	if conf.ConnTimeOut == 0 {
		conf.ConnTimeOut = 3 * time.Second
	}
	if conf.WriteTimeOut == 0 {
		conf.WriteTimeOut = 1200 * time.Millisecond
	}
	if conf.ReadTimeOut == 0 {
		conf.ReadTimeOut = 1200 * time.Millisecond
	}
}

func InitMySQL(conf MysqlConf, logger logger.Interface) (*gorm.DB, error) {
	conf.checkConf()

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=%s&readTimeout=%s&writeTimeout=%s&parseTime=True&loc=Asia%%2FShanghai",
		conf.User,
		conf.Password,
		conf.Addr,
		conf.Database,
		conf.ConnTimeOut,
		conf.ReadTimeOut,
		conf.WriteTimeOut,
	)

	if conf.Charset != "" {
		dsn = dsn + "&charset=" + conf.Charset
	}

	c := &gorm.Config{
		SkipDefaultTransaction:                   true,
		NamingStrategy:                           nil,
		FullSaveAssociations:                     false,
		Logger:                                   logger,
		NowFunc:                                  nil,
		DryRun:                                   false,
		PrepareStmt:                              false,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		AllowGlobalUpdate:                        false,
		ClauseBuilders:                           nil,
		ConnPool:                                 nil,
		Dialector:                                nil,
		Plugins:                                  nil,
	}

	client, err := gorm.Open(mysql.Open(dsn), c)
	if err != nil {
		return client, err
	}

	sqlDB, err := client.DB()
	if err != nil {
		return client, err
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)

	// SetMaxOpenConns 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)

	// SetConnMaxLifetime 设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(conf.ConnMaxLifeTime)

	// only for go version >= 1.15 设置最大空闲连接时间
	sqlDB.SetConnMaxIdleTime(conf.ConnMaxIdlTime)

	return client, nil
}
