package integration

import (
	"GinStart/Repository/Dao"
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
	s.server = startup.InitWireServer()
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
				Msg:  "success",
				Data: 1,
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
	Title   string `json:"title"`
	Content string `json:"content"`
}
