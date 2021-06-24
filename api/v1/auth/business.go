package auth

import (
	"errors"

	"../../../config"
	"../../../model"
	"../../../service/authService"
	"../../../service/authService/businessService"
	"../../../service/authService/permission"
	"../../response"

	"github.com/labstack/echo"
)

// @Title loginBusiness
// @Description Login a business.
// @Accept  json
// @Produce	json
// @Param   email       form   string   true	"Business Email."
// @Param   password	form   string 	true	"Business Password."
// @Success 200 {object} BusinessForm 			"Returns login business"
// @Failure 400 {object} response.BasicResponse "err.business.bind"
// @Failure 400 {object} response.BasicResponse "err.business.incorrect"
// @Failure 400 {object} response.BasicResponse "err.business.token"
// @Resource /business/login
// @Router /business/login [post]
func loginBusiness(c echo.Context) error {
	business := &model.Business{}
	if err := c.Bind(&business); err != nil {
		return response.KnownErrJSON(c, "err.business.bind", err)
	}

	// generate hash password
	//	business.Password = crypto.GenerateHash(business.Password)
	// check business crediential
	b, err := businessService.LoginByInfo(business)
	if err != nil {
		return response.KnownErrJSON(c, "err.business.incorrect", errors.New("Incorrect email or password"))
	}

	// generate encoded token and send it as response.
	t, err := permission.GenerateToken(b.ID, config.RoleBusiness)
	if err != nil {
		return response.KnownErrJSON(c, "err.business.token",
			errors.New("Something went wrong. Please check token creating"))
	}

	// retreive by public business
	//	publicBusiness := &model.PublicBusiness{Business: business}
	return response.SuccessInterface(c, BusinessForm{t, b})
}

// @Title registerBusiness
// @Description Register a business.
// @Accept  json
// @Produce	json
// @Param   email       form   string   true	"Business Email."
// @Param   password	form   string 	true	"Business Password."
// @Success 200 {object} BusinessForm			"Returns registered business"
// @Failure 400 {object} response.BasicResponse "err.business.bind"
// @Failure 400 {object} response.BasicResponse "err.business.exist"
// @Failure 400 {object} response.BasicResponse "err.business.create"
// @Failure 400 {object} response.BasicResponse "err.business.token"
// @Resource /business/register
// @Router /business/register [post]
func registerBusiness(c echo.Context) error {
	business := &model.Business{}
	if err := c.Bind(business); err != nil {
		return response.KnownErrJSON(c, "err.business.bind", err)
	}

	// check existed email
	if _, err := businessService.ReadBusinessByEmail(business.Email); err == nil {
		return response.KnownErrJSON(c, "err.business.exist",
			errors.New("Same email is existed. Please input other email"))
	}

	// create business with registered info
	business, err := businessService.CreateBusiness(business)
	if err != nil {
		return response.KnownErrJSON(c, "err.business.create", err)
	}

	// send to verification email to business email
	authService.SendVerifyCode(business.Email, config.RoleBusiness, config.EmailMethod)

	// generate encoded token and send it as response.
	// t, err := permission.GenerateToken(business.ID, config.RoleBusiness)
	// if err != nil {
	// 	return response.KnownErrJSON(c, "err.business.token",
	// 		errors.New("Something went wrong. Please check token creating"))
	// }

	// retreive by public business
	// publicBusiness := &model.PublicBusiness{Business: business}
	return response.SuccessInterface(c, business)
}
