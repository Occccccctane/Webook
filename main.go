package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
	"log"
	"time"
)

func main() {
	InitViperWatch()
	//InitViperRemoteWatch()
	server := InitWireServer()

	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func InitLogger() {
	//生产环境日志模板
	//logger , err := zap.NewProduction()
	//开发环境
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}
func InitViper() {
	// 为方便更换配置路径，使用命令行参数
	cf := pflag.String("config", "Config/dev.yaml", "配置文件路径")
	// 用于解析参数，这一步以后，变量才有值
	pflag.Parse()
	viper.SetConfigType("yaml")
	viper.SetConfigFile(*cf)
	//读取配置
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func InitViperRemote() {
	err := viper.AddRemoteProvider("etcd3", "http://127.0.0.1:12379", "/Ginstart")
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

func InitViperWatch() {
	cf := pflag.String("config", "Config/dev.yaml", "配置文件路径")
	// 用于解析参数，这一步以后，变量才有值
	pflag.Parse()
	viper.SetConfigType("yaml")
	viper.SetConfigFile(*cf)
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println(viper.Get("test.name"))
	})
	//读取配置
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func InitViperRemoteWatch() {
	//不推荐使用，若是WatchRemoteConfig监听提前会导致并发问题
	err := viper.AddRemoteProvider("etcd3", "http://127.0.0.1:12379", "/Ginstart")
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("远程配置文件发生改变")
	})
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			err := viper.WatchRemoteConfig()
			if err != nil {
				panic(err)
			}
			log.Println("Remote:", viper.Get("test.name"))
			time.Sleep(time.Second * 5)
		}
	}()
}
