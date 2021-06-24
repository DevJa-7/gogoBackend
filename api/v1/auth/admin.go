package auth

import (
	"errors"

	"../../../config"
	"../../../model"
	"../../../service/authService"
	"../../../service/authService/adminService"
	"../../../service/authService/permission"
	"../../../util/crypto"
	"../../response"

	"github.com/labstack/echo"
)

// @Title loginAdmin
// @Description Login a admin.
// @Accept  json
// @Produce	json
// @Param   email       form   string   true	"Admin Email."
// @Param   password	form   string 	true	"Admin Password."
// @Success 200 {object} UserForm 				"Returns login admin"
// @Failure 400 {object} response.BasicResponse "err.admin.bind"
// @Failure 400 {object} response.BasicResponse "err.admin.incorrect"
// @Failure 400 {object} response.BasicResponse "err.admin.token"
// @Resource /admin/login
// @Router /admin/login [post]
func loginAdmin(c echo.Context) error {
	admin := &model.Admin{}
	if err := c.Bind(&admin); err != nil {
		return response.KnownErrJSON(c, "err.admin.bind", err)
	}

	// generate hash password
	admin.Password = crypto.GenerateHash(admin.Password)
	// check admin crediential
	admin, err := adminService.LoginByInfo(admin)
	if err != nil {
		return response.KnownErrJSON(c, "err.admin.incorrect", errors.New("Incorrect email or password"))
	}
	// check verify status
	if !admin.Verify.IsVerified {
		return response.KnownErrJSON(c, "err.admin.verify", errors.New("You are not verifed yet"))
	}
	// generate encoded token and send it as response.
	t, err := permission.GenerateToken(admin.ID, config.RoleAdmin)
	if err != nil {
		return response.KnownErrJSON(c, "err.admin.token",
			errors.New("Something went wrong. Please check token creating"))
	}

	// retreive by public admin
	publicAdmin := &model.PublicAdmin{Admin: admin}
	return response.SuccessInterface(c, UserForm{t, publicAdmin})
}

// @Title registerAdmin
// @Description Register a admin.
// @Accept  json
// @Produce	json
// @Param   email       form   string   true	"Admin Email."
// @Param   password	form   string 	true	"Admin Password."
// @Success 200 {object} UserForm				"Returns registered admin"
// @Failure 400 {object} response.BasicResponse "err.admin.bind"
// @Failure 400 {object} response.BasicResponse "err.admin.exist"
// @Failure 400 {object} response.BasicResponse "err.admin.create"
// @Failure 400 {object} response.BasicResponse "err.admin.token"
// @Resource /admin/register
// @Router /admin/register [post]
func registerAdmin(c echo.Context) error {
	admin := &model.Admin{}
	if err := c.Bind(admin); err != nil {
		return response.KnownErrJSON(c, "err.admin.bind", err)
	}

	// check existed email
	if a, err := adminService.ReadAdminByEmail(admin.Email); err == nil {
		if a.Verify.IsVerified {
			return response.KnownErrJSON(c, "err.admin.exist",
				errors.New("Same email is existed. Please input other email"))
		}
		adminService.DeleteAdmin(a.ID)
	}

	// create admin with registered info
	admin, err := adminService.CreateAdmin(admin)
	if err != nil {
		return response.KnownErrJSON(c, "err.admin.create", err)
	}
	// send to verification email to admin email
	authService.SendVerifyCode(admin.Email, config.RoleAdmin, config.EmailMethod)

	// generate encoded token and send it as response.
	// t, err := permission.GenerateToken(admin.ID, config.RoleAdmin)
	// if err != nil {
	// 	return response.KnownErrJSON(c, "err.admin.token",
	// 		errors.New("Something went wrong. Please check token creating"))
	// }

	// retreive by public admin
	publicAdmin := &model.PublicAdmin{Admin: admin}
	return response.SuccessInterface(c, publicAdmin)
}
