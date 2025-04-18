package Ioc

import (
	"GinStart/Repository/Dao"
	"GinStart/pkg/logger"
	"context"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func InitMongoDB() *mongo.Database {
	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
			fmt.Println(evt.Command)
		},
	}

	opt := options.Client().
		ApplyURI("mongodb://root:root@127.0.0.1:27017/").
		SetMonitor(monitor)

	client, err := mongo.Connect(opt)
	if err != nil {
		panic(err)
	}

	//初始化集合
	err = Dao.InitCollection(client.Database("webook"))
	if err != nil {
		panic(err)
	}
	return client.Database("webook")
}

func InitSnowflake() *snowflake.Node {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	return node
}

// 函数衍生类型实现接口
type gormLoggerFunc func(msg string, fields ...logger.Field)

func (f gormLoggerFunc) Printf(msg string, v ...interface{}) {
	f(msg, logger.Field{
		Key:   "args",
		Value: v,
	})
}
