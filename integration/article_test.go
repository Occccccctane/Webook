package integration

import (
	"GinStart/Repository/Dao"
	ijwt "GinStart/Web/Jwt"
	"GinStart/integration/startup"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ArticleHandlerSuite struct {
	suite.Suite
	db     *gorm.DB
	server *gin.Engine
}

func TestArticleHandler(t *testing.T) {
	suite.Run(t, &ArticleHandlerSuite{})
}
func (s *ArticleHandlerSuite) SetupSuite() {
	s.db = startup.InitDB()
	server := gin.Default()
	hdl := startup.InitArticleHandler()
	server.Use(func(ctx *gin.Context) {
		ctx.Set("user", ijwt.UserClaims{
			Uid: 123,
		})
	})
	hdl.RegisterRoute(server)
	s.server = server
}
func (s *ArticleHandlerSuite) TearDownTest() {
	s.db.Exec("truncate table `articles`")
}

func (s *ArticleHandlerSuite) TestEdit() {
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
				err := s.db.Where("author_id = ?", 123).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				assert.True(t, art.Id > 0)
				assert.Equal(t, "测试标题", art.Title)
				assert.Equal(t, "测试内容", art.Content)
				assert.Equal(t, int64(123), art.AuthorId)
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
				err := s.db.Create(&Dao.Article{
					Id:       2,
					AuthorId: 123,
					Title:    "测试标题",
					Content:  "测试内容",
					Ctime:    456,
					Utime:    789,
					Status:   1,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				//	验证保存成功
				var art Dao.Article
				err := s.db.Where("id = ?", 2).First(&art).Error
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
				err := s.db.Create(&Dao.Article{
					Id:       3,
					AuthorId: 234,
					Title:    "测试标题",
					Content:  "测试内容",
					Ctime:    456,
					Utime:    789,
					Status:   1,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				//	验证数据不变
				var art Dao.Article
				err := s.db.Where("id = ?", 3).First(&art).Error
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
			var red Result[int64]
			err = json.NewDecoder(recorder.Body).Decode(&red)
			assert.NoError(t, err)
			assert.Equal(t, tc.ExpectedReq, red)
		})
	}
}

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"date"`
}

type Article struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
