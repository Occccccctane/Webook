package Ioc

import (
	"GinStart/Config"
	"GinStart/Repository/Dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(Config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}

	err1 := Dao.InitTables(db)
	if err1 != nil {
		panic(err1)
	}
	return db
}
