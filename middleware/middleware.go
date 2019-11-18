package middleware

import "github.com/gin-gonic/gin"

func CORS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		if ctx.Request.Method == "OPTIONS" {
			ctx.Header("Content-Type", "application/json; charset=utf-8")
			ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			ctx.Header("Access-Control-Allow-Headers", "X-ACCESS-TOKEN,X-RUNNING-ENV,Content-Type")
			ctx.Header("Access-Control-Max-Age", "3600")
			ctx.Status(204)
			ctx.Writer.Write([]byte(""))
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
