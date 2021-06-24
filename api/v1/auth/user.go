package auth

import (
	"errors"

	"../../../config"
	"../../../model"
	"../../../service/authService"
	"../../../service/authService/permission"
	"../../../service/authService/userService"
	"../../../util/crypto"
	"../../../util/timeHelper"
	"../../response"

	"github.com/labstack/echo"
)

// @Title loginUser
// @Description Login a user.
// @Accept  json
// @Produce	json
// @Param   email       form   string   true	"User Email."
// @Param   password	form   string 	true	"User Password."
// @Success 200 {object} UserForm 				"Returns login user"
// @Failure 400 {object} response.BasicResponse "err.user.bind"
// @Failure 400 {object} response.BasicResponse "err.user.incorrect"
// @Failure 400 {object} response.BasicResponse "err.user.token"
// @Resource /user/login
// @Router /user/login [post]
func loginUser(c echo.Context) error {
	user := &model.User{}
	if err := c.Bind(&user); err != nil {
		return response.KnownErrJSON(c, "err.user.bind", err)
	}

	lastLogin := &model.LastLogin{
		Date: timeHelper.GetCurrentTime(),
		IP:   c.Request().RemoteAddr,
	}
	onesignalPlayerID := user.OneSignalPlayerID
	platform := user.Platform

	// generate hash password
	user.Password = crypto.GenerateHash(user.Password)
	// check user crediential
	user, err := userService.LoginByInfo(user)
	if err != nil {
		return response.KnownErrJSON(c, "err.user.incorrect", errors.New("Incorrect email or password"))
	}
	// check verify status
	if !user.Verify.IsVerified {
		return response.KnownErrJSON(c, "err.user.verify", errors.New("You are not verifed yet"))
	}

	// generate encoded token and send it as response.
	t, err := permission.GenerateToken(user.ID, config.RoleUser)
	if err != nil {
		return response.KnownErrJSON(c, "err.user.token",
			errors.New("Something went wrong. Please check token creating"))
	}

	// update last login info
	if err := userService.UpdateLoginInfo(user.ID, lastLogin, onesignalPlayerID, platform); err != nil {
		return response.KnownErrJSON(c, "err.user.update", errors.New("Login info update is failed"))
	}

	// retreive by public user
	publicUser := &model.PublicUser{User: user}
	return response.SuccessInterface(c, UserForm{t, publicUser})
}

// @Title registerUser
// @Description Register a user.
// @Accept  json
// @Produce	json
// @Param   firstname	form   string   true	"User Firstname."
// @Param   lastname   	form   string   true	"User Lastname."
// @Param   email       form   string   true	"User Email."
// @Param   password	form   string 	true	"User Password."
// @Success 200 {object} UserForm				"Returns registered user"
// @Failure 400 {object} response.BasicResponse "err.user.bind"
// @Failure 400 {object} response.BasicResponse "err.user.exist"
// @Failure 400 {object} response.BasicResponse "err.user.create"
// @Failure 400 {object} response.BasicResponse "err.user.token"
// @Resource /user/register
// @Router /user/register [post]
func registerUser(c echo.Context) error {
	user := &model.User{}
	if err := c.Bind(user); err != nil {
		return response.KnownErrJSON(c, "err.user.bind", err)
	}
	// check user phone number
	if len(user.Phone) == 0 || len(user.PhoneCode) == 0 {
		return response.KnownErrJSON(c, "err.user.phone", errors.New("Phone number is invalid. Please input phone number again"))
	}
	// check existed email
	if u, err := userService.ReadUserByEmail(user.Email); err == nil {
		if u.Verify.IsVerified {
			return response.KnownErrJSON(c, "err.user.exist",
				errors.New("Same email is existed. Please input other email"))
		}
		userService.DeleteUser(u.ID)
	}
	// create user with registered info
	user, err := userService.CreateUser(user)
	if err != nil {
		return response.KnownErrJSON(c, "err.user.create", err)
	}
	// send to verification email to user email
	authService.SendVerifyCode(user.Email, config.RoleUser, config.TwilloMethod)

	// generate encoded token and send it as response.
	// t, err := permission.GenerateToken(user.ID, config.RoleUser)
	// if err != nil {
	// 	return response.KnownErrJSON(c, "err.user.token",
	// 		errors.New("Something went wrong. Please check token creating"))
	// }

	// retreive by public user
	publicUser := &model.PublicUser{User: user}
	return response.SuccessInterface(c, publicUser)
}

func loginUserWithFacebook(c echo.Context) error {
	user := &model.User{}
	if err := c.Bind(user); err != nil {
		return response.KnownErrJSON(c, "err.user.bind", err)
	}
	lastLogin := &model.LastLogin{
		Date: timeHelper.GetCurrentTime(),
		IP:   c.Request().RemoteAddr,
	}
	onesignalPlayerID := user.OneSignalPlayerID
	platform := user.Platform

	// read user with facebook id
	user, err := userService.ReadUserByFacebookID(user.FacebookUserID)
	if err != nil {
		return response.KnownErrJSON(c, "err.user.read", errors.New("You are not registered yet"))
	}
	// check verify status
	if !user.Verify.IsVerified {
		return response.KnownErrJSON(c, "err.user.verify", errors.New("You are not verifed yet"))
	}

	// generate encoded token and send it as response.
	t, err := permission.GenerateToken(user.ID, config.RoleUser)
	if err != nil {
		return response.KnownErrJSON(c, "err.user.token",
			errors.New("Something went wrong. Please check token creating"))
	}

	// update last login info
	if err := userService.UpdateLoginInfo(user.ID, lastLogin, onesignalPlayerID, platform); err != nil {
		return response.KnownErrJSON(c, "err.user.update", errors.New("Login info update is failed"))
	}

	// retreive by public user
	publicUser := &model.PublicUser{User: user}
	return response.SuccessInterface(c, UserForm{t, publicUser})
}
