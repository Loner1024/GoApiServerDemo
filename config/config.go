package config

import (
	"strings"

	"github.com/fsnotify/fsnotify"

	"github.com/lexkong/log"
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
	// 初始化日志包
	c.initLog()
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

func (c *Config) initLog() {
	passLagerCfg := log.PassLagerCfg{
		Writers:        viper.GetString("log.writers"),       // 输出位置，选择 file 会将日志记录到 logger_file 指定的日志文件中，选择 stdout 会将日志输出到标准输出，当然也可以两者同时选择
		LoggerLevel:    viper.GetString("log.logger_level"),  // 日志级别，DEBUG、INFO、WARN、ERROR、FATAL
		LoggerFile:     viper.GetString("log.logger_file"),   // 日志文件
		LogFormatText:  viper.GetBool("log.log_format_text"), // 日志输出格式，JSON 或者 plaintext，true 会输出成非 JSON 格式，
		RollingPolicy:  viper.GetString("log.rollingPolicy"), // rotate 依据，如果选 daily 则根据天进行转存，如果是 size 则根据大小进行转存
		LogRotateDate:  viper.GetInt("log.log_rotate_date"),  // rotate 转存时间
		LogRotateSize:  viper.GetInt("log.log_rotate_size"),  // rotate 转存大小
		LogBackupCount: viper.GetInt("log.log_backup_count"), // 当日志文件达到转存标准时，log 系统会将该日志文件进行压缩备份，这里指定了备份文件的最大个数
	}
	log.InitWithConfig(&passLagerCfg)
}

// 监听文件变化的 func
func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置发生变化，执行响应操作
		log.Infof("Config file changed: %s", e.Name)
	})
}
