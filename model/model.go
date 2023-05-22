package model

import (
	"errors"
	"time"

	"gorm.io/plugin/soft_delete"
)

type SqlType int8

const (
	SqlTypeTemplate SqlType = iota + 1 // template
	SqlTypeChain                       // chain
)

func SqlTypeFromStr(tp string) (SqlType, error) {
	switch tp {
	case "template", "":
		return SqlTypeTemplate, nil
	case "chain":
		return SqlTypeChain, nil
	default:
		return SqlTypeTemplate, errors.New("error sql type")
	}
}

type ApiType string

type DbType string

const (
	DbTypePostgres DbType = "postgres"
	// TODO: support mysql
	// DbTypeMysql    DbType = "mysql"
)

// TODO: db config
type Database struct {
	Modelx
	DeletedAt soft_delete.DeletedAt `json:"deleteAt" gorm:"uniqueIndex:idx_db_key"`
	DBKey     string                `json:"dbKey" gorm:"column:db_key;type:varchar(200);uniqueIndex:idx_db_key;not null"` // db key desc
	Type      DbType                `json:"dbType" gorm:"column:db_type;not null;index"`
	Host      string                `json:"host" gorm:"column:host;type:varchar(100);not null"`
	Port      int                   `json:"port" gorm:"column:port;not null"`
	User      string                `json:"user" gorm:"column:user;not null"`
	Password  string                `json:"password" gorm:"column:password;not null"`
	DBName    string                `json:"dbName" gorm:"column:db_name;not null"`
}

// base sql template or sql chain
type Apis struct {
	Modelx
	DeletedAt             soft_delete.DeletedAt `json:"deleteAt" gorm:"uniqueIndex:idx_api_name"`
	ApiType               ApiType               `json:"apiType" gorm:"column:api_type;default:'';not null;index"`
	Name                  string                `json:"name" gorm:"column:api_name;type:varchar(200);uniqueIndex:idx_api_name;not null"` // api name
	Method                string                `json:"method" gorm:"column:method;type:varchar(50);not null;index"`                     // api name// http method
	Description           string                `json:"description" gorm:"column:description;type:varchar(200);not null"`                // description
	SqlType               SqlType               `json:"sqlType" gorm:"column:sql_type;not null;index"`
	SqlTemplate           JSONMap               `json:"sqlTemplate" gorm:"column:sql_template;not null"`                      // sql template
	SqlTemplateParameters JSONMap               `json:"sqlTemplateParameters" gorm:"column:sql_template_parameters;not null"` // template parameters
	SqlTemplateResult     JSONMap               `json:"sqlTemplateResult" gorm:"column:sql_template_result;not null"`         // format result
	// SqlChainParameters    string                `json:"sqlChainParameters"`
}

// soft delete
type Modelx struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt" gorm:"<-:create"`
	UpdatedAt time.Time `json:"deletedAt"`
}

func (m *Modelx) GetID() uint {
	return m.ID
}
