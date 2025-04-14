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

	var testCases []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)
		req    Article

		// 预期响应
		wantCode   int
		wantResult Result[int64]
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
			assert.Equal(t, tc.wantResult, result)
			tc.after(t)
		})
	}
}
