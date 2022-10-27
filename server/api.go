package server

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/meimeitou/makabaka/config"
	"github.com/meimeitou/makabaka/db"
	"github.com/meimeitou/makabaka/model"
	"github.com/meimeitou/makabaka/model/exec"
	"github.com/meimeitou/makabaka/pkg/bind"
	"golang.org/x/exp/maps"
	"gorm.io/gorm"
)

type (
	ApiCreateData struct {
		DB                    string                 `json:"db"`
		ID                    uint                   `json:"id"`
		Name                  string                 `json:"name"`
		Method                string                 `json:"method"`
		Description           string                 `json:"desc"`
		SqlType               string                 `json:"sqlType"`
		SqlTemplate           map[string]string      `json:"sqlTemplate"`
		SqlTemplateParameters map[string]interface{} `json:"sqlTemplateParameters"`
		SqlTemplateResult     map[string]interface{} `json:"sqlTemplateResult"`
	}
	QueryParams struct {
		Test               bool                   `form:"test"`
		TemplateParameters map[string]interface{} `form:"templateParameters"`
		ChainParameters    map[string]interface{} `form:"chainParameters"`
	}
	apiListParams struct {
		queryPage
		Name string `form:"name"`
	}
)

func ProxyList(c *gin.Context) {
	responseOkWithData(c, config.Cfg.Proxy)
}

func ApiCreate(c *gin.Context) {
	data := ApiCreateData{}
	if err := c.ShouldBindJSON(&data); err != nil {
		responseError(c, 400, err)
		return
	}
	conn, err := config.DBSet.GetDB(data.DB)
	if err != nil {
		responseError(c, 400, err)
		return
	}
	sqlType, err := model.SqlTypeFromStr(data.SqlType)
	if err != nil {
		responseError(c, 400, err)
		return
	}
	a := model.Apis{
		Modelx: model.Modelx{
			ID: data.ID,
		},
		Name:                  data.Name,
		Method:                data.Method,
		Description:           data.Description,
		SqlType:               sqlType,
		SqlTemplate:           model.FromMapString(data.SqlTemplate),
		SqlTemplateParameters: data.SqlTemplateParameters,
		SqlTemplateResult:     data.SqlTemplateResult,
	}
	api := conn.GetQuery().Apis
	if a.ID > 0 {
		if err := api.Save(&a); err != nil {
			responseError(c, 500, err)
			return
		}
	} else {
		if err := api.Create(&a); err != nil {
			responseError(c, 500, err)
			return
		}
	}
	responseOkWithData(c, a)
}

func ApiDelete(c *gin.Context) {
	db := c.Param("db")
	conn, err := config.DBSet.GetDB(db)
	if err != nil {
		responseError(c, 400, err)
		return
	}
	apiName := c.Param("api")
	apis := conn.GetQuery().Apis
	info, err := apis.WithContext(c.Request.Context()).Where(apis.Name.Eq(apiName)).Delete()
	if err != nil {
		responseError(c, 500, err)
		return
	}
	responseOkWithData(c, info)
}

func ApiList(c *gin.Context) {
	db := c.Param("db")
	conn, err := config.DBSet.GetDB(db)
	if err != nil {
		responseError(c, 400, err)
		return
	}
	payload := apiListParams{}
	apis := conn.GetQuery().Apis
	start, size, err := GetRecordWindows(c, &payload)
	if err != nil {
		responseError(c, 400, err)
		return
	}
	var (
		data  []*model.Apis
		count int64
	)
	if payload.Name != "" {
		data, count, err = apis.Where(apis.Name.Like(fmt.Sprintf("%%%s%%", payload.Name))).FindByPage(start, size)
	} else {
		data, count, err = apis.FindByPage(start, size)
	}
	if err != nil {
		Logger(c).Error(err)
		responseError(c, 500, err)
		return
	}
	responseOkWithData(c, map[string]interface{}{
		"count": count,
		"item":  data,
	})
}

func Query(c *gin.Context) {
	method := c.Request.Method
	body := QueryParams{}
	// bind from query and body
	if err := c.ShouldBind(&body); err != nil {
		responseError(c, 400, fmt.Errorf("请求body错误, %v", err))
		return
	}
	Logger(c).Info(body)
	db := c.Param("db")
	conn, err := config.DBSet.GetDB(db)
	if err != nil {
		responseError(c, 400, err)
		return
	}
	apiName := c.Param("name")
	query := conn.GetQuery().Apis
	res, err := query.GetWithNameAndMethod(apiName, method)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			responseError(c, 404, errors.New("接口不存在"))
			return
		}
		responseError(c, 500, err)
		return
	}
	if res.SqlType == model.SqlTypeChain {
		responseError(c, 400, errors.New("chain not support now."))
		return
	}
	if body.TemplateParameters == nil {
		body.TemplateParameters = make(map[string]interface{})
	}
	if externalData, exits := c.Get("external"); exits {
		if src, ok := externalData.(map[string]interface{}); ok {
			maps.Copy(body.TemplateParameters, src)
		}
	}
	// parameter valid
	if err := bind.ValidateInput(body.TemplateParameters, res.SqlTemplateParameters); err != nil {
		responseError(c, 400, err)
		return
	}
	sqls, err := res.SqlTemplate.ToMapString()
	if err != nil {
		responseError(c, 500, err)
		return
	}
	// sql exec
	if body.Test && config.Cfg.EnableTest {
		queryTest(c, conn, sqls, body.TemplateParameters)
	} else {
		queryMulti(c, conn, sqls, body.TemplateParameters)
	}
}

func queryMulti(c *gin.Context, conn *db.Conn, sqls map[string]string, values map[string]interface{}) {
	res := make(map[string]interface{})
	keys := []string{}
	rawsql := []string{}
	for key, sql := range sqls {
		builder := exec.NewQueryBuilder(conn, sql, values)
		data, err := builder.Exec()
		if err != nil {
			responseError(c, 400, err)
			return
		}
		rawsql = append(rawsql, builder.GetRawSql())
		keys = append(keys, key)
		if len(data) == 1 {
			res[key] = data[0]
		} else {
			res[key] = data
		}
	}
	if len(res) == 1 { // one
		if config.Cfg.EnableTest {
			responseMsgWithData(c, strings.Join(rawsql, ";"), res[keys[0]])
			return
		}
		responseOkWithData(c, res[keys[0]])
	} else if len(res) > 1 { // more
		if config.Cfg.EnableTest {
			responseMsgWithData(c, strings.Join(rawsql, ";"), res)
			return
		}
		responseOkWithData(c, res)
	} else { // zero
		responseOk(c)
	}

}

func queryTest(c *gin.Context, conn *db.Conn, sqls map[string]string, values map[string]interface{}) {
	res := make(map[string]string)
	for key, sql := range sqls {
		raw, err := exec.NewQueryBuilder(conn, sql, values).TemplateParse()
		if err != nil {
			responseError(c, 400, err)
			return
		}
		res[key] = raw
	}
	responseOkWithData(c, res)
}
