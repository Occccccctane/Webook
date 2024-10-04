package Ioc

import (
	"GinStart/Repository/Dao"
	"GinStart/pkg/logger"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func InitDB(l logger.Logger) *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}

	var cfg = Config{
		DSN: "root:aaa@tcp(localhost:3306)/ginstart",
	}
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
			//慢查询日志    0 为所有日志都打出来
			SlowThreshold: 0,
			LogLevel:      glogger.Info,
		}),
	})
	if err != nil {
		panic(err)
	}

	err1 := Dao.InitTables(db)
	if err1 != nil {
		panic(err1)
	}
	return db
}

// 函数衍生类型实现接口
type gormLoggerFunc func(msg string, fields ...logger.Field)

func (f gormLoggerFunc) Printf(msg string, v ...interface{}) {
	f(msg, logger.Field{
		Key:   "args",
		Value: v,
	})
}
