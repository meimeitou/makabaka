package postgres

import (
	"fmt"
	"testing"

	"github.com/meimeitou/makabaka/db"
	log "github.com/sirupsen/logrus"
)

func TestPg(t *testing.T) {
	st := PostgresStorage{
		NetworkDB: db.NetworkDB{
			Host:     "localhost",
			Port:     5432,
			Database: "makabaka",
			User:     "makabaka",
			Password: "makabaka",
			Debug:    true,
		},
	}
	conn, err := st.Open(log.StandardLogger())
	if err != nil {
		panic(err)
	}
	tb, err := NewPostgresTable(conn, st.Database)
	if err != nil {
		panic(err)
	}
	tables, err := tb.GetTables()
	if err != nil {
		panic(err)
	}
	for _, item := range tables {
		fmt.Println(item)
		if err := tb.GetColumnsOfTable(item); err != nil {
			panic(err)
		}
		fmt.Println(item.Columns)
	}
}
