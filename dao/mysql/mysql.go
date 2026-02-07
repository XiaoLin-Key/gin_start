package mysql

import (
	"fmt"
	"gin_start/settings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func InitMySQL() (err error) {
	conf := settings.Conf.MySQL

	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)
	db, err = gorm.Open(mysql.Open(dns), &gorm.Config{
		// 开启这个，控制台会打印出它执行的每一行 SQL
		Logger: logger.Default.LogMode(logger.Info),
	})

	// 检查数据库连接是否成功
	if err != nil {
		panic("数据库连接失败")
	}

	// 获取底层的 sql.DB 对象
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(conf.MaxIdleConns) // 最大空闲连接数
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns) // 最大打开连接数
	sqlDB.SetConnMaxLifetime(60)             // 连接最大生命周期（秒）

	return sqlDB.Ping()
}

func GetDB() *gorm.DB {
	return db
}
