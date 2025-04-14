package startup

import (
	"context"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

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

	return client.Database("webook")
}

func InitSnowflake() *snowflake.Node {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	return node
}
