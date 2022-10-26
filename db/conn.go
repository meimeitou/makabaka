package db

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/mattn/go-sqlite3"
	"github.com/meimeitou/makabaka/model/query"
	"gorm.io/gorm"
)

type Conn struct {
	config NetworkDB
	db     *gorm.DB
	q      *query.Query
}

func NewConn(db *gorm.DB, config NetworkDB) *Conn {
	q := query.Use(db)
	return &Conn{
		db:     db,
		config: config,
		q:      q,
	}
}

func (c *Conn) GetDB() *gorm.DB {
	return c.db
}

func (c *Conn) GetQuery() *query.Query {
	return c.q
}

func (c *Conn) Raw(sql string, values ...interface{}) *gorm.DB {
	return c.db.Raw(sql, values...)
}

func (c *Conn) ScanRows(rows *sql.Rows, dest interface{}) error {
	return c.db.ScanRows(rows, dest)
}

func (c *Conn) MustExec(sql string, values ...interface{}) {
	if err := c.db.Exec(sql, values...).Error; err != nil {
		panic(fmt.Sprintf("Error executing SQL %s: %v", sql, err))
	}
}

func (c *Conn) SourceSQL(filename string) {
	if c.config.Host != "localhost" &&
		c.config.Host != "127.0.0.1" {
		panic("only support init localhost db!")
	}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("Error reading %s: %v", filename, err))
	}
	sqlList := strings.Split(string(content), ";")
	for _, sql := range sqlList {
		cleaned := strings.TrimSpace(sql)
		if cleaned == "" {
			continue
		}
		c.MustExec(cleaned)
	}
}

func (c *Conn) AutoMigrate(dst ...interface{}) error {
	return c.db.Migrator().AutoMigrate(dst...)
}

func IsDuplicateError(err error, dbtype string) bool {
	if err != nil {
		switch dbtype {
		case "sqlite":
			if err, ok := err.(sqlite3.Error); ok {
				if err.ExtendedCode == sqlite3.ErrConstraintUnique {
					return true
				}
			}
		case "postgres":
			if err, ok := err.(*pgconn.PgError); ok {
				if err.Code == pgerrcode.UniqueViolation {
					return true
				}
			}
		}
	}
	return false
}

func TranslateErrors(err error, dbtype string) error {
	if err != nil {
		switch dbtype {
		case "sqlite":
			if err, ok := err.(sqlite3.Error); ok {
				if err.ExtendedCode == sqlite3.ErrConstraintUnique {
					return errors.New("record already exists")

				}
				// other errors handling from sqlite3
			}
		case "postgres":
			if err, ok := err.(*pgconn.PgError); ok {
				if err.Code == pgerrcode.UniqueViolation {
					return errors.New("record already exists")
				}
				// other errors handling from pgconn
			}
		}
	}
	return err
}
