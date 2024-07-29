package integration

import (
	Handler "GinStart/Web"
	"GinStart/integration/startup"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}
func TestUserHandler_SendSMS(t *testing.T) {
	rdb := startup.InitRedis()
	server := startup.InitWireServer()

	testCase := []struct {
		name string
		//准备数据和验证清理
		before func(t *testing.T)
		after  func(t *testing.T)

		phone        string
		ExpectedCode int
		ExpectedBody Handler.Result
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {
				//	不需要
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()

				key := "phone_code:login:13333333333"
				code, err := rdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, len(code) > 0)

				dur, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, dur > time.Minute*9)
				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},

			phone:        "13333333333",
			ExpectedCode: http.StatusOK,
			ExpectedBody: Handler.Result{Msg: "发送成功"},
		},
		{
			name: "手机号为空",
			before: func(t *testing.T) {
				//	不需要
			},
			after: func(t *testing.T) {
			},

			phone:        "",
			ExpectedCode: http.StatusOK,
			ExpectedBody: Handler.Result{Code: 400, Msg: "电话格式错误"},
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:13333333333"
				err := rdb.Set(ctx, key, "123456", time.Minute*9+time.Second*50).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()

				key := "phone_code:login:13333333333"
				code, err := rdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "123456", code)
				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},

			phone:        "13333333333",
			ExpectedCode: http.StatusOK,
			ExpectedBody: Handler.Result{Code: 400, Msg: "短信发送太频繁"},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			reqString := fmt.Sprintf(`{"phone": "%s"}`, tc.phone)
			req, err := http.NewRequest(http.MethodPost, "/users/login_sms/code/send",
				bytes.NewReader([]byte(reqString)))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.ExpectedCode, recorder.Code)
			var red Handler.Result
			err = json.NewDecoder(recorder.Body).Decode(&red)
			assert.NoError(t, err)
			assert.Equal(t, tc.ExpectedBody, red)
		})
	}
}
