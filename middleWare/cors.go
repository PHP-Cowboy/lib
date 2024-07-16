package middleWare

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func Cors() gin.HandlerFunc {

	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"token", "Origin", "Authorization", "Content-Type", "Access-Token", "Package-Id", "Version"},
		ExposeHeaders:    []string{"Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 6 * time.Hour,
	}
	return cors.New(corsConfig)
}

// 屏蔽所有GET请求,尽量测试用的GET
func BlockGetRequests() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查请求方法是否为GET
		if c.Request.Method == "GET" && gin.Mode() == gin.ReleaseMode {
			c.JSON(403, gin.H{"error": "Requests are not allowed"})
			c.Abort() // 终止请求处理
		}
		// 如果不是GET请求，则继续执行后续的路由处理
		c.Next()
	}
}

// 获取包名
func GetPackageID(request *gin.Context) string {
	return request.Request.Header.Get("Package-ID")
}

func GetVersion(request *gin.Context) string {
	return request.Request.Header.Get("version")
}
