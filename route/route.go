package route

import (
	"PrometheusProject/middleware"
	"PrometheusProject/v1/admin"
	"PrometheusProject/v1/view"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	var app = gin.New()
	app.Use(gin.Recovery(), middleware.CORS())
	admin.Mapping("/admin", app)
	view.Mapping("/view", app)
	return app
}

