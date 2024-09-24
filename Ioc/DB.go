package Ioc

import (
	"GinStart/Repository/Dao"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
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

	db, err := gorm.Open(mysql.Open(cfg.DSN))
	if err != nil {
		panic(err)
	}

	err1 := Dao.InitTables(db)
	if err1 != nil {
		panic(err1)
	}
	return db
}
