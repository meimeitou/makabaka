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
	logger          *logrus.Logger
	addr            string
	prefix          string
	querymiddleware []gin.HandlerFunc
	adminMiddleware gin.HandlerFunc
}

func NewServer(logger *logrus.Logger, addr, prefix string) *Server {
	return &Server{
		logger: logger,
		addr:   addr,
		prefix: prefix,
	}
}

func (s *Server) QueryMiddleware(m ...gin.HandlerFunc) {
	s.querymiddleware = append(s.querymiddleware, m...)
}

func (s *Server) AdminMiddleware(m gin.HandlerFunc) {
	s.adminMiddleware = m
}

// 消息api
func (s *Server) Run(g *run.Group) {
	r := gin.New()
	r.Use(gLogger(s.logger))
	r.Use(externalMsgMiddleware(s.logger))
	r.Use(gin.Recovery())
	s.RouterRegist(r, s.prefix)

	srv := &http.Server{
		Addr:         s.addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	g.Add(
		func() error {
			if err := srv.ListenAndServe(); err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					s.logger.Printf("listen: %s\n", err)
					return nil
				}
				s.logger.Error(err)
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
	{
		var admin *gin.RouterGroup
		if s.adminMiddleware != nil {
			admin = root.Group("/admin", s.adminMiddleware)
		} else {
			admin = root.Group("/admin")
		}
		admin.GET("/proxy/list", ProxyList)
		admin.POST("/apis/create", ApiCreate)
		admin.GET("/apis/:db/list", ApiList)
		admin.GET("/apis/:db/detail/:id", ApiGet)
		admin.DELETE("/apis/:db/:api", ApiDelete)
	}
	{
		// query api
		var query *gin.RouterGroup
		if len(s.querymiddleware) > 0 {
			query = root.Group("/query", s.querymiddleware...)
		} else {
			query = root.Group("/query")
		}
		query.Any("/:db/:name", Query)
	}
}
