package exec

import (
	"database/sql"
	"strings"

	"github.com/meimeitou/makabaka/db"
	"github.com/meimeitou/makabaka/pkg/tpl"
	"github.com/sirupsen/logrus"
)

type QueryBuilder struct {
	format map[string]string
	query  string
	values map[string]interface{}
	db     *db.Conn
	rawsql string
}

func NewQueryBuilder(database *db.Conn, query string, values map[string]interface{}) *QueryBuilder {
	return &QueryBuilder{
		query:  query,
		db:     database,
		values: values,
	}
}

func (q *QueryBuilder) GetRawSql() string {
	return q.rawsql
}

func (q *QueryBuilder) WithFormat(format map[string]string) {
	q.format = format
}

func (q *QueryBuilder) TemplateParse() (string, error) {
	logrus.Debug(q.query, q.values)
	return tpl.Parse(q.query, q.values)
}

func (q *QueryBuilder) Exec() ([]map[string]interface{}, error) {
	//Use raw query
	var (
		err  error
		rows *sql.Rows
	)
	query, err := q.TemplateParse()
	if err != nil {
		return nil, err
	}
	q.rawsql = query
	if strings.Contains(query, "@") && len(q.values) > 0 {
		rows, err = q.db.Raw(query, q.values).Rows()
	} else {
		rows, err = q.db.Raw(query).Rows()
	}
	if err != nil {
		return nil, err
	}
	var list []map[string]interface{}
	for rows.Next() {
		// ScanRows scan a row into user
		data := make(map[string]interface{})
		if err := q.db.ScanRows(rows, &data); err != nil {
			return nil, err
		}
		list = append(list, data)
	}

	return list, nil
}
