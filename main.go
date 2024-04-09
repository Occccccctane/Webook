package main

import (
	"GinStart/MiddleWare"
	"GinStart/Repository"
	"GinStart/Repository/Dao"
	"GinStart/Service"
	"GinStart/Web"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := InitDB()
	server := InitServer()
	//初始化
	InitUserHdl(db, server)

	//useSessionCheck(server)
	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/Gin"))
	if err != nil {
		panic(err)
	}
	err1 := Dao.InitTables(db)
	if err1 != nil {
		panic(err1)
	}
	return db
}

func InitServer() *gin.Engine {
	c := gin.Default()

	cross := &MiddleWare.CrossDomain{}
	//登录验证
	useJWTCheck(c)
	c.Use(
		//跨域
		cross.CrossDomainHandler(),
	)
	return c
}

func InitUserHdl(db *gorm.DB, server *gin.Engine) {
	ud := Dao.NewUserDao(db)
	ur := Repository.NewUserRepository(ud)
	us := Service.NewUserService(ur)
	hdl := Handler.NewUserHandler(us)
	hdl.RegisterRoute(server)
}

func useJWTCheck(server *gin.Engine) {
	login := &MiddleWare.LoginJWTBuilder{}
	server.Use(login.CheckLogin())
}

func useSessionCheck(server *gin.Engine) {
	login := &MiddleWare.LoginBuilder{}
	//使用Cookie存储
	//store := cookie.NewStore([]byte("secret"))
	//使用内存作为存储Session信息的载体
	//store := memstore.NewStore([]byte("uT4WEEyz2oSIkzEVwSBwIpwMoHGZ70N2FNtXJFDfCsRVUFa0PKU53UfzCdwZs8I8"),
	//	[]byte("VgNulil8EVFkEi7nWNozTqHr7bVwreh9Pn4CPvCPvjhpEDzLMVYKeCaXQePKnBxW"))
	//使用redis存储
	store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
		[]byte("ppSik8fZfCugefcqWNeh54adKgtN1Fmp"),
		[]byte("zysHnwMiU2jPJ59NmrmBLQcZNT3FPysv"),
	)
	if err != nil {
		panic(err)
	}
	server.Use(
		// 登录校验
		sessions.Sessions("ssid", store),
		login.CheckLogin(),
	)
}
