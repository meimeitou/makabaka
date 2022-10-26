package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/run"
	"github.com/sirupsen/logrus"
)

type Server struct {
	logger     *logrus.Logger
	addr       string
	prefix     string
	middleware gin.HandlerFunc
}

func NewServer(logger *logrus.Logger, addr, prefix string) *Server {
	return &Server{
		logger: logger,
		addr:   addr,
		prefix: prefix,
	}
}

func (s *Server) WithMiddlewarem(m gin.HandlerFunc) {
	s.middleware = m
}

// 消息api
func (s *Server) Run(g *run.Group) {
	r := gin.New()
	r.Use(gLogger(s.logger))
	r.Use(externalMsgMiddleware(s.logger))
	r.Use(gin.Recovery())
	s.RouterRegist(r, s.prefix)

	srv := &http.Server{
		Addr:    s.addr,
		Handler: r,
	}
	g.Add(
		func() error {
			if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
				s.logger.Printf("listen: %s\n", err)
			}
			return nil
		},
		func(e error) {
			s.logger.Info("stopping server...")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				s.logger.Fatal("Server forced to shutdown:", err)
			}
		},
	)
}

func (s *Server) RouterRegist(r *gin.Engine, prefix string) {
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code": "PAGE_NOT_FOUND",
			"msg":  "Page not found"})
	})
	root := r.Group(prefix)
	root.GET("/admin/proxy/list", ProxyList)
	root.POST("/admin/apis/create", ApiCreate)
	root.GET("/admin/apis/:db/list", ApiList)
	// query api
	var query *gin.RouterGroup
	if s.middleware != nil {
		query = root.Group("/query", s.middleware)
	} else {
		query = root.Group("/query")
	}
	query.Any("/:db/:name", Query)
}
