package test

import (
	"fmt"
	"testing"

	"github.com/meimeitou/makabaka/db"
	"github.com/meimeitou/makabaka/db/postgres"
	"github.com/meimeitou/makabaka/model/exec"
	log "github.com/sirupsen/logrus"
)

func TestSql(t *testing.T) {
	pg := postgres.PostgresStorage{
		NetworkDB: db.NetworkDB{
			Host:     "localhost",
			Port:     5432,
			Database: "makabaka",
			User:     "makabaka",
			Password: "makabaka",
		},
	}
	storage, err := pg.Open(log.StandardLogger())
	if err != nil {
		panic(err)
	}
	sql := `select * from apis
	where api_name=@name
	{{ if gt .ct 0 }} limit {{ .ct }} {{ end }}`
	query := exec.NewQueryBuilder(storage, sql, map[string]interface{}{"ct": 1, "name": "listApi"})
	data, err := query.Exec()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
	query = exec.NewQueryBuilder(storage, "update apis set method=@method", map[string]interface{}{"method": "POST"})
	data, err = query.Exec()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func TestTemplateParse(t *testing.T) {
	sql := `select api_name,method,sql_type from apis where api_name like @name {{ if gt (toType .limit "int") 0 }} limit {{ .limit }} {{ end }}`
	query := exec.NewQueryBuilder(nil, sql, map[string]interface{}{"limit": 1.1, "name": "test"})
	fmt.Println(query.TemplateParse())
}
