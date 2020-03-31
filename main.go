package main

import (
	"errors"
	"main/config"
	"main/model"
	"main/router"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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

	// Set gin mode
	gin.SetMode(viper.GetString("runmode"))

	// 测试 log
	//for {
	//	log.Info("111111111111111111111111111111111111111111")
	//	time.Sleep(100 * time.Millisecond)
	//}

	// 创建 gin 引擎
	g := gin.New()

	// 中间件：非具体业务类型的代码
	// gin 中间件 HandleFunc 是用于自定义中间件的方法
	middlewares := []gin.HandlerFunc{}

	// 加载路由
	router.Load(g, middlewares...)

	// init db
	model.DB.Init()
	defer model.DB.Close()

	// Ping the server to make sure the router is working.
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Info("The router has been deployed successfully.")
	}()

	log.Infof("Start to listening the incoming requests on http address: %s", viper.GetString("url"))
	log.Info(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

// pingServer func
func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		resp, err := http.Get(viper.GetString("url") + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}
		// 停 1s 再 ping
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("Cannot connect to the router.")
}
