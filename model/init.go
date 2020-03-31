package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
)

// 用于管理多个数据库
type Database struct {
	Self   *gorm.DB
	Docker *gorm.DB
}

var DB *Database

// 调用 GetSelfDB() 和 GetDockerDB() 方法来同时创建两个 Database 的数据库对象。
func (db *Database) Init() {
	DB = &Database{
		Self:   GetSelfDB(),
		Docker: GetDockerDB(),
	}
}

func openDB(name string) *gorm.DB {
	db, err := gorm.Open("sqlite3", "/tmp/gorm.db")
	if err != nil {
		log.Errorf(err, "Database connection failed. Database name: %s", name)
	}
	log.Infof("Database connection successful. Database name:%s", name)
	defer db.Close()
	return db
}

func GetSelfDB() *gorm.DB {
	return openDB(viper.GetString("db.name"))
}

func GetDockerDB() *gorm.DB {
	return openDB(viper.GetString("docker_db.username"))
}

func (db *Database) Close() {
	DB.Self.Close()
	DB.Docker.Close()
}
