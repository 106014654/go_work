package web

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go_work/user_webook/internal/domain"
	"go_work/user_webook/internal/service"
	"net/http"
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
	ug.POST("/login", c.login)
	ug.POST("/edit", c.edit)
	ug.POST("/profile", c.profile)
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
		ctx.String(http.StatusOK, "系统异常")
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
		ctx.String(http.StatusOK, "系统错误")
	}
	fmt.Println(user.Id, user.Email)
	session := sessions.Default(ctx)
	session.Set("user_id", user.Id)
	session.Save()
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
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "修改成功")
	return
}

func (c *UserHandler) profile(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	id := sess.Get("user_id")

	value, ok := id.(int64)

	if !ok {
		return
	}

	user, err := c.uservice.GetUserInfo(ctx, value)

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if user.NickName == "" && user.Birth == "" && user.Introduction == "" {
		ctx.String(http.StatusOK, "您还未填写任何个人信息")
		return
	}

	userInfo := fmt.Sprintf("user info :昵称:%s, 生日:%s, 个人简介:%s", user.NickName, user.Birth, user.Introduction)

	ctx.String(http.StatusOK, userInfo)
	return
}