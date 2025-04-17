package integration

import (
	"GinStart/Repository/Dao"
	ijwt "GinStart/Web/Jwt"
	"GinStart/integration/startup"
	"bytes"
	"context"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type ArticleMongoDBHandlerSuite struct {
	suite.Suite
	db      *mongo.Database
	col     *mongo.Collection
	liveCol *mongo.Collection
	server  *gin.Engine
}

func TestArticleMongoDBHandler(t *testing.T) {
	suite.Run(t, &ArticleMongoDBHandlerSuite{})
}
func (s *ArticleMongoDBHandlerSuite) SetupSuite() {
	s.db = startup.InitMongoDB()
	err := Dao.InitCollection(s.db)
	if err != nil {
		return
	}
	s.col = s.db.Collection("articles")
	s.liveCol = s.db.Collection("published_articles")
	node, err := snowflake.NewNode(1)
	assert.NoError(s.T(), err)

	server := gin.Default()
	hdl := startup.InitArticleHandler(Dao.NewMongoDBArticleDao(s.db, node))
	server.Use(func(ctx *gin.Context) {
		ctx.Set("user", ijwt.UserClaims{
			Uid: 123,
		})
	})
	hdl.RegisterRoute(server)
	s.server = server
}
func (s *ArticleMongoDBHandlerSuite) TearDownTest() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := s.col.DeleteMany(ctx, bson.D{})
	assert.NoError(s.T(), err)
	_, err = s.liveCol.DeleteMany(ctx, bson.D{})
	assert.NoError(s.T(), err)
}

func (s *ArticleMongoDBHandlerSuite) TestEdit() {
	t := s.T()
	testCase := []struct {
		name         string
		before       func(t *testing.T)
		after        func(t *testing.T)
		art          Article
		ExpectedCode int
		ExpectedReq  Result[int64]
	}{
		{
			name:   "新建成功",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				//	验证保存成功
				var art Dao.Article
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := s.col.
					FindOne(ctx, bson.D{bson.E{
						Key:   "author_id",
						Value: 123,
					}}).
					Decode(&art)
				assert.NoError(t, err)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				assert.True(t, art.Id > 0)
				art.Utime = 0
				art.Ctime = 0
				art.Id = 0
				assert.Equal(t, Dao.Article{
					Title:    "测试标题",
					Content:  "测试内容",
					AuthorId: 123,
					Status:   1,
				}, art)
				s.TearDownTest()
			},
			art: Article{
				Title:   "测试标题",
				Content: "测试内容",
			},
			ExpectedCode: http.StatusOK,
			ExpectedReq: Result[int64]{
				Code: 200,
				Msg:  "保存成功",
				Data: 1,
			},
		},
		{
			name: "修改帖子",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				_, err := s.col.InsertOne(ctx, Dao.Article{
					Id:       2,
					AuthorId: 123,
					Title:    "测试标题",
					Content:  "测试内容",
					Ctime:    456,
					Utime:    789,
					// 假设已经发布
					Status: 2,
				})
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				//	验证保存成功
				var art Dao.Article
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := s.col.
					FindOne(ctx, bson.D{bson.E{
						Key:   "id",
						Value: 2,
					}}).
					Decode(&art)
				assert.NoError(t, err)
				assert.True(t, art.Utime > 789)
				art.Utime = 0
				assert.Equal(t, Dao.Article{
					Id:       2,
					Title:    "新的标题",
					Content:  "新的内容",
					AuthorId: 123,
					Status:   1,
					Ctime:    456,
				}, art)
				s.TearDownTest()
			},
			art: Article{
				Id:      2,
				Title:   "新的标题",
				Content: "新的内容",
			},
			ExpectedCode: http.StatusOK,
			ExpectedReq: Result[int64]{
				Code: 200,
				Msg:  "保存成功",
				Data: 2,
			},
		},
		{
			name: "修改别人帖子",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				_, err := s.col.InsertOne(ctx, Dao.Article{
					Id:       3,
					AuthorId: 234,
					Title:    "测试标题",
					Content:  "测试内容",
					Ctime:    456,
					Utime:    789,
					Status:   1,
				})
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				//	验证数据不变
				var art Dao.Article
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := s.col.FindOne(ctx, bson.D{bson.E{
					Key:   "id",
					Value: 3,
				}}).Decode(&art)
				assert.NoError(t, err)
				assert.Equal(t, Dao.Article{
					Id:       3,
					Title:    "测试标题",
					Content:  "测试内容",
					AuthorId: 234,
					Status:   1,
					Ctime:    456,
					Utime:    789,
				}, art)
				s.TearDownTest()
			},
			art: Article{
				Id:      3,
				Title:   "新的标题",
				Content: "新的内容",
			},
			ExpectedCode: http.StatusInternalServerError,
			ExpectedReq: Result[int64]{
				Code: 500,
				Msg:  "系统错误",
				Data: 0,
			},
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			//序列化传入json数据
			reqString, err := json.Marshal(tc.art)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/articles/edit",
				bytes.NewReader(reqString))

			// 设置请求头和接收器
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()

			s.server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.ExpectedCode, recorder.Code)
			var res Result[int64]
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			if tc.ExpectedReq.Data > 0 {
				assert.True(t, res.Data > 0)
			}
		})
	}
}

func (s *ArticleMongoDBHandlerSuite) TestArticle_Publish() {
	t := s.T()

	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)
		req    Article

		// 预期响应
		wantCode   int
		wantResult Result[int64]
	}{
		{
			name: "新建帖子并发表",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
				// 验证一下数据
				var art Dao.Article
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := s.col.FindOne(ctx, bson.D{bson.E{
					Key:   "author_id",
					Value: 123,
				}}).Decode(&art)
				assert.NoError(t, err)
				assert.Equal(t, "hello，你好", art.Title)
				assert.Equal(t, "随便试试", art.Content)
				assert.Equal(t, int64(123), art.AuthorId)
				assert.Equal(t, uint8(2), art.Status)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				var publishedArt Dao.PublishedArticle
				err = s.liveCol.FindOne(ctx, bson.D{bson.E{
					Key:   "author_id",
					Value: 123,
				}}).Decode(&publishedArt)
				assert.NoError(t, err)
				assert.Equal(t, "hello，你好", publishedArt.Title)
				assert.Equal(t, "随便试试", publishedArt.Content)
				assert.Equal(t, int64(123), publishedArt.AuthorId)
				assert.Equal(t, uint8(2), publishedArt.Status)
				assert.True(t, publishedArt.Ctime > 0)
				assert.True(t, publishedArt.Utime > 0)
			},
			req: Article{
				Title:   "hello，你好",
				Content: "随便试试",
			},
			wantCode: 200,
			wantResult: Result[int64]{
				Code: 200,
				Msg:  "保存成功",
				Data: 1,
			},
		},
		{
			// 制作库有，但是线上库没有
			name: "更新帖子并新发表",
			before: func(t *testing.T) {
				// 模拟已经存在的帖子
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				_, err := s.col.InsertOne(ctx, &Dao.Article{
					Id:       2,
					Title:    "我的标题",
					Content:  "我的内容",
					Ctime:    456,
					Status:   1,
					Utime:    234,
					AuthorId: 123,
				})
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 验证一下数据
				var art Dao.Article
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := s.col.FindOne(ctx, bson.M{
					"id": 2,
				}).Decode(&art)
				assert.NoError(t, err)
				assert.Equal(t, "新的标题", art.Title)
				assert.Equal(t, "新的内容", art.Content)
				assert.Equal(t, uint8(2), art.Status)
				assert.Equal(t, int64(123), art.AuthorId)
				// 创建时间没变
				assert.Equal(t, int64(456), art.Ctime)
				// 更新时间变了
				assert.True(t, art.Utime > 234)
				var publishedArt Dao.PublishedArticle
				err = s.liveCol.FindOne(ctx, bson.M{
					"id": 2,
				}).Decode(&publishedArt)
				assert.NoError(t, err)
				assert.Equal(t, "新的标题", art.Title)
				assert.Equal(t, "新的内容", art.Content)
				assert.Equal(t, int64(123), art.AuthorId)
				assert.True(t, publishedArt.Ctime > 0)
				assert.Equal(t, uint8(2), publishedArt.Status)
				assert.True(t, publishedArt.Utime > 0)
			},
			req: Article{
				Id:      2,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: 200,
			wantResult: Result[int64]{
				Code: 200,
				Msg:  "保存成功",
				Data: 2,
			},
		},
		{
			name: "更新帖子，并且重新发表",
			before: func(t *testing.T) {
				art := Dao.Article{
					Id:       3,
					Title:    "我的标题",
					Content:  "我的内容",
					Ctime:    456,
					Status:   1,
					Utime:    234,
					AuthorId: 123,
				}
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				_, err := s.col.InsertOne(ctx, &art)
				assert.NoError(t, err)
				part := Dao.PublishedArticle(art)
				_, err = s.liveCol.InsertOne(ctx, &part)
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var art Dao.Article
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				res := s.col.FindOne(ctx, bson.M{
					"id": 3,
				})
				err := res.Decode(&art)
				assert.NoError(t, err)
				assert.Equal(t, "新的标题", art.Title)
				assert.Equal(t, "新的内容", art.Content)
				assert.Equal(t, int64(123), art.AuthorId)
				assert.Equal(t, uint8(2), art.Status)
				// 创建时间没变
				assert.Equal(t, int64(456), art.Ctime)
				// 更新时间变了
				assert.True(t, art.Utime > 234)

				var part Dao.PublishedArticle
				err = s.liveCol.FindOne(ctx, bson.M{
					"id": 3,
				}).Decode(&part)
				assert.NoError(t, err)
				assert.Equal(t, "新的标题", part.Title)
				assert.Equal(t, "新的内容", part.Content)
				assert.Equal(t, int64(123), part.AuthorId)
				assert.Equal(t, uint8(2), part.Status)
				// 创建时间没变
				assert.Equal(t, int64(456), part.Ctime)
				// 更新时间变了
				assert.True(t, part.Utime > 234)
			},
			req: Article{
				Id:      3,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: 200,
			wantResult: Result[int64]{
				Code: 200,
				Msg:  "保存成功",
				Data: 3,
			},
		},
		{
			name: "更新别人的帖子，并且发表失败",
			before: func(t *testing.T) {
				art := Dao.Article{
					Id:      4,
					Title:   "我的标题",
					Content: "我的内容",
					Ctime:   456,
					Utime:   234,
					Status:  1,
					// 注意。这个 AuthorID 我们设置为另外一个人的ID
					AuthorId: 789,
				}
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				_, err := s.col.InsertOne(ctx, &art)
				assert.NoError(t, err)
				part := Dao.PublishedArticle(Dao.Article{
					Id:       4,
					Title:    "我的标题",
					Content:  "我的内容",
					Ctime:    456,
					Status:   2,
					Utime:    234,
					AuthorId: 789,
				})
				_, err = s.liveCol.InsertOne(ctx, &part)
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 更新应该是失败了，数据没有发生变化
				var art Dao.Article
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := s.liveCol.FindOne(ctx, bson.M{
					"id": 4,
				}).Decode(&art)
				assert.NoError(t, err)
				assert.Equal(t, "我的标题", art.Title)
				assert.Equal(t, "我的内容", art.Content)
				assert.Equal(t, int64(456), art.Ctime)
				assert.Equal(t, int64(234), art.Utime)
				assert.Equal(t, uint8(1), art.Status)
				assert.Equal(t, int64(789), art.AuthorId)

				var part Dao.PublishedArticle
				// 数据没有变化
				err = s.liveCol.FindOne(ctx, bson.M{
					"id": 4,
				}).Decode(&part)
				assert.NoError(t, err)
				assert.Equal(t, "我的标题", part.Title)
				assert.Equal(t, "我的内容", part.Content)
				assert.Equal(t, int64(789), part.AuthorId)
				assert.Equal(t, uint8(2), part.Status)
				// 创建时间没变
				assert.Equal(t, int64(456), part.Ctime)
				// 更新时间变了
				assert.Equal(t, int64(234), part.Utime)
			},
			req: Article{
				Id:      4,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: 500,
			wantResult: Result[int64]{
				Code: 500,
				Msg:  "系统错误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			data, err := json.Marshal(tc.req)
			// 不能有 error
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost,
				"/articles/publish", bytes.NewReader(data))
			assert.NoError(t, err)
			req.Header.Set("Content-Type",
				"application/json")
			recorder := httptest.NewRecorder()

			s.server.ServeHTTP(recorder, req)
			code := recorder.Code
			assert.Equal(t, tc.wantCode, code)
			if code != http.StatusOK {
				return
			}
			// 反序列化为结果
			// 利用泛型来限定结果必须是 int64
			var result Result[int64]
			err = json.Unmarshal(recorder.Body.Bytes(), &result)
			assert.NoError(t, err)
			if tc.wantResult.Data > 0 {
				assert.True(t, result.Data > 0)
			}
			tc.after(t)
		})
	}
}
