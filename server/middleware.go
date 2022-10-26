package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type GetRequestMoreData func() map[string]interface{}

func externalDataMiddleware(f GetRequestMoreData) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := f()
		if len(data) > 0 {
			c.Set("external", data)
		}
	}
}

func externalMsgMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("Logger", logger)
	}
}
