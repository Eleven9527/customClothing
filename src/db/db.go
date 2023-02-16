package db

import (
	"customClothing/src/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	db *gorm.DB
)

func Db() *gorm.DB {
	return db
}

func InitDb() {
	conn, err := gorm.Open("mysql", config.Cfg().DbCfg.Dsn)
	if err != nil {
		fmt.Println("初始化数据库失败:", err.Error())
		panic(err)
	}

	conn.DB().SetMaxIdleConns(config.Cfg().DbCfg.MaxIdle)
	conn.DB().SetMaxOpenConns(config.Cfg().DbCfg.MaxOpen)

	db = conn
}
