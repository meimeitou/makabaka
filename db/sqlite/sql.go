package sqlite

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	sqlite3 "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
// _ db.Storage = new(SQLite3)
)

type SQLite3 struct {
	File string `json:"file"`
}

func (s *SQLite3) Open(logger *log.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite3.Open(s.File), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database")
	}
	return db, nil
}
