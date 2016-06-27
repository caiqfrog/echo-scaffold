package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/echo/engine/standard"
	//"log"

	"github.com/caiqfrog/echo-scaffold/utils"
	"github.com/caiqfrog/echo-scaffold/controllers"
	//"github.com/caiqfrog/echo-scaffold/mongo"
	"github.com/caiqfrog/echo-scaffold/test"
)

func main() {
	// 连接mango服务器
	//if e := mongo.Dial("mongodb://127.0.0.1:27017", "colledge", 10); nil != e {
	//	log.Fatal("dial: ", e)
	//	return
	//}
	//defer mongo.Close()
	// 摘要认证
	//digest := mw.NewDigest("example.com", func(user, realm string) string {
	//	// 获取密码
	//	// 测试，密码为hello
	//	return "hello"
	//})
	//digest.Allow("*", "/")
	//digest.Allow("*", "/login")
	//digest.Allow("*", "/register")
	//digest.Allow("*", "/static/*")

	go func() {
		// 创建并启动http服务
		e := echo.New()
		// 摘要认证中间件
		//e.Use(digest.Process)
		// 跨域访问中间件
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"http://127.0.0.1:3002", "http://localhost:3002"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAuthorization},
			AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		}))
		// router
		controllers.Register(e, test.TestController{Prefix:""})

		e.Run(standard.New(":8080"))
	}()

	// TODO 资源管理
	// 中断信号捕捉
	utils.Signal()
}
