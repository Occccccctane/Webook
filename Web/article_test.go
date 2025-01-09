package Handler

import (
	"GinStart/Domain"
	"GinStart/Service"
	svcmock "GinStart/Service/mocks"
	ijwt "GinStart/Web/Jwt"
	"GinStart/pkg/logger"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArticleHandler_Publish(t *testing.T) {
	testCases := []struct {
		name         string
		mock         func(ctrl *gomock.Controller) Service.ArticleService
		reqBody      string
		ExpectedCode int
		ExpectedReq  Result
	}{
		{
			name: "新建并发布成功",
			mock: func(ctrl *gomock.Controller) Service.ArticleService {
				svc := svcmock.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), Domain.Article{
					Content: "测试内容",
					Title:   "测试标题",
					Author: Domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				return svc
			},
			reqBody: `{
				"title": "测试标题",
				"content": "测试内容"
			}`,
			ExpectedCode: http.StatusOK,
			ExpectedReq: Result{
				Code: 200,
				Msg:  "保存成功",
				Data: float64(1),
			},
		},
		{
			name: "帖子已存在但发表成功",
			mock: func(ctrl *gomock.Controller) Service.ArticleService {
				svc := svcmock.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), Domain.Article{
					Id:      12,
					Content: "测试内容",
					Title:   "测试标题",
					Author: Domain.Author{
						Id: 123,
					},
				}).Return(int64(12), nil)
				return svc
			},
			reqBody: `{
				"id": 12,
				"title": "测试标题",
				"content": "测试内容"
			}`,
			ExpectedCode: http.StatusOK,
			ExpectedReq: Result{
				Code: 200,
				Msg:  "保存成功",
				Data: float64(12),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			//创建控制器
			artSvc := tc.mock(ctrl)
			hdl := NewArticleHandler(artSvc, logger.NewNopLogger())

			//准备服务器,注册路由
			server := gin.Default()
			server.Use(func(c *gin.Context) {
				c.Set("user", &ijwt.UserClaims{
					Uid: 123,
				})
			})
			hdl.RegisterRoute(server)

			//准备请求和记录的recorder
			//序列化传入json数据
			req, err := http.NewRequest(http.MethodPost, "/articles/publish",
				bytes.NewBufferString(tc.reqBody))

			// 设置请求头和接收器
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()

			//服务器接收响应
			server.ServeHTTP(recorder, req)

			//	断言结果
			var res Result
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.ExpectedCode, recorder.Code)
			assert.Equal(t, tc.ExpectedReq, res)
		})
	}
}
