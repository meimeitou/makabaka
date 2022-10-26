package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseMsg struct {
	ErrorCode int         `json:"errno"`
	ErrorMsg  string      `json:"errmsg"`
	Data      interface{} `json:"data"`
}

func responseOk(c *gin.Context) {
	c.JSON(http.StatusOK, ResponseMsg{
		ErrorCode: 0,
		ErrorMsg:  "success",
	})
}

func responseError(c *gin.Context, code int, err error) {
	c.JSON(http.StatusOK, ResponseMsg{
		ErrorCode: code,
		ErrorMsg:  fmt.Sprintf("%v", err),
	})
}

func responseOkWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ResponseMsg{
		ErrorCode: 0,
		ErrorMsg:  "success",
		Data:      data,
	})
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
