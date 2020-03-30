package main

import (
	"errors"
	"log"
	"main/router"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
)

var (
	cfg = pflag.StringP("config", "c", "", "apiserver config file path")
)

func main() {
	pflag.Parse()
	// 初始化 config
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	// 创建 gin 引擎
	g := gin.New()

	// 中间件：非具体业务类型的代码
	// gin 中间件 HandleFunc 是用于自定义中间件的方法
	middlewares := []gin.HandlerFunc{}

	// 加载路由
	router.Load(g, middlewares...)

	// Ping the server to make sure the router is working.
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Print("The router has been deployed successfully.")
	}()

	log.Printf("Start to listening the incoming requests on http address: %s", ":8080")
	log.Printf(http.ListenAndServe(":8080", g).Error())
}

// pingServer func
func pingServer() error {
	for i := 0; i < 10; i++ {
		resp, err := http.Get("http://127.0.0.1/" + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}
		// 停 1s 再 ping
		log.Print("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("Cannot connect to the router.")
}