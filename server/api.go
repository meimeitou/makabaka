package server

import (
	"errors"
	"fmt"
	"net/http"
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
		ApiType               string                 `json:"apiType"`
		Description           string                 `json:"description"`
		SqlType               string                 `json:"sqlType"` // template chain
		SqlTemplate           map[string]string      `json:"sqlTemplate"`
		SqlTemplateParameters map[string]interface{} `json:"sqlTemplateParameters"`
		SqlTemplateResult     map[string]interface{} `json:"sqlTemplateResult"`
	}

	apiListParams struct {
		queryPage
		Name    string `form:"name"`
		ApiType string `form:"apiType"`
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

func (s *Server) ProxyList(c *gin.Context) {
	s.responseOkWithData(c, config.Cfg.Proxy)
}

func (s *Server) ApiCreate(c *gin.Context) {
	data := ApiCreateData{}
	if err := c.ShouldBindJSON(&data); err != nil {
		s.responseError(c, 400, err)
		return
	}
	Logger(c).Info(data)
	if err := data.valid(); err != nil {
		s.responseError(c, 400, err)
		return
	}
	conn, err := config.DBSet.GetDB(data.DB)
	if err != nil {
		s.responseError(c, 400, err)
		return
	}
	sqlType, err := model.SqlTypeFromStr(data.SqlType)
	if err != nil {
		s.responseError(c, 400, err)
		return
	}
	a := model.Apis{
		Modelx: model.Modelx{
			ID: data.ID,
		},
		Name:                  data.Name,
		ApiType:               model.ApiType(data.ApiType),
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
			s.responseError(c, 500, err)
			return
		}
	} else {
		if err := api.Create(&a); err != nil {
			s.responseError(c, 500, err)
			return
		}
	}
	s.responseOkWithData(c, a)
}

func (s *Server) ApiDelete(c *gin.Context) {
	db := c.Param("db")
	conn, err := config.DBSet.GetDB(db)
	if err != nil {
		s.responseError(c, 400, err)
		return
	}
	apiName := c.Param("api")
	apis := conn.GetQuery().Apis
	info, err := apis.WithContext(c.Request.Context()).Where(apis.Name.Eq(apiName)).Delete()
	if err != nil {
		s.responseError(c, 500, err)
		return
	}
	s.responseOkWithData(c, info)
}

func (s *Server) ApiList(c *gin.Context) {
	db := c.Param("db")
	conn, err := config.DBSet.GetDB(db)
	if err != nil {
		s.responseError(c, 400, err)
		return
	}
	payload := apiListParams{}
	apis := conn.GetQuery().Apis
	start, size, err := GetRecordWindows(c, &payload)
	if err != nil {
		s.responseError(c, 400, err)
		return
	}

	sql := apis.Select(apis.ID, apis.Name, apis.Method, apis.ApiType, apis.Description, apis.CreatedAt)
	if payload.Name != "" {
		sql = sql.Where(apis.Name.Like(fmt.Sprintf("%%%s%%", payload.Name)))
	}
	if payload.ApiType != "" {
		sql = sql.Where(apis.ApiType.Eq(payload.ApiType))
	}
	var (
		data  []*model.Apis
		count int64
	)
	data, count, err = sql.FindByPage(start, size)
	if err != nil {
		Logger(c).Error(err)
		s.responseError(c, 500, err)
		return
	}
	s.responseOkWithData(c, map[string]interface{}{
		"count": count,
		"item":  data,
	})
}

func (s *Server) ApiGet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.responseError(c, 400, fmt.Errorf("参数错误：%v", err))
		return
	}
	db := c.Param("db")
	conn, err := config.DBSet.GetDB(db)
	if err != nil {
		s.responseError(c, 400, err)
		return
	}
	apis := conn.GetQuery().Apis
	var data *model.Apis
	if data, err = apis.Where(apis.ID.Eq(uint(id))).First(); err != nil {
		s.responseError(c, 500, err)
		return
	}
	s.responseOkWithData(c, data)
}

func (s *Server) healthz(c *gin.Context) {
	s.responseOk(c)
}

func (s *Server) ready(c *gin.Context) {
	if config.DBInited {
		s.responseOk(c)
		return
	}
	c.AbortWithStatus(http.StatusMisdirectedRequest)
}
