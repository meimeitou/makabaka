package postgres

import (
	"errors"
	"strings"

	"github.com/meimeitou/makabaka/db"
)

var (
	_ db.Database = new(PostgresTable)
)

type PostgresTable struct {
	dbname string
	db     *db.Conn
}

func NewPostgresTable(conn *db.Conn, dbname string) (*PostgresTable, error) {
	if conn == nil {
		return nil, errors.New("db not init")
	}
	return &PostgresTable{
		dbname: dbname,
		db:     conn,
	}, nil
}

// GetTables gets all tables for a given schema by name.
func (pg *PostgresTable) GetTables() (tables []*db.Table, err error) {
	err = pg.db.Raw(`SELECT table_name
	FROM information_schema.tables
	WHERE table_type = 'BASE TABLE'
	AND table_schema = ?
	ORDER BY table_name`, pg.dbname).Scan(&tables).Error
	return
}

func (pg *PostgresTable) Connect() (err error) {
	return nil
}

// func (pg *Postgres) PrepareGetColumnsOfTableStmt() (err error) {
// 	return nil
// }

// get table columns
func (pg *PostgresTable) GetColumnsOfTable(table *db.Table) (err error) {
	pg.db.Raw(`
	SELECT
		ic.ordinal_position,
		ic.column_name,
		ic.data_type,
		ic.column_default,
		ic.is_nullable,
		ic.character_maximum_length,
		ic.numeric_precision,
		itc.constraint_name,
		itc.constraint_type
	FROM information_schema.columns AS ic
		LEFT JOIN information_schema.key_column_usage AS ikcu ON ic.table_name = ikcu.table_name
		AND ic.table_schema = ikcu.table_schema
		AND ic.column_name = ikcu.column_name
		LEFT JOIN information_schema.table_constraints AS itc ON ic.table_name = itc.table_name
		AND ic.table_schema = itc.table_schema
		AND ikcu.constraint_name = itc.constraint_name
	WHERE ic.table_name = ?
	AND ic.table_schema = ?
	ORDER BY ic.ordinal_position
`, table.Name, pg.dbname).Scan(&table.Columns)
	return nil
}

func (pg *PostgresTable) IsPrimaryKey(column db.Column) bool {
	return strings.Contains(column.ConstraintType.String, "PRIMARY KEY")
}

func (pg *PostgresTable) IsAutoIncrement(column db.Column) bool {
	return strings.Contains(column.DefaultValue.String, "nextval")
}

func (pg *PostgresTable) IsNullable(column db.Column) bool {
	return false
}

func (pg *PostgresTable) GetStringDatatypes() []string {
	return []string{
		"character varying",
		"varchar",
		"character",
		"char",
		"uuid",
	}
}
func (pg *PostgresTable) IsString(column db.Column) bool {
	return db.IsStringInSlice(column.DataType, pg.GetStringDatatypes())
}

func (pg *PostgresTable) GetTextDatatypes() []string {
	return []string{
		"text",
	}
}
func (pg *PostgresTable) IsText(column db.Column) bool {
	return db.IsStringInSlice(column.DataType, pg.GetTextDatatypes())
}

func (pg *PostgresTable) GetIntegerDatatypes() []string {
	return []string{
		"smallint",
		"integer",
		"bigint",
		"smallserial",
		"serial",
		"bigserial",
	}
}

func (pg *PostgresTable) IsInteger(column db.Column) bool {
	return db.IsStringInSlice(column.DataType, pg.GetIntegerDatatypes())
}

func (pg *PostgresTable) GetFloatDatatypes() []string {
	return []string{
		"numeric",
		"decimal",
		"real",
		"double precision",
	}
}
func (pg *PostgresTable) IsFloat(column db.Column) bool {
	return db.IsStringInSlice(column.DataType, pg.GetFloatDatatypes())
}

func (pg *PostgresTable) GetTemporalDatatypes() []string {
	return []string{
		"time",
		"timestamp",
		"time with time zone",
		"timestamp with time zone",
		"time without time zone",
		"timestamp without time zone",
		"date",
	}
}

func (pg *PostgresTable) IsTemporal(column db.Column) bool {
	return db.IsStringInSlice(column.DataType, pg.GetTemporalDatatypes())
}
