package startup

import (
	"GinStart/Repository/Dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:aaa@tcp(localhost:3306)/ginstart"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = Dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
