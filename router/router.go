package router

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/handler"
)

// New .
func New() *gin.Engine {
	router := gin.Default()

	handler.Bind(router)

	if gin.Mode() == gin.DebugMode {
		// Debug uri - /debug/pprof/
		pprof.Register(router)
	}

	return router
}
