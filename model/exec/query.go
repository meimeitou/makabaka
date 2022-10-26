package exec

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
	"text/template"

	"github.com/meimeitou/makabaka/db"
	"github.com/sirupsen/logrus"
)

type QueryBuilder struct {
	format map[string]string
	query  string
	values map[string]interface{}
	db     *db.Conn
}

func NewQueryBuilder(database *db.Conn, query string, values map[string]interface{}) *QueryBuilder {
	return &QueryBuilder{
		query:  query,
		db:     database,
		values: values,
	}
}

func (q *QueryBuilder) WithFormat(format map[string]string) {
	q.format = format
}

func (q *QueryBuilder) TemplateParse() (string, error) {
	b := &bytes.Buffer{}
	t := template.New("").Funcs(toType)
	// t.Funcs()
	logrus.Debug(q.query, q.values)
	tpl, err := t.Parse(q.query)
	if err != nil {
		return "", fmt.Errorf("模板格式错误: %v", err)
	}
	err = tpl.Execute(b, q.values)
	if err != nil {
		return "", fmt.Errorf("查询格式验证错误%v", err)
	}
	return b.String(), nil
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
	if strings.Contains(query, "@") && len(q.values) > 0 {
		rows, err = q.db.Raw(query, q.values).Rows()
	} else {
		rows, err = q.db.Raw(query).Rows()
	}
	if err != nil {
		return nil, err
	}
	// columns, er := rows.ColumnTypes()
	// if er != nil {
	// 	log.Println(er)
	// 	return nil, er
	// }

	//make一個臨時儲存的地方，並賦予指標
	// cache := make([]interface{}, columnLength)
	// for index := range cache {
	// 	var a interface{}
	// 	cache[index] = &a
	// }

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

// func (q *Query) FormatRow(data *map[string]interface{}) {
// 	for key, formats := range q.format {
// 		if value, ok := (*data)[key]; ok {
// 			fm := strings.Split(formats, ",")
// 		}
// 	}
// }
