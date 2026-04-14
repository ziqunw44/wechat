package dao

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"wechat/internal/config"
	"wechat/internal/model"
)

var GormDB *gorm.DB

// InitGorm 连接 MySQL 并 AutoMigrate。需在 config.LoadConfig() 之后调用。
func InitGorm() error {
	conf := config.GetConfig()
	user := conf.MysqlConfig.User
	password := conf.MysqlConfig.Password
	host := conf.MysqlConfig.Host
	port := conf.MysqlConfig.Port
	dbName := conf.MysqlConfig.DatabaseName
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)
	var err error
	GormDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return GormDB.AutoMigrate(
		&model.UserInfo{},
		&model.GroupInfo{},
		&model.UserContact{},
		&model.Session{},
		&model.ContactApply{},
		&model.Message{},
	)
}

// MustInitGorm 用于进程入口：失败则退出。
func MustInitGorm() {
	if err := InitGorm(); err != nil {
		log.Fatal(err)
	}
}
