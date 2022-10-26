package db

import (
	"database/sql"

	"gorm.io/gorm"
)

type DBType string

const (
	DBTypePostgresql DBType = "pg"
	DBTypeMySQL      DBType = "mysql"
	DBTypeSQLite     DBType = "sqlite3"
)

var (
	// dbTypeToDriverMap maps the database type to the driver names.
	dbTypeToDriverMap = map[DBType]string{
		DBTypePostgresql: "postgres",
		DBTypeMySQL:      "mysql",
		DBTypeSQLite:     "sqlite3",
	}
)

// Database interface for the concrete databases.
type Database interface {
	GetTables() (tables []*Table, err error)
	// PrepareGetColumnsOfTableStmt() (err error)
	GetColumnsOfTable(table *Table) (err error)

	IsPrimaryKey(column Column) bool
	IsAutoIncrement(column Column) bool
	IsNullable(column Column) bool

	GetStringDatatypes() []string
	IsString(column Column) bool

	GetTextDatatypes() []string
	IsText(column Column) bool

	GetIntegerDatatypes() []string
	IsInteger(column Column) bool

	GetFloatDatatypes() []string
	IsFloat(column Column) bool

	GetTemporalDatatypes() []string
	IsTemporal(column Column) bool
}

// Table has a name and a set (slice) of columns.
type Table struct {
	Name    string   `json:"table_name" gorm:"column:table_name"`
	Columns []Column `json:"columns" gorm:"-"`
}

func (t *Table) ToString() string {
	return t.Name
}

// Column stores information about a column.
type Column struct {
	OrdinalPosition        int            `json:"ordinal_position" gorm:"column:ordinal_position"`
	Name                   string         `json:"column_name" gorm:"column:column_name"`
	DataType               string         `json:"data_type" gorm:"column:data_type"`
	DefaultValue           sql.NullString `json:"column_default" gorm:"column:column_default"`
	IsNullable             string         `json:"is_nullable" gorm:"column:is_nullable"`
	CharacterMaximumLength sql.NullInt64  `json:"character_maximum_length" gorm:"column:character_maximum_length"`
	NumericPrecision       sql.NullInt64  `json:"numeric_precision" gorm:"column:numeric_precision"`
	ColumnKey              string         `json:"column_key" gorm:"column:column_key"`           // mysql specific
	Extra                  string         `json:"extra" gorm:"column:extra"`                     // mysql specific
	ConstraintName         sql.NullString `json:"constraint_name" gorm:"column:constraint_name"` // pg specific
	ConstraintType         sql.NullString `json:"constraint_type" gorm:"column:constraint_type"` // pg specific
}

// GeneralDatabase represents a base "class" database - for all other concrete
// databases it implements partly the Database interface.
type GeneralDatabase struct {
	db *gorm.DB
}

// New creates a new Database based on the given type in the settings.
func New() Database {

	// var db Database

	// switch s.DbType {
	// case settings.DBTypeSQLite:
	// 	db = NewSQLite(s)
	// case settings.DBTypeMySQL:
	// 	db = NewMySQL(s)
	// case settings.DBTypePostgresql:
	// 	fallthrough
	// default:
	// 	db = NewPostgresql(s)
	// }

	return nil
}

// Connect establishes a connection to the database with the given DSN.
// It pings the database to ensure it is reachable.
func (gdb *GeneralDatabase) Connect(dsn string) (err error) {
	return
}

// Close closes the database connection.
func (gdb *GeneralDatabase) Close() error {
	return nil
}

// IsNullable returns true if the column is a nullable column.
func (gdb *GeneralDatabase) IsNullable(column Column) bool {
	return column.IsNullable == "YES"
}

// isStringInSlice checks if needle (string) is in haystack ([]string).
func IsStringInSlice(needle string, haystack []string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}
