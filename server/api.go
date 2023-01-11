package server

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/meimeitou/makabaka/config"
	"github.com/meimeitou/makabaka/model"
)

type (
	ApiCreateData struct {
		DB                    string                 `json:"db"`
		ID                    uint                   `json:"id"`
		Name                  string                 `json:"name"`
		Method                string                 `json:"method"`
		Description           string                 `json:"description"`
		SqlType               string                 `json:"sqlType"` // template chain
		SqlTemplate           map[string]string      `json:"sqlTemplate"`
		SqlTemplateParameters map[string]interface{} `json:"sqlTemplateParameters"`
		SqlTemplateResult     map[string]interface{} `json:"sqlTemplateResult"`
	}

	apiListParams struct {
		queryPage
		Name string `form:"name"`
	}
)

func (a *ApiCreateData) valid() error {
	if a.Name == "" || a.DB == "" || a.Method == "" {
		return errors.New("name/db/Method can not empty")
	}
	if a.SqlTemplate == nil {
		return errors.New("template can not empty")
	}
	return nil
}

func ProxyList(c *gin.Context) {
	responseOkWithData(c, config.Cfg.Proxy)
}

func ApiCreate(c *gin.Context) {
	data := ApiCreateData{}
	if err := c.ShouldBindJSON(&data); err != nil {
		responseError(c, 400, err)
		return
	}
	Logger(c).Info(data)
	if err := data.valid(); err != nil {
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
	sql := apis.Select(apis.ID, apis.Name, apis.Method, apis.SqlType, apis.CreatedAt)
	if payload.Name != "" {
		data, count, err = sql.Where(apis.Name.Like(fmt.Sprintf("%%%s%%", payload.Name))).FindByPage(start, size)
	} else {
		data, count, err = sql.FindByPage(start, size)
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

func ApiGet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		responseError(c, 400, fmt.Errorf("参数错误：%v", err))
		return
	}
	db := c.Param("db")
	conn, err := config.DBSet.GetDB(db)
	if err != nil {
		responseError(c, 400, err)
		return
	}
	apis := conn.GetQuery().Apis
	var data *model.Apis
	if data, err = apis.Where(apis.ID.Eq(uint(id))).First(); err != nil {
		responseError(c, 500, err)
		return
	}
	responseOkWithData(c, data)
}
