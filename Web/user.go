package Handler

import (
	"GinStart/Domain"
	"GinStart/Service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 正则常量
const (
	emailRegex    = "^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$"
	passwordRegex = "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)[a-zA-Z\\d]{8,72}$"
)

func (h *UserHandler) RegisterRoute(server *gin.Engine) {

	user := server.Group("/users")
	user.POST("/signup", h.Signup)
	user.POST("/login", h.Login)
	user.POST("/edit", h.Edit)
	user.GET("/profile", h.Profile)
}

type UserHandler struct {
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	svc            *Service.UserService
}

// NewUserHandler 正则预加载
func NewUserHandler(svc *Service.UserService) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegex, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegex, regexp.None),
		svc:            svc,
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
	case Service.EmailUniqueErr:
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

func (h *UserHandler) Login(c *gin.Context) {

}

func (h *UserHandler) Edit(c *gin.Context) {

}

func (h *UserHandler) Profile(c *gin.Context) {

}
