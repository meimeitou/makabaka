package postgres

import (
	"fmt"
	"time"

	"github.com/meimeitou/makabaka/db"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
)

var (
	_ db.Storage = new(PostgresStorage)
)

type PostgresStorage struct {
	db.NetworkDB
	SSL db.SSL `json:"ssl" yaml:"ssl"`
}

func (p *PostgresStorage) Open(logger *log.Logger) (*db.Conn, error) {
	conn, err := p.open(logger)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (p *PostgresStorage) open(logger *log.Logger) (*db.Conn, error) {
	sslmode := "disable"
	if p.SSL.Mode != "" {
		sslmode = p.SSL.Mode
	}
	connParams := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Shanghai",
		p.Host, p.Port, p.User, p.Password, p.Database, sslmode)
	logger.Info(connParams)
	newLogger := gormlog.New(
		logger,
		gormlog.Config{
			SlowThreshold: 3 * time.Second, // 慢 SQL 阈值
			LogLevel:      gormlog.Silent,  // Log level
			Colorful:      false,           // 禁用彩色打印
		},
	)
	gormdb, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  connParams,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   newLogger,
	})

	if err != nil {
		return nil, err
	}
	sqlDB, err := gormdb.DB()
	if err != nil {
		return nil, err
	}
	// default connection pool set
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Second * 500)
	if p.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(p.MaxIdleConns)
	}
	if p.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(p.MaxOpenConns)
	}
	if p.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(p.ConnMaxLifetime))
	}
	if p.Debug {
		gormdb = gormdb.Debug()
	}
	return db.NewConn(gormdb, p.NetworkDB), nil
}
