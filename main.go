package main

import (
	"GinStart/Config"
	"GinStart/MiddleWare"
	"GinStart/Repository"
	"GinStart/Repository/Dao"
	"GinStart/Service"
	"GinStart/Web"
	"GinStart/pkg/middleware/ratelimit"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func main() {
	db := InitDB()
	server := InitServer()
	//初始化
	InitUserHdl(db, server)
	//server := gin.Default()
	//server.GET("/hello", func(context *gin.Context) {
	//	context.String(http.StatusOK, "登录成功")
	//})
	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}

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

func InitServer() *gin.Engine {
	c := gin.Default()

	cross := &MiddleWare.CrossDomain{}
	// 登录验证
	useJWTCheck(c)
	// useSessionCheck(server)

	// 创建一个Redis服务器
	redisClient := redis.NewClient(&redis.Options{
		Addr: Config.Config.Redis.Addr,
	})
	c.Use(
		//跨域
		cross.CrossDomainHandler(),
		ratelimit.NewBuilder(redisClient, time.Second, 1).Build(),
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
	store := cookie.NewStore([]byte("secret"))
	//使用内存作为存储Session信息的载体
	//store := memstore.NewStore([]byte("uT4WEEyz2oSIkzEVwSBwIpwMoHGZ70N2FNtXJFDfCsRVUFa0PKU53UfzCdwZs8I8"),
	//	[]byte("VgNulil8EVFkEi7nWNozTqHr7bVwreh9Pn4CPvCPvjhpEDzLMVYKeCaXQePKnBxW"))
	//使用redis存储
	//store, err := redis.NewStore(16, "tcp",Config.Config.Redis.Addr, "",
	//	[]byte("ppSik8fZfCugefcqWNeh54adKgtN1Fmp"),
	//	[]byte("zysHnwMiU2jPJ59NmrmBLQcZNT3FPysv"),
	//)
	//if err != nil {
	//	panic(err)
	//}
	server.Use(
		// 登录校验
		sessions.Sessions("ssid", store),
		login.CheckLogin(),
	)
}
