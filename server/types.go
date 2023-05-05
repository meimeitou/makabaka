package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResInterface interface {
	CodeTag() string
	MessageTag() string
	DataTag() string
}

type ResponseMsg struct {
	ErrorCode   int         `json:"errno"`
	ErrorMsg    string      `json:"errmsg"`
	Data        interface{} `json:"data"`
	responseTag ResInterface
}

func (r *ResponseMsg) New() *ResponseMsg {
	return &ResponseMsg{
		ErrorCode: 0,
	}
}

func (r *ResponseMsg) CodeTag() string {
	return "code"
}

func (r *ResponseMsg) MessageTag() string {
	return "msg"
}

func (r *ResponseMsg) DataTag() string {
	return "data"
}

func (r *ResponseMsg) WithTag(tag ResInterface) {
	r.responseTag = tag
}

func (r ResponseMsg) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{}
	data[r.responseTag.CodeTag()] = r.ErrorCode
	data[r.responseTag.MessageTag()] = r.ErrorMsg
	data[r.responseTag.DataTag()] = r.Data
	return json.Marshal(data)
}

func (s *Server) responseOk(c *gin.Context) {
	res := ResponseMsg{
		ErrorCode: 0,
		ErrorMsg:  "success",
	}
	res.WithTag(s.responseTag)
	c.JSON(http.StatusOK, res)
}

func (s *Server) responseError(c *gin.Context, code int, err error) {
	res := ResponseMsg{
		ErrorCode: code,
		ErrorMsg:  fmt.Sprintf("%v", err),
	}
	res.WithTag(s.responseTag)
	c.JSON(http.StatusOK, res)
}

func (s *Server) responseOkWithData(c *gin.Context, data interface{}) {
	res := ResponseMsg{
		ErrorCode: 0,
		ErrorMsg:  "success",
		Data:      data,
	}
	res.WithTag(s.responseTag)
	c.JSON(http.StatusOK, res)
}

func (s *Server) responseMsgWithData(c *gin.Context, msg string, data interface{}) {
	res := ResponseMsg{
		ErrorCode: 0,
		ErrorMsg:  msg,
		Data:      data,
	}
	res.WithTag(s.responseTag)
	c.JSON(http.StatusOK, res)
}

type queryPage struct {
	PageSize int `form:"pageSize"`
	PageNum  int `form:"pageNum"`
}

func (q *queryPage) GetPageNum() int {
	return q.PageNum
}

func (q *queryPage) GetPageSize() int {
	return q.PageSize
}

func GetRecordWindows(c *gin.Context, queryStruct pageQuery) (start, size int, err error) {
	if err = c.ShouldBind(queryStruct); err != nil {
		return
	}
	var pageNum = 1
	if queryStruct.GetPageNum() > 0 {
		pageNum = queryStruct.GetPageNum()
	}
	size = 8
	if queryStruct.GetPageSize() > 0 {
		size = queryStruct.GetPageSize()
	}

	start = 0
	if pageNum > 0 {
		start = (pageNum - 1) * size
	}
	return
}

type pageQuery interface {
	GetPageNum() int
	GetPageSize() int
}
