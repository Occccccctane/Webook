package Handler

import (
	"GinStart/Domain"
	"GinStart/Service"
	ijwt "GinStart/Web/Jwt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
)

// 正则常量
const (
	emailRegex    = "^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$"
	passwordRegex = "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)[a-zA-Z\\d]{8,72}$"
	phoneRegex    = "^1[3-9]\\d{9}$"

	bizLogin = "login"
)

func (h *UserHandler) RegisterRoute(server *gin.Engine) {

	user := server.Group("/users")
	user.POST("/signup", h.Signup)
	user.POST("/login", h.Login)
	user.POST("logout", h.Logout)
	user.POST("/edit", h.Edit)
	user.GET("/profile", h.Profile)
	user.GET("/refresh_token", h.RefreshToken)

	//	验证码相关
	user.POST("/login_sms/code/send", h.LoginSMSCode)
	user.POST("/login_sms", h.LoginSMS)
}

type UserHandler struct {
	ijwt.Handler
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	phoneRexExp    *regexp.Regexp
	svc            Service.UserService
	codeSvc        Service.CodeService
}

// NewUserHandler 正则预加载
func NewUserHandler(svc Service.UserService, codeSvc Service.CodeService,
	hdl ijwt.Handler) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegex, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegex, regexp.None),
		phoneRexExp:    regexp.MustCompile(phoneRegex, regexp.None),
		svc:            svc,
		codeSvc:        codeSvc,
		Handler:        hdl,
	}
}

func (h *UserHandler) Signup(c *gin.Context) {

	type signUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req signUpReq
	err1 := c.Bind(&req)
	if err1 != nil {
		return
	}

	// 校验邮箱格式
	isEmailTrue, err := h.emailRexExp.MatchString(req.Email)
	if err != nil {
		c.JSON(http.StatusNotImplemented, gin.H{
			"code": "501",
			"msg":  "系统错误",
		})
		return
	}
	if !isEmailTrue {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "邮箱格式错误",
		})
		return
	}

	//校验密码
	isPasswordTrue, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		c.JSON(http.StatusNotImplemented, gin.H{
			"code": "501",
			"msg":  "系统错误",
		})
		return
	}
	if !isPasswordTrue {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "密码格式错误，应包括大小写字母和数字，并大于8位",
		})
		return
	}

	//校验两次密码
	if req.ConfirmPassword != req.Password {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "两次密码不一致",
		})
		return
	}

	//service层逻辑调用
	err = h.svc.Signup(c, Domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	//错误处理
	switch err {
	case nil:
		c.JSON(http.StatusOK, gin.H{
			"code": "200",
		})
	case Service.ErrUserUnique:
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": "500",
			"msg":  "邮箱已注册",
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": "500",
			"msg":  "注册失败",
		})

	}

}

func (h *UserHandler) Login(ctx *gin.Context) {
	type logINReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req logINReq
	err1 := ctx.Bind(&req)
	if err1 != nil {
		return
	}

	u, err2 := h.svc.Login(ctx, req.Email, req.Password)
	switch err2 {
	case nil:
		//sess := sessions.Default(ctx)
		//sess.Set("UserId", u.Id)
		//sess.Options(sessions.Options{
		//	MaxAge:   900, //15分钟
		//	HttpOnly: true,
		//})
		//err := sess.Save()
		//if err != nil {
		//	ctx.JSON(http.StatusInternalServerError, gin.H{
		//		"code": "500",
		//		"msg":  "系统错误",
		//	})
		//	return
		//}

		//换成JWT处理
		h.SetLoginToken(ctx, u.Id)
		ctx.JSON(http.StatusOK, gin.H{
			"code": "200",
		})
	case Service.ErrInvalidUserOrPassword:
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "账号或密码错误",
		})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": "500",
			"msg":  "系统错误",
		})
	}

}

func (h *UserHandler) Edit(c *gin.Context) {
	type editReq struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		NewPassword string `json:"newPassword"`
		Nickname    string `json:"nickname"`
		Birthday    string `json:"birthday"`
		Info        string `json:"info"`
	}

	var req editReq
	err1 := c.Bind(&req)
	if err1 != nil {
		return
	}

	//校验密码
	isPasswordTrue, err := h.passwordRexExp.MatchString(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusNotImplemented, gin.H{
			"code": "501",
			"msg":  "系统错误",
		})
		return
	}
	if !isPasswordTrue {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "密码格式错误，应包括大小写字母和数字，并大于8位",
		})
		return
	}
	if len(req.Nickname) > 15 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "信息不能大于50位",
		})
		return
	}
	if len(req.Info) > 50 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "信息不能大于50位",
		})
		return
	}

	err2 := h.svc.Edit(c, req.NewPassword, Domain.User{
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Nickname,
		Birthday: req.Birthday,
		Info:     req.Info,
	})

	switch err2 {
	case nil:
		c.JSON(http.StatusOK, gin.H{
			"code": "200",
		})
	case Service.ErrInvalidUserOrPassword:
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "账号或密码错误",
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": "500",
			"msg":  "系统错误",
		})
	}
}

func (h *UserHandler) Profile(c *gin.Context) {
	//从上下文取出，断言为UserClaims类型
	us, ok := c.Get("user")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	us1, _ := us.(ijwt.UserClaims)

	u, _ := h.svc.FindById(c, us1.Uid)
	c.JSON(http.StatusOK, gin.H{
		"code":     "200",
		"Id":       u.Id,
		"Email":    u.Email,
		"Nickname": u.Nickname,
		"Birthday": u.Birthday,
		"Info":     u.Info,
	})
}

func (h *UserHandler) LoginSMS(ctx *gin.Context) {
	type req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var Req req
	err := ctx.Bind(&Req)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "请求失败",
		})
		return
	}
	ok, err := h.codeSvc.Verify(ctx, bizLogin, Req.Phone, Req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "系统异常",
		})
		zap.L().Error("验证码校验失败",
			// 开发环境能作为debug用
			// 生产环境不能暴露手机号
			zap.String("phone", Req.Phone),
			zap.Error(err))
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{Code: 400, Msg: "验证码不正确"})
		return
	}
	u, err2 := h.svc.FindOrCreate(ctx, Req.Phone)
	if err2 != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "查找创建失败",
		})
		return
	}
	h.SetLoginToken(ctx, u.Id)
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
	})
}

func (h *UserHandler) LoginSMSCode(ctx *gin.Context) {
	type req struct {
		Phone string `json:"phone"`
	}
	var Req req
	err := ctx.Bind(&Req)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "请求失败",
		})
		return
	}
	isPhoneTrue, err := h.phoneRexExp.MatchString(Req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 501, Msg: "系统错误"})
		return
	}
	if !isPhoneTrue {
		ctx.JSON(http.StatusOK, Result{Code: 400, Msg: "电话格式错误"})
		return
	}

	err1 := h.codeSvc.Send(ctx, bizLogin, Req.Phone)
	switch err1 {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case Service.ErrCodeSendTooMany:
		//少数可以接受
		//频繁触发是有人异常触发
		zap.L().Warn("短信发送太频繁", zap.String("phone", Req.Phone))
		ctx.JSON(http.StatusOK, Result{Code: 400, Msg: "短信发送太频繁"})
	default:
		ctx.JSON(http.StatusOK, Result{Code: 501, Msg: "系统错误"})
	}
}

func (h *UserHandler) RefreshToken(ctx *gin.Context) {
	tokenStr := h.ExtractToken(ctx)
	var rc ijwt.RefreshClaims
	token, err := jwt.ParseWithClaims(tokenStr, &rc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.RcJwtKey, nil
	})
	//解析失败，401,未授权
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// jwt没承诺非法就返回错误，加入校验保底
	if token == nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	//校验ssid
	err = h.CheckSession(ctx, rc.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	h.SetJWTToken(ctx, rc.Uid, rc.Ssid)
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "OK",
		"code": "200",
	})
}

func (h *UserHandler) Logout(ctx *gin.Context) {
	err := h.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "已登出",
	})
}
