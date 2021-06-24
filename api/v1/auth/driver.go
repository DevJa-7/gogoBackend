package auth

import (
	"errors"

	"../../../config"
	"../../../model"
	"../../../service/authService"
	"../../../service/authService/driverService"
	"../../../service/authService/permission"
	"../../../util/crypto"
	"../../../util/timeHelper"
	"../../response"

	"github.com/labstack/echo"
)

// @Title loginDriver
// @Description Login a driver.
// @Accept  json
// @Produce	json
// @Param   email       form   string   true	"Driver Email."
// @Param   password	form   string 	true	"Driver Password."
// @Success 200 {object} UserForm 			"Returns login driver"
// @Failure 400 {object} response.BasicResponse "err.driver.bind"
// @Failure 400 {object} response.BasicResponse "err.driver.incorrect"
// @Failure 400 {object} response.BasicResponse "err.driver.token"
// @Resource /driver/login
// @Router /driver/login [post]
func loginDriver(c echo.Context) error {
	driver := &model.Driver{}
	if err := c.Bind(&driver); err != nil {
		return response.KnownErrJSON(c, "err.driver.bind", err)
	}

	lastLogin := &model.LastLogin{
		Date: timeHelper.GetCurrentTime(),
		IP:   c.Request().RemoteAddr,
	}
	onesignalPlayerID := driver.OneSignalPlayerID
	platform := driver.Platform

	// generate hash password
	driver.Password = crypto.GenerateHash(driver.Password)
	// check driver crediential
	driver, err := driverService.LoginByInfo(driver)
	if err != nil {
		return response.KnownErrJSON(c, "err.driver.incorrect", errors.New("Incorrect email or password"))
	}

	// generate encoded token and send it as response.
	t, err := permission.GenerateToken(driver.ID, config.RoleDriver)
	if err != nil {
		return response.KnownErrJSON(c, "err.driver.token",
			errors.New("Something went wrong. Please check token creating"))
	}

	// update last login info
	if err := driverService.UpdateLoginInfo(driver.ID, lastLogin, onesignalPlayerID, platform); err != nil {
		return response.KnownErrJSON(c, "err.driver.update", errors.New("Login info update is failed"))
	}

	// retreive by public driver
	publicDriver := &model.PublicDriver{Driver: driver}
	return response.SuccessInterface(c, UserForm{t, publicDriver})
}

// @Title registerDriver
// @Description Register a driver.
// @Accept  json
// @Produce	json
// @Param   firstname	form   string   true	"Driver Firstname."
// @Param   lastname   	form   string   true	"Driver Lastname."
// @Param   email       form   string   true	"Driver Email."
// @Param   password	form   string 	true	"Driver Password."
// @Success 200 {object} UserForm				"Returns registered driver"
// @Failure 400 {object} response.BasicResponse "err.driver.bind"
// @Failure 400 {object} response.BasicResponse "err.driver.exist"
// @Failure 400 {object} response.BasicResponse "err.driver.create"
// @Failure 400 {object} response.BasicResponse "err.driver.token"
// @Resource /driver/register
// @Router /driver/register [post]
func registerDriver(c echo.Context) error {
	driver := &model.Driver{}
	if err := c.Bind(driver); err != nil {
		return response.KnownErrJSON(c, "err.driver.bind", err)
	}

	// check existed email
	if d, err := driverService.ReadDriverByEmail(driver.Email); err == nil {
		if d.Verify.IsVerified {
			return response.KnownErrJSON(c, "err.driver.exist",
				errors.New("Same email is existed. Please input other email"))
		}
		driverService.DeleteDriver(d.ID)
	}
	// create driver with registered info
	driver, err := driverService.CreateDriver(driver)
	if err != nil {
		return response.KnownErrJSON(c, "err.driver.create", err)
	}
	// send to verification email to driver email
	authService.SendVerifyCode(driver.Email, config.RoleDriver, config.TwilloMethod)

	// generate encoded token and send it as response.
	// t, err := permission.GenerateToken(driver.ID, config.RoleDriver)
	// if err != nil {
	// 	return response.KnownErrJSON(c, "err.driver.token",
	// 		errors.New("Something went wrong. Please check token creating"))
	// }

	// retreive by public driver
	publicDriver := &model.PublicDriver{Driver: driver}
	return response.SuccessInterface(c, publicDriver)
}
