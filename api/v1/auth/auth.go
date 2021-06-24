package auth

import (
	"errors"

	"../../../config"
	"../../../model"
	"../../../service/authService"
	"../../../service/authService/permission"
	"../../response"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// UserForm struct.
type UserForm struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

type BusinessForm struct {
	Token    string      `json:"token"`
	Business interface{} `json:"business"`
}

// Init inits authorization apis
// @Title Auth
// @Description Auth's router group.
func Init(parentRoute *echo.Group) {
	// init admin
	initAdmin(parentRoute)
	// init business
	initBusiness(parentRoute)
	// init driver
	initDriver(parentRoute)
	// init user
	initUser(parentRoute)

	parentRoute.POST("/forgotPassword", forgotPassword)
	parentRoute.POST("/verifyCode", verifyCode)

	route := parentRoute.Group("/me")
	route.Use(middleware.JWT([]byte(config.AuthTokenKey)))

	route.POST("/resetPassword", permission.AuthRequired(resetPassword))
	route.GET("", permission.AuthRequired(readMe))
}

func initAdmin(parentRoute *echo.Group) {
	// admin auth
	parentRoute.POST("/admin/login", loginAdmin)
	parentRoute.POST("/admin/register", registerAdmin)
}

func initBusiness(parentRoute *echo.Group) {
	// business auth
	parentRoute.POST("/business/login", loginBusiness)
	parentRoute.POST("/business/register", registerBusiness)
}

func initDriver(parentRoute *echo.Group) {
	// driver auth
	parentRoute.POST("/driver/login", loginDriver)
	parentRoute.POST("/driver/register", registerDriver)
}

func initUser(parentRoute *echo.Group) {
	// user auth
	parentRoute.POST("/user/login", loginUser)
	parentRoute.POST("/user/register", registerUser)
	parentRoute.POST("/user/facebook/login", loginUserWithFacebook)
}

// @Title forgotPassword
// @Description Forgot Password.
// @Accept  json
// @Produce	json
// @Param   email       form   string   true	"Email."
// @Param   role        form   string   true	"Client role."
// @Success 200 {object} string					"Returns result message"
// @Failure 400 {object} response.BasicResponse "err.email.read"
// @Resource /forgotPassword
// @Router /forgotPassword [post]
func forgotPassword(c echo.Context) error {
	v := &model.Verify{}
	if err := c.Bind(&v); err != nil {
		return response.KnownErrJSON(c, "err.verify.bind", err)
	}
	email := v.Email
	role := v.Role

	// handle forgot password
	if ok := authService.SendVerifyCode(email, role, config.EmailMethod); !ok {
		return response.KnownErrJSON(c, "err.email.read", errors.New("Email is not existed"))
	}
	return response.SuccessJSON(c, "Server has sent email to you. Please check your email and reset password.")
}

// @Title verifyCode
// @Description Verify code.
// @Accept  json
// @Produce	json
// @Param   email       form   string   true	"User Email."
// @Param   role        form   string   true	"Client role."
// @Param   code        form   string   true	"Veryfy code."
// @Success 200 {object} {object}				"Returns token to reset password"
// @Failure 400 {object} response.BasicResponse "err.email.verify"
// @Failure 400 {object} response.BasicResponse "err.user.read"
// @Resource /verifyCode
// @Router /verifyCode [post]
func verifyCode(c echo.Context) error {
	v := &model.Verify{}
	if err := c.Bind(&v); err != nil {
		return response.KnownErrJSON(c, "err.verify.bind", err)
	}
	email := v.Email
	code := v.Code
	role := v.Role

	// check email with verify code
	objid, result, err := authService.CheckVerifyCode(email, code, role)
	if err != nil {
		return response.KnownErrJSON(c, "err.email.verify", err)
	}

	// Generate encoded token and send it as response.
	t, err := permission.GenerateToken(objid, role)
	if err != nil {
		return response.KnownErrJSON(c, "err.auth.token", err)
	}

	return response.SuccessInterface(c, UserForm{t, result})
}

func resetPassword(c echo.Context) error {
	pwd := &model.Password{}
	if err := c.Bind(&pwd); err != nil {
		return response.KnownErrJSON(c, "err.password.bind", err)
	}

	objid, role := permission.InfoFromToken(c)

	result, err := authService.ResetPassword(objid, role, pwd)
	if err != nil {
		return response.KnownErrJSON(c, "err.password.read", err)
	}

	return response.SuccessInterface(c, result)
}

// @Title readMe
// @Description Read self profile information.
// @Accept  json
// @Produce	json
// @Success 200 {object} {object}				"Returns self profile information"
// @Failure 400 {object} response.BasicResponse "err.profile.read"
// @Resource /me
// @Router /me [get]
func readMe(c echo.Context) error {
	objid, role := permission.InfoFromToken(c)
	// read user by objid from token
	result, err := authService.ReadMe(objid, role)
	if err != nil {
		return response.KnownErrJSON(c, "err.profile.read", err)
	}

	return response.SuccessInterface(c, result)
}
