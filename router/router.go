package router

import (
	"main/handler/sd"
	"main/handler/user"
	"main/router/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// 设置 HTTP Header
	g.Use(gin.Recovery())     // 遇到异常或 bug 时，通过 gin.Recovery() 恢复 API 服务器
	g.Use(middleware.NoCache) // 强制浏览器不使用缓存
	// 在正式跨域之前，浏览器会根据需要发起一次预检（也就是option请求），用来让服务端返回允许的方法（如get、post），被跨域访问的Origin（来源或者域），还有是否需要Credentials(认证信息)等。
	g.Use(middleware.Options) // 浏览器跨域 OPTIONS 请求设置
	g.Use(middleware.Secure)  // 安全设置

	// 404 NoRoute 默认返回 404
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	u := g.Group("/v1/user")
	{
		u.POST("", user.Create)       // 创建用户
		u.DELETE("/:id", user.Delete) // 删除用户
		u.PUT("/:id", user.Update)    // 更新用户
		u.GET("", user.List)          // 用户列表
		u.GET("/:username", user.Get) // 获取指定用户的详细信息

	}

	// Group 创建路由组，具有相同路径的前缀的路由，通用中间件
	// 健康检查处理程序，检查 API Server 的状态：健康状况、服务器硬盘、CPU 和内存使用量
	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}
	return g
}
