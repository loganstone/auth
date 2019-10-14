package handler

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/middleware"
)

func newRouter() *gin.Engine {
	router := gin.New()
	router.Use(middleware.LogFormat())
	router.Use(middleware.RequestID())
	router.Use(gin.Recovery())

	bind(router)

	if gin.Mode() == gin.DebugMode {
		// Debug uri - /debug/pprof/
		pprof.Register(router)
	}

	return router
}

// New .
func New() *gin.Engine {
	return newRouter()
}

// NewTest .
func NewTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return newRouter()
}
