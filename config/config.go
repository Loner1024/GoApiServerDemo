package config

import (
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

type Config struct {
	Name string
}

func Init(cfg string) error {
	c := Config{Name: cfg}

	// 初始化配置文件
	if err := c.initConfig(); err != nil {
		return err
	}
	c.watchConfig()
	return nil
}

func (c *Config) initConfig() error {
	if c.Name != "" {
		viper.SetConfigFile(c.Name) // 指定了配置文件，就去使用
	} else {
		viper.AddConfigPath("conf") // 解析默认配置文件
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")     // 设置配置文件格式
	viper.AutomaticEnv()            // 读取匹配的环境变量
	viper.SetEnvPrefix("APISERVER") // 设置环境变量前缀
	/*
		实例化的 NewReplacer
		有两个方法：
		Replace() 返回s的所有替换进行完后的拷贝。
		WriteString向w中写入s的所有替换进行完后的拷贝。
		可以提供多个替换，按传入顺序进行，old,new,old,new...
	*/
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

// 监听文件变化的 func
func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置发生变化，执行响应操作
		log.Printf("Config file changed: %s", e.Name)
	})
}
