package Handler

import (
	"GinStart/Domain"
	"GinStart/Service"
	svcmock "GinStart/Service/mocks"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_Signup(t *testing.T) {
	TestCase := []struct {
		name string
		//准备两个service实例，userHandler用到了两个service，用mock生成出来
		mock func(ctrl *gomock.Controller) (Service.UserService, Service.CodeService)
		//构造请求，服务器预期收到的请求
		reqBuilder func(t *testing.T) *http.Request
		//预期输出
		ExpectedCode int
		ExpectedBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) (Service.UserService, Service.CodeService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), Domain.User{
					Email:    "1333@163.com",
					Password: "Aaa123123123",
				}).Return(nil)
				codeSvc := svcmock.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
				"email":"1333@163.com",
				"password":"Aaa123123123",
				"confirmPassword":"Aaa123123123"
			}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			ExpectedCode: http.StatusOK,
			ExpectedBody: `{"code":"200"}`,
		},
		{
			name: "Bind错误",
			mock: func(ctrl *gomock.Controller) (Service.UserService, Service.CodeService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				codeSvc := svcmock.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`
				"email":"1333@163.com",
				"password":"Aaa1231
				"confirmPassword":"Aaa123123123"
			}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			ExpectedCode: http.StatusBadRequest,
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) (Service.UserService, Service.CodeService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				codeSvc := svcmock.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
				"email":"1333163.com",
				"password":"Aaa123123123",
				"confirmPassword":"Aaa123123123"
			}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: `{"code":"400","msg":"邮箱格式错误"}`,
		},
		//{
		//	name: "邮箱冲突",
		//	mock: func(ctrl *gomock.Controller) (Service.UserService, Service.CodeService) {
		//		userSvc := svcmock.NewMockUserService(ctrl)
		//		userSvc.EXPECT().Signup(gomock.Any(), Domain.User{
		//			Email:    "1333@163.com",
		//			Password: "Aaa123123123",
		//		}).Return(Service.ErrUserUnique)
		//		codeSvc := svcmock.NewMockCodeService(ctrl)
		//		return userSvc, codeSvc
		//	},
		//	reqBuilder: func(t *testing.T) *http.Request {
		//		req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
		//		"email":"1333@163.com",
		//		"password":"Aaa123123123",
		//		"confirmPassword":"Aaa123123123"
		//	}`)))
		//		req.Header.Set("Content-Type", "application/json")
		//		assert.NoError(t, err)
		//
		//		return req
		//	},
		//	ExpectedCode: http.StatusInternalServerError,
		//	ExpectedBody: `{"code":"500","msg":"注册失败"}`,
		//},
	}
	for _, tc := range TestCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			//创建控制器
			userSvc, codeSvc := tc.mock(ctrl)
			hdl := NewUserHandler(userSvc, codeSvc)

			//准备服务器,注册路由
			server := gin.Default()
			hdl.RegisterRoute(server)

			//准备请求和记录的recorder
			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()

			//服务器接收响应
			server.ServeHTTP(recorder, req)
			log.Println(recorder.Code)
			log.Println(recorder.Body)

			//	断言结果
			assert.Equal(t, tc.ExpectedCode, recorder.Code)
			assert.Equal(t, tc.ExpectedBody, recorder.Body.String())
		})
	}
}

func TestEmailPatten(t *testing.T) {
	TestCase := []struct {
		name  string
		email string
		match bool
	}{
		{
			name:  "不带@",
			email: "123123",
			match: false,
		},
		{
			name:  "带@没后缀",
			email: "45457asd@",
			match: false,
		},
		{
			name:  "正确",
			email: "123123@122.com",
			match: true,
		},
	}
	h := NewUserHandler(nil, nil)
	for _, tc := range TestCase {
		t.Run(tc.name, func(t *testing.T) {
			match, err := h.emailRexExp.MatchString(tc.email)
			require.NoError(t, err)
			assert.Equal(t, tc.match, match)
		})
	}
}

//func TestMock(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	//	mock实现，模拟实现
//	userSvc := svcmock.NewMockUserService(ctrl)
//	userSvc.EXPECT().Signup(gomock.Any(), Domain.User{
//		Id:    1,
//		Email: "123@qq.com",
//	}).Return(nil)
//}
