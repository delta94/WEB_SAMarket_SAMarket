package api

import (
	"fmt"

	"net/http"
	"sam/config"
	"sam/middleware"
	"sam/models"

	"github.com/gin-gonic/gin"
)

func InitAuthRouter(rg *gin.RouterGroup) {
	router := rg.Group("/auth")
	{
		router.POST("/login", login)
		router.GET("/logout", logout)
		router.POST("/register", register)
		if config.Settings.Server.Mode == "debug" {
			router.Use(middleware.TokenAuth)
			router.GET("/checkSession", checkSession)
		}
	}
}

// login godoc
// @Summary 로그인
// @ID login
// @name Login
// @Accept  json
// @Produce  json
// @Param payload body LoginRequest true "로그인 정보"
// @Router /auth/login [post]
// @Success 200 {object} LoginResult
func login(c *gin.Context) {
	var rq *LoginRequest
	err := c.ShouldBindJSON(&rq)
	if err != nil {
		fmt.Println(err)
	}
	user := models.UserStore.GetUserByIDAndPW(rq.LoginID, rq.Password)
	result := &LoginResult{*user, models.ChatStore.GetUnreadCount(user.ID)}
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error"})
	} else {
		middleware.GenerateToken(user, c)
		c.JSON(http.StatusOK, result)
	}
}

// logout godoc
// @Summary 로그아웃
// @Description 로그아웃
// @name Login
// @Accept  json
// @Produce  json
// @Router /auth/logout [get]
// @Success 200 {object} models.User
func logout(c *gin.Context) {
	//토큰 쿠키 expire
	c.SetCookie("token", "", 0, "", "", false, false)
	c.JSON(http.StatusOK, gin.H{"logout ok": "logout ok"})
}

// register godoc
// @Summary 회원가입
// @Description 회원가입, 비밀번호는 sha256 + salt("samarket")
// @Description id, unit 제외하고 보낼것
// @ID register
// @name register
// @Accept  json
// @Produce  json
// @Param json body RegisterRequest true "회원가입 정보"
// @Router /auth/register [post]
// @Success 200 {object} models.User
func register(c *gin.Context) {
	var rq RegisterRequest
	err := c.ShouldBindJSON(&rq)
	if err != nil {
		fmt.Println(err)
	}
	user := models.User{LoginID: rq.LoginID, Name: rq.Name, Password: rq.Password, Mil: rq.Mil, UnitID: rq.UnitID}
	models.UserStore.AddUser(user)
	c.JSON(200, "okay")
}

// checkSession godoc
// @Security ApiKeyAuth
// @Summary 세션 체크
// @Description 세션 체크
// @Router /auth/checkSession [get]
// @Success 200 {object} models.User
func checkSession(c *gin.Context) {
	val, _ := c.Get("user")
	if user, ok := val.(models.User); ok {
		c.JSON(200, user)
	}
}