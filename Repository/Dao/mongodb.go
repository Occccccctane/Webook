package Dao

import (
	"context"
	"errors"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

type MongoDBArticleDao struct {
	node    *snowflake.Node
	col     *mongo.Collection
	liveCol *mongo.Collection
}

func NewMongoDBArticleDao(mdb *mongo.Database, node *snowflake.Node) ArticleDao {
	return &MongoDBArticleDao{
		node:    node,
		col:     mdb.Collection("articles"),
		liveCol: mdb.Collection("published_articles"),
	}
}
func (m *MongoDBArticleDao) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	art.Id = m.node.Generate().Int64()
	_, err := m.col.InsertOne(ctx, &art)
	return art.Id, err
}

func (m *MongoDBArticleDao) Update(ctx context.Context, art Article) error {
	filter := map[string]interface{}{
		"author_id": art.AuthorId,
		"id":        art.Id,
	}
	now := time.Now().UnixMilli()
	set := bson.M{
		"$set": map[string]interface{}{
			"title":   art.Title,
			"content": art.Content,
			"status":  art.Status,
			"utime":   now,
		},
	}
	res, err := m.col.UpdateOne(ctx, filter, set)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.New("更新数据失败")
	}
	return nil
}

func (m *MongoDBArticleDao) Sync(ctx context.Context, art Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBArticleDao) SyncStatus(ctx context.Context, uid int64, aid int64, status uint8) error {
	//TODO implement me
	panic("implement me")
}
