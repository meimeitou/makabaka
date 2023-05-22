package server

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/meimeitou/makabaka/config"
	"github.com/meimeitou/makabaka/db"
	"github.com/meimeitou/makabaka/model/exec"
	"github.com/meimeitou/makabaka/pkg/bind"
	"golang.org/x/exp/maps"
	"gorm.io/gorm"
)

const (
	keepParametersTest = "test"
)

type (
	QueryParams map[string]interface{}
)

func (q QueryParams) GetTest() (bool, error) {
	data, ok := q[keepParametersTest]
	if !ok {
		return false, nil
	}
	switch v := data.(type) {
	case string:
		return v == "true" || v == "True", nil
	case bool:
		return v, nil
	default:
		return false, errors.New("system keep parameter (test) type error, must bool")
	}
}

// query for type template
func (s *Server) Query(c *gin.Context) {
	body, err := getQueryParameters(c)
	if err != nil {
		s.responseError(c, 400, err)
		return
	}
	// valid and get test
	test, err := body.GetTest()
	if err != nil {
		s.responseError(c, 400, err)
		return
	}

	db := c.Param("db")
	conn, err := config.DBSet.GetDB(db)
	if err != nil {
		s.responseError(c, 400, err)
		return
	}
	apiName := c.Param("name")
	query := conn.GetQuery().Apis
	res, err := query.GetWithNameAndMethod(apiName, c.Request.Method)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.responseError(c, 404, errors.New("api not exist"))
			return
		}
		s.responseError(c, 500, err)
		return
	}
	if s.checkAdminRequest(c, string(res.ApiType)) {
		s.responseError(c, 403, errors.New("用户权限错误"))
		return
	}
	// externalData can set by upstream server
	if externalData, exits := c.Get(externalData); exits {
		if src, ok := externalData.(map[string]interface{}); ok {
			maps.Copy(body, src)
		}
	}
	// parameter valid
	if err := bind.ValidateInput(body, res.SqlTemplateParameters); err != nil {
		s.responseError(c, 400, err)
		return
	}
	sqls, err := res.SqlTemplate.ToMapString()
	if err != nil {
		s.responseError(c, 500, err)
		return
	}
	// sql exec
	if test && config.Cfg.EnableTest {
		s.queryTest(c, conn, sqls, body)
	} else {
		s.queryMulti(c, conn, sqls, body)
	}
}

func getQueryParameters(c *gin.Context) (QueryParams, error) {
	body := make(QueryParams)
	// bind from query and body
	/*
		method = GET, default from form
		method = DELETE/POST/PUT/... will read header Content-Type
		Content-Type define:

			MIMEJSON              = "application/json"
			MIMEHTML              = "text/html"
			MIMEXML               = "application/xml"
			MIMEXML2              = "text/xml"
			MIMEPlain             = "text/plain"
			MIMEPOSTForm          = "application/x-www-form-urlencoded"
			MIMEMultipartPOSTForm = "multipart/form-data"
			MIMEPROTOBUF          = "application/x-protobuf"
			MIMEMSGPACK           = "application/x-msgpack"
			MIMEMSGPACK2          = "application/msgpack"
			MIMEYAML              = "application/x-yaml"
			MIMETOML              = "application/toml"
	*/
	b := binding.Default(c.Request.Method, c.ContentType())
	switch b {
	case binding.Form, binding.FormMultipart, binding.FormPost:
		data := make(map[string]string)
		if err := c.ShouldBindWith(&data, b); err != nil {
			return nil, err
		}
		for k, v := range data {
			body[k] = v
		}
	default:
		if err := c.ShouldBindWith(&body, b); err != nil {
			return nil, err
		}
	}
	Logger(c).Debug(body)
	return body, nil
}

func (s *Server) queryMulti(c *gin.Context, conn *db.Conn, sqls map[string]string, values map[string]interface{}) {
	res := make(map[string]interface{})
	keys := []string{}
	rawsql := []string{}
	for key, sql := range sqls {
		builder := exec.NewQueryBuilder(conn, sql, values)
		data, err := builder.Exec()
		if err != nil {
			s.responseError(c, 400, err)
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
			s.responseMsgWithData(c, strings.Join(rawsql, ";"), res[keys[0]])
			return
		}
		s.responseOkWithData(c, res[keys[0]])
	} else if len(res) > 1 { // more
		if config.Cfg.EnableTest {
			s.responseMsgWithData(c, strings.Join(rawsql, ";"), res)
			return
		}
		s.responseOkWithData(c, res)
	} else { // zero
		s.responseOk(c)
	}

}

func (s *Server) queryTest(c *gin.Context, conn *db.Conn, sqls map[string]string, values map[string]interface{}) {
	res := make(map[string]string)
	for key, sql := range sqls {
		raw, err := exec.NewQueryBuilder(conn, sql, values).TemplateParse()
		if err != nil {
			s.responseError(c, 400, err)
			return
		}
		res[key] = raw
	}
	s.responseOkWithData(c, res)
}
