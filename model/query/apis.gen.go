// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/meimeitou/makabaka/model"
)

func newApis(db *gorm.DB) apis {
	_apis := apis{}

	_apis.apisDo.UseDB(db)
	_apis.apisDo.UseModel(&model.Apis{})

	tableName := _apis.apisDo.TableName()
	_apis.ALL = field.NewAsterisk(tableName)
	_apis.ID = field.NewUint(tableName, "id")
	_apis.CreatedAt = field.NewTime(tableName, "created_at")
	_apis.UpdatedAt = field.NewTime(tableName, "updated_at")
	_apis.DeletedAt = field.NewUint(tableName, "deleted_at")
	_apis.Name = field.NewString(tableName, "api_name")
	_apis.Method = field.NewString(tableName, "method")
	_apis.Description = field.NewString(tableName, "description")
	_apis.SqlType = field.NewInt8(tableName, "sql_type")
	_apis.SqlTemplate = field.NewField(tableName, "sql_template")
	_apis.SqlTemplateParameters = field.NewField(tableName, "sql_template_parameters")
	_apis.SqlTemplateResult = field.NewField(tableName, "sql_template_result")

	_apis.fillFieldMap()

	return _apis
}

type apis struct {
	apisDo

	ALL                   field.Asterisk
	ID                    field.Uint
	CreatedAt             field.Time
	UpdatedAt             field.Time
	DeletedAt             field.Uint
	Name                  field.String
	Method                field.String
	Description           field.String
	SqlType               field.Int8
	SqlTemplate           field.Field
	SqlTemplateParameters field.Field
	SqlTemplateResult     field.Field

	fieldMap map[string]field.Expr
}

func (a apis) Table(newTableName string) *apis {
	a.apisDo.UseTable(newTableName)
	return a.updateTableName(newTableName)
}

func (a apis) As(alias string) *apis {
	a.apisDo.DO = *(a.apisDo.As(alias).(*gen.DO))
	return a.updateTableName(alias)
}

func (a *apis) updateTableName(table string) *apis {
	a.ALL = field.NewAsterisk(table)
	a.ID = field.NewUint(table, "id")
	a.CreatedAt = field.NewTime(table, "created_at")
	a.UpdatedAt = field.NewTime(table, "updated_at")
	a.DeletedAt = field.NewUint(table, "deleted_at")
	a.Name = field.NewString(table, "api_name")
	a.Method = field.NewString(table, "method")
	a.Description = field.NewString(table, "description")
	a.SqlType = field.NewInt8(table, "sql_type")
	a.SqlTemplate = field.NewField(table, "sql_template")
	a.SqlTemplateParameters = field.NewField(table, "sql_template_parameters")
	a.SqlTemplateResult = field.NewField(table, "sql_template_result")

	a.fillFieldMap()

	return a
}

func (a *apis) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := a.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (a *apis) fillFieldMap() {
	a.fieldMap = make(map[string]field.Expr, 11)
	a.fieldMap["id"] = a.ID
	a.fieldMap["created_at"] = a.CreatedAt
	a.fieldMap["updated_at"] = a.UpdatedAt
	a.fieldMap["deleted_at"] = a.DeletedAt
	a.fieldMap["api_name"] = a.Name
	a.fieldMap["method"] = a.Method
	a.fieldMap["description"] = a.Description
	a.fieldMap["sql_type"] = a.SqlType
	a.fieldMap["sql_template"] = a.SqlTemplate
	a.fieldMap["sql_template_parameters"] = a.SqlTemplateParameters
	a.fieldMap["sql_template_result"] = a.SqlTemplateResult
}

func (a apis) clone(db *gorm.DB) apis {
	a.apisDo.ReplaceDB(db)
	return a
}

type apisDo struct{ gen.DO }

type IApisDo interface {
	gen.SubQuery
	Debug() IApisDo
	WithContext(ctx context.Context) IApisDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IApisDo
	WriteDB() IApisDo
	As(alias string) gen.Dao
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IApisDo
	Not(conds ...gen.Condition) IApisDo
	Or(conds ...gen.Condition) IApisDo
	Select(conds ...field.Expr) IApisDo
	Where(conds ...gen.Condition) IApisDo
	Order(conds ...field.Expr) IApisDo
	Distinct(cols ...field.Expr) IApisDo
	Omit(cols ...field.Expr) IApisDo
	Join(table schema.Tabler, on ...field.Expr) IApisDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IApisDo
	RightJoin(table schema.Tabler, on ...field.Expr) IApisDo
	Group(cols ...field.Expr) IApisDo
	Having(conds ...gen.Condition) IApisDo
	Limit(limit int) IApisDo
	Offset(offset int) IApisDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IApisDo
	Unscoped() IApisDo
	Create(values ...*model.Apis) error
	CreateInBatches(values []*model.Apis, batchSize int) error
	Save(values ...*model.Apis) error
	First() (*model.Apis, error)
	Take() (*model.Apis, error)
	Last() (*model.Apis, error)
	Find() ([]*model.Apis, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Apis, err error)
	FindInBatches(result *[]*model.Apis, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Apis) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IApisDo
	Assign(attrs ...field.AssignExpr) IApisDo
	Joins(fields ...field.RelationField) IApisDo
	Preload(fields ...field.RelationField) IApisDo
	FirstOrInit() (*model.Apis, error)
	FirstOrCreate() (*model.Apis, error)
	FindByPage(offset int, limit int) (result []*model.Apis, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IApisDo
	UnderlyingDB() *gorm.DB
	schema.Tabler

	GetWithNameAndMethod(name string, method string) (result *model.Apis, err error)
}

//SELECT * FROM @@table WHERE api_name = @name AND method = @method limit 1
func (a apisDo) GetWithNameAndMethod(name string, method string) (result *model.Apis, err error) {
	params := make(map[string]interface{}, 0)

	var generateSQL strings.Builder
	params["name"] = name
	params["method"] = method
	generateSQL.WriteString("SELECT * FROM apis WHERE api_name = @name AND method = @method limit 1 ")

	var executeSQL *gorm.DB
	if len(params) > 0 {
		executeSQL = a.UnderlyingDB().Raw(generateSQL.String(), params).Take(&result)
	} else {
		executeSQL = a.UnderlyingDB().Raw(generateSQL.String()).Take(&result)
	}
	err = executeSQL.Error
	return
}

func (a apisDo) Debug() IApisDo {
	return a.withDO(a.DO.Debug())
}

func (a apisDo) WithContext(ctx context.Context) IApisDo {
	return a.withDO(a.DO.WithContext(ctx))
}

func (a apisDo) ReadDB() IApisDo {
	return a.Clauses(dbresolver.Read)
}

func (a apisDo) WriteDB() IApisDo {
	return a.Clauses(dbresolver.Write)
}

func (a apisDo) Clauses(conds ...clause.Expression) IApisDo {
	return a.withDO(a.DO.Clauses(conds...))
}

func (a apisDo) Returning(value interface{}, columns ...string) IApisDo {
	return a.withDO(a.DO.Returning(value, columns...))
}

func (a apisDo) Not(conds ...gen.Condition) IApisDo {
	return a.withDO(a.DO.Not(conds...))
}

func (a apisDo) Or(conds ...gen.Condition) IApisDo {
	return a.withDO(a.DO.Or(conds...))
}

func (a apisDo) Select(conds ...field.Expr) IApisDo {
	return a.withDO(a.DO.Select(conds...))
}

func (a apisDo) Where(conds ...gen.Condition) IApisDo {
	return a.withDO(a.DO.Where(conds...))
}

func (a apisDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IApisDo {
	return a.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (a apisDo) Order(conds ...field.Expr) IApisDo {
	return a.withDO(a.DO.Order(conds...))
}

func (a apisDo) Distinct(cols ...field.Expr) IApisDo {
	return a.withDO(a.DO.Distinct(cols...))
}

func (a apisDo) Omit(cols ...field.Expr) IApisDo {
	return a.withDO(a.DO.Omit(cols...))
}

func (a apisDo) Join(table schema.Tabler, on ...field.Expr) IApisDo {
	return a.withDO(a.DO.Join(table, on...))
}

func (a apisDo) LeftJoin(table schema.Tabler, on ...field.Expr) IApisDo {
	return a.withDO(a.DO.LeftJoin(table, on...))
}

func (a apisDo) RightJoin(table schema.Tabler, on ...field.Expr) IApisDo {
	return a.withDO(a.DO.RightJoin(table, on...))
}

func (a apisDo) Group(cols ...field.Expr) IApisDo {
	return a.withDO(a.DO.Group(cols...))
}

func (a apisDo) Having(conds ...gen.Condition) IApisDo {
	return a.withDO(a.DO.Having(conds...))
}

func (a apisDo) Limit(limit int) IApisDo {
	return a.withDO(a.DO.Limit(limit))
}

func (a apisDo) Offset(offset int) IApisDo {
	return a.withDO(a.DO.Offset(offset))
}

func (a apisDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IApisDo {
	return a.withDO(a.DO.Scopes(funcs...))
}

func (a apisDo) Unscoped() IApisDo {
	return a.withDO(a.DO.Unscoped())
}

func (a apisDo) Create(values ...*model.Apis) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Create(values)
}

func (a apisDo) CreateInBatches(values []*model.Apis, batchSize int) error {
	return a.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (a apisDo) Save(values ...*model.Apis) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Save(values)
}

func (a apisDo) First() (*model.Apis, error) {
	if result, err := a.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Apis), nil
	}
}

func (a apisDo) Take() (*model.Apis, error) {
	if result, err := a.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Apis), nil
	}
}

func (a apisDo) Last() (*model.Apis, error) {
	if result, err := a.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Apis), nil
	}
}

func (a apisDo) Find() ([]*model.Apis, error) {
	result, err := a.DO.Find()
	return result.([]*model.Apis), err
}

func (a apisDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Apis, err error) {
	buf := make([]*model.Apis, 0, batchSize)
	err = a.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (a apisDo) FindInBatches(result *[]*model.Apis, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return a.DO.FindInBatches(result, batchSize, fc)
}

func (a apisDo) Attrs(attrs ...field.AssignExpr) IApisDo {
	return a.withDO(a.DO.Attrs(attrs...))
}

func (a apisDo) Assign(attrs ...field.AssignExpr) IApisDo {
	return a.withDO(a.DO.Assign(attrs...))
}

func (a apisDo) Joins(fields ...field.RelationField) IApisDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Joins(_f))
	}
	return &a
}

func (a apisDo) Preload(fields ...field.RelationField) IApisDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Preload(_f))
	}
	return &a
}

func (a apisDo) FirstOrInit() (*model.Apis, error) {
	if result, err := a.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Apis), nil
	}
}

func (a apisDo) FirstOrCreate() (*model.Apis, error) {
	if result, err := a.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Apis), nil
	}
}

func (a apisDo) FindByPage(offset int, limit int) (result []*model.Apis, count int64, err error) {
	result, err = a.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = a.Offset(-1).Limit(-1).Count()
	return
}

func (a apisDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = a.Count()
	if err != nil {
		return
	}

	err = a.Offset(offset).Limit(limit).Scan(result)
	return
}

func (a apisDo) Scan(result interface{}) (err error) {
	return a.DO.Scan(result)
}

func (a apisDo) Delete(models ...*model.Apis) (result gen.ResultInfo, err error) {
	return a.DO.Delete(models)
}

func (a *apisDo) withDO(do gen.Dao) *apisDo {
	a.DO = *do.(*gen.DO)
	return a
}
