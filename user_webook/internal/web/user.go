package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"

	"go_work/user_webook/internal/domain"
	"go_work/user_webook/internal/service"
)

const (
	nickNameRegexPattern  = "^[\\u4E00-\\u9FA5A-Za-z0-9]{2,8}$"
	emailRegexPattern     = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	birthRegexPattern     = "^\\d{4}-\\d{1,2}-\\d{1,2}$"
	passwordRegexPattern  = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	introduceRegexPattern = "^[\\u4E00-\\u9FA5A-Za-z0-9]{2,30}$"
)

type UserHandler struct {
	uservice        *service.UserService
	emailRegexp     *regexp2.Regexp
	passwordRegexp  *regexp2.Regexp
	nickenameRegexp *regexp2.Regexp
	birthRegexp     *regexp2.Regexp
	introduceRegexp *regexp2.Regexp
}

func NewUserHandler(uservice *service.UserService) *UserHandler {
	return &UserHandler{
		uservice:        uservice,
		emailRegexp:     regexp2.MustCompile(emailRegexPattern, regexp2.None),
		passwordRegexp:  regexp2.MustCompile(passwordRegexPattern, regexp2.None),
		nickenameRegexp: regexp2.MustCompile(nickNameRegexPattern, regexp2.None),
		birthRegexp:     regexp2.MustCompile(birthRegexPattern, regexp2.None),
		introduceRegexp: regexp2.MustCompile(introduceRegexPattern, regexp2.None),
	}
}

func (c *UserHandler) RegisteRoute(server *gin.Engine) {

	ug := server.Group("/users")

	ug.POST("/signup", c.signUp)
	//ug.POST("/login", c.login)
	ug.POST("/login", c.loginJWT)
	ug.POST("/edit", c.edit)
	ug.POST("/profile", c.profile)
	//ug.POST("/profile", c.profileJWT)
	ug.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "success")
	})
}

func (c *UserHandler) signUp(ctx *gin.Context) {
	type signUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req signUpReq

	if err := ctx.Bind(&req); err != nil {
		return
	}

	isEmail, err := c.emailRegexp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	fmt.Println(req.Email)
	fmt.Println(isEmail)

	if !isEmail {
		ctx.String(http.StatusOK, "你的邮箱格式不对")
		return
	}

	isPassword, err := c.passwordRegexp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	fmt.Println(req.Password, isPassword)
	if !isPassword {
		ctx.String(http.StatusOK, "密码必须包含数字、特殊字符，并且长度不能小于 8 位")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入的密码不相同")
		return
	}

	err = c.uservice.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "hello, 注册成功")
}

func (c *UserHandler) login(ctx *gin.Context) {
	type signUpReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req signUpReq

	if err := ctx.Bind(&req); err != nil {
		return
	}
	fmt.Println(req.Email, req.Password)
	user, err := c.uservice.Login(ctx, req.Email, req.Password)

	if err == service.ErrInvalidEmailOrPassword {
		ctx.String(http.StatusOK, "邮箱或密码不正确")
		return
	}

	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
	}
	fmt.Println(user.Id, user.Email)
	session := sessions.Default(ctx)
	session.Set("user_id", user.Id)
	session.Save()
	ctx.String(http.StatusOK, "登录成功")
	return
}

type UserClaims struct {
	jwt.RegisteredClaims
	// 声明你自己的要放进去 token 里面的数据
	Uid int64
	// 自己随便加
	UserAgent string
}

func (c *UserHandler) loginJWT(ctx *gin.Context) {
	type signUpReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req signUpReq

	if err := ctx.Bind(&req); err != nil {
		return
	}
	fmt.Println(req.Email, req.Password)
	user, err := c.uservice.Login(ctx, req.Email, req.Password)

	if err == service.ErrInvalidEmailOrPassword {
		ctx.String(http.StatusOK, "邮箱或密码不正确")
		return
	}

	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
	}

	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		Uid:       user.Id,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
	}
	ctx.Header("x-jwt-token", tokenStr)

	ctx.String(http.StatusOK, "登录成功")
	return
}

func (c *UserHandler) edit(ctx *gin.Context) {
	type userDetailReq struct {
		NickName     string `json:"nick_name"`
		Birth        string `json:"birth"`
		Introduction string `json:"introduction"`
	}

	var req userDetailReq

	if err := ctx.Bind(&req); err != nil {
		return
	}

	isNickName, err := c.nickenameRegexp.MatchString(req.NickName)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if !isNickName {
		ctx.String(http.StatusOK, "昵称可包含中文，数字，字母，长度2~8")
		return
	}

	isBirth, err := c.birthRegexp.MatchString(req.Birth)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if !isBirth {
		ctx.String(http.StatusOK, "不符合时间格式如`2006-01-01`")
		return
	}

	isIntro, err := c.introduceRegexp.MatchString(req.Introduction)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if !isIntro {
		ctx.String(http.StatusOK, "个人简介可包含中文，数字，字母，长度2~30")
		return
	}

	sess := sessions.Default(ctx)
	id := sess.Get("user_id")

	value, ok := id.(int64)

	if !ok {
		return
	}

	err = c.uservice.EditUserDetail(ctx, value, req.NickName, req.Birth, req.Introduction)

	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "修改成功")
	return
}

func (c *UserHandler) profile(ctx *gin.Context) {
	u, err := c.uservice.Profile(ctx, 1)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "获取用户信息失败")
		return
	}
	str := fmt.Sprintf("user id :%d,email:%s", u.Id, u.Email)
	ctx.String(http.StatusOK, str)
	return
}

func (c *UserHandler) profileJWT(ctx *gin.Context) {
	ca, _ := ctx.Get("claims")
	// 你可以断定，必然有 claims
	//if !ok {
	//	// 你可以考虑监控住这里
	//	ctx.String(http.StatusOK, "系统错误")
	//	return
	//}
	// ok 代表是不是 *UserClaims
	claims, ok := ca.(*UserClaims)
	if !ok {
		// 你可以考虑监控住这里
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	println(claims.Uid)
	ctx.String(http.StatusOK, "你的 profile")
}
