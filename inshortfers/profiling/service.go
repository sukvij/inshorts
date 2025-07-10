package profiling

import (
	"net/http/pprof" // Import pprof

	"github.com/gin-gonic/gin"
)

func Profiling(app *gin.Engine) {
	// Create a new router group for pprof endpoints
	pprofGroup := app.Group("/debug/pprof")
	{
		pprofGroup.GET("/", gin.WrapF(pprof.Index))
		pprofGroup.GET("/profile", gin.WrapF(pprof.Profile))
		pprofGroup.GET("/heap", gin.WrapF(pprof.Handler("heap").ServeHTTP))
		pprofGroup.GET("/goroutine", gin.WrapF(pprof.Handler("goroutine").ServeHTTP))
		pprofGroup.GET("/mutex", gin.WrapF(pprof.Handler("mutex").ServeHTTP))
	}
}
