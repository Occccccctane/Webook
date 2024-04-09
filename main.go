package main

import (
	"GinStart/MiddleWare"
	"GinStart/Repository"
	"GinStart/Repository/Dao"
	"GinStart/Service"
	"GinStart/Web"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	db := InitDB()
	server := InitServer()
	//初始化
	InitUserHdl(db, server)
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

	c.Use(
		cors.New(cors.Config{
			//AllowAllOrigins: true,允许所有域名，比较危险
			//AllowedOrigins: []string{"http://localhost:3000","http://aaa"},枚举允许的域名

			AllowCredentials: true,                     //是否允许带cookie等用户凭据，正常都需要允许
			AllowedHeaders:   []string{"content-type"}, // 允许的请求头

			//AllowedMethods: []string{"POST"}, 允许的请求方法，最好不配置

			//允许字符串的检查方法，如果传入的字符串包含相关的字段则放行
			AllowOriginFunc: func(origin string) bool {

				if strings.Contains(origin, "localhost") { //判断包含该字段
					// if strings.HasPrefix(origin, "http://localhost")判定包含前缀
					return true
				}
				return strings.Contains(origin, "xxx.com") //返回一个表达式，上面用于判断是否是本机调试，下面用于判断是否是公司的域名
			},

			MaxAge: 12 * time.Hour, //检测时间长度
		}),
		// 登录校验
		sessions.Sessions("ssid", store),
		login.CheckLogin(),
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
