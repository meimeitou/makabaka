package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/meimeitou/makabaka/pkg/bind"
	"github.com/meimeitou/makabaka/pkg/tpl"
	"github.com/oklog/run"
	"github.com/sirupsen/logrus"
)

type Server struct {
	router           *gin.Engine
	logger           *logrus.Logger
	addr             string
	prefix           string
	querymiddleware  []gin.HandlerFunc
	adminMiddleware  gin.HandlerFunc
	responseTag      ResInterface
	checkApinRequest CheckApiRequest
}

func NewServer(logger *logrus.Logger, addr, prefix string) *Server {
	r := gin.New()
	r.Use(gLogger(logger))
	r.Use(externalMsgMiddleware(logger))
	r.Use(gin.Recovery())

	return &Server{
		router:           r,
		logger:           logger,
		addr:             addr,
		prefix:           prefix,
		responseTag:      &ResponseMsg{},
		checkApinRequest: defaultCheckApiRequest,
	}
}

func (r *Server) WithTag(tag ResInterface) {
	r.responseTag = tag
}

// add custom validation
func (r *Server) RegisterBindValidator(tag string, fn validator.Func, callValidationEvenIfNull ...bool) {
	bind.RegisterValidation(tag, fn, callValidationEvenIfNull...)
}

// add template functions
func (r *Server) RegisterTemplateFunction(funcMap template.FuncMap) {
	tpl.RegisterFunction(funcMap)
}

// middleware for query api
func (s *Server) QueryMiddleware(m ...gin.HandlerFunc) {
	s.querymiddleware = append(s.querymiddleware, m...)
}

// check api type auth
func (s *Server) WithCheckApinRequest(f CheckApiRequest) {
	s.checkApinRequest = f
}

func (s *Server) AdminMiddleware(m gin.HandlerFunc) {
	s.adminMiddleware = m
}

// get gin router group
func (s *Server) GetRouterGroup(prefix string) *gin.RouterGroup {
	return s.router.Group(prefix)
}

func (s *Server) Run(g *run.Group) {
	s.routerRegist(s.router, s.prefix)
	srv := &http.Server{
		Addr:         s.addr,
		Handler:      s.router,
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

func (s *Server) routerRegist(r *gin.Engine, prefix string) {
	r.NoRoute(func(c *gin.Context) {
		s.responseError(c, 404, fmt.Errorf("Page not found"))
	})
	r.GET("/healthz", s.healthz)
	r.GET("/ready", s.ready)
	root := r.Group(prefix)
	{
		// admin
		var admin *gin.RouterGroup
		if s.adminMiddleware != nil {
			admin = root.Group("/admin", s.adminMiddleware)
		} else {
			admin = root.Group("/admin")
		}
		admin.GET("/proxy/list", s.ProxyList)
		admin.POST("/apis/create", s.ApiCreate)
		admin.GET("/apis/:db/list", s.ApiList)
		admin.GET("/apis/:db/detail/:id", s.ApiGet)
		admin.DELETE("/apis/:db/:api", s.ApiDelete)
	}
	{
		// query api
		var query *gin.RouterGroup
		if len(s.querymiddleware) > 0 {
			query = root.Group("/query", s.querymiddleware...)
		} else {
			query = root.Group("/query")
		}
		query.Any("/:db/:name", s.Query)
	}
}
