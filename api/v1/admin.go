package v1

import (
	"errors"
	"strconv"

	"../../config"
	"../../model"
	"../../service/authService/adminService"
	"../../service/authService/permission"
	"../../service/roleService"
	"../response"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2/bson"
)

// InitAdmins inits admin CRUD apis
// @Title Admins
// @Description Admins's router group.
func InitAdmins(parentRoute *echo.Group) {
	route := parentRoute.Group("/admins")
	route.Use(middleware.JWT([]byte(config.AuthTokenKey)))

	route.POST("", permission.AuthRequired(createAdmin))
	route.GET("/:id", permission.AuthRequired(readAdmin))
	route.PUT("/:id", permission.AuthRequired(updateAdmin))
	route.DELETE("/:id", permission.AuthRequired(deleteAdmin))

	route.GET("", permission.AuthRequired(readAdmins))
	route.GET("/url/groups", permission.AuthRequired(readAdminURLGroups))

	adminService.InitService()
}

// @Title createAdmin
// @Description Create a admin.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Param   email       	form   	string  true	"Admin Email."
// @Param   password		form   	string 	true	"Admin Password."
// @Success 200 {object} model.PublicAdmin 		"Returns created admin"
// @Failure 400 {object} response.BasicResponse "err.admin.bind"
// @Failure 400 {object} response.BasicResponse "err.admin.create"
// @Resource /admins
// @Router /admins [post]
func createAdmin(c echo.Context) error {
	admin := &model.Admin{}
	if err := c.Bind(admin); err != nil {
		return response.KnownErrJSON(c, "err.admin.bind", err)
	}

	// create admin
	admin, err := adminService.CreateAdmin(admin)
	if err != nil {
		return response.KnownErrJSON(c, "err.admin.create", err)
	}

	publicAdmin := &model.PublicAdmin{Admin: admin}
	return response.SuccessInterface(c, publicAdmin)
}

// @Title readAdmin
// @Description Read a admin.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"Admin ID."
// @Success 200 {object} model.PublicAdmin 		"Returns read admin"
// @Failure 400 {object} response.BasicResponse "err.admin.bind"
// @Failure 400 {object} response.BasicResponse "err.admin.read"
// @Resource /admins
// @Router /admins/{id} [get]
func readAdmin(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.admin.bind", errors.New("Retreived object id is invalid"))
	}

	objid := bson.ObjectIdHex(c.Param("id"))
	admin, err := adminService.ReadAdmin(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.admin.read", err)
	}

	publicAdmin := &model.PublicAdmin{Admin: admin}
	return response.SuccessInterface(c, publicAdmin)
}

// @Title updateAdmin
// @Description Update a admin.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Param   id				path   	string  true	"Admin ID."
// @Param   avatar      	form   	string  true	"Admin Avatar"
// @Param   firstname		form   	string  true	"Admin Firstname"
// @Param   lastname		form   	string  true	"Admin Lastname"
// @Param   email	    	form   	string  true	"Admin Email"
// @Param   birth      		form   	Time   	true	"Admin Birth"
// @Success 200 {object} model.PublicAdmin 		"Returns read admin"
// @Failure 400 {object} response.BasicResponse "err.admin.bind"
// @Failure 400 {object} response.BasicResponse "err.admin.read"
// @Resource /admins
// @Router /admins/{id} [put]
func updateAdmin(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.admin.bind", errors.New("Retreived object id is invalid"))
	}

	admin := &model.Admin{}
	if err := c.Bind(admin); err != nil {
		return response.KnownErrJSON(c, "err.admin.bind", err)
	}

	objid := bson.ObjectIdHex(c.Param("id"))
	admin, err := adminService.UpdateAdmin(objid, admin)
	if err != nil {
		return response.KnownErrJSON(c, "err.admin.read", err)
	}

	publicAdmin := &model.PublicAdmin{Admin: admin}
	return response.SuccessInterface(c, publicAdmin)
}

// @Title deleteAdmin
// @Description Delete a admin.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"Admin ID."
// @Success 200 {object} response.BasicResponse "Admin is deleted correctly."
// @Failure 400 {object} response.BasicResponse "err.admin.bind"
// @Failure 400 {object} response.BasicResponse "err.admin.delete"
// @Resource /admins
// @Router /admins/{id} [delete]
func deleteAdmin(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return response.KnownErrJSON(c, "err.admin.bind", errors.New("Retreived object id is invalid"))
	}

	// delete admin with object id
	objid := bson.ObjectIdHex(c.Param("id"))
	err := adminService.DeleteAdmin(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.admin.delete", err)
	}

	return response.SuccessJSON(c, "Admin is deleted correctly.")
}

// @Title readAdmins
// @Description Read admins with parameters.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string 	true	"Bearer {token}"
// @Param   query			form   	string  false	"Will search string."
// @Param   offset			form    int		false	"Offset for pagination."
// @Param   count 			form    int		false	"Count that will show per page."
// @Param   field			form    string  false	"Sort field."
// @Param   sort			form    int		false	"Sort direction. 0:default, 1:Ascending, -1:Descending"
// @Success 200 {object} ListForm 				"Admin is deleted correctly."
// @Failure 400 {object} response.BasicResponse "err.admin.read"
// @Resource /admins
// @Router /admins [get]
func readAdmins(c echo.Context) error {
	query := c.FormValue("query")
	offset, _ := strconv.Atoi(c.FormValue("offset"))
	count, _ := strconv.Atoi(c.FormValue("count"))
	field := c.FormValue("field")
	sort, _ := strconv.Atoi(c.FormValue("sort"))

	// read admins with params
	admins, total, err := adminService.ReadAdmins(query, offset, count, field, sort)
	if err != nil {
		return response.KnownErrJSON(c, "err.admin.read", err)
	}

	// retreive by public format
	publicAdmins := []*model.PublicAdmin{}
	for _, admin := range admins {
		publicAdmin := &model.PublicAdmin{Admin: admin}
		publicAdmins = append(publicAdmins, publicAdmin)
	}

	return response.SuccessInterface(c, &model.ListForm{total, publicAdmins})
}

// @Title readAdminURLGroups
// @Description Delete a admin.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"Admin ID."
// @Success 200 {object} response.BasicResponse "Admin is deleted correctly."
// @Failure 400 {object} response.BasicResponse "err.admin.bind"
// @Failure 400 {object} response.BasicResponse "err.admin.delete"
// @Resource /admins
// @Router /admins/{id} [delete]
func readAdminURLGroups(c echo.Context) error {
	objid, _ := permission.InfoFromToken(c)
	admin, err := adminService.ReadAdmin(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.admin.read", err)
	}

	urlGroups, err := roleService.ReadURLGroupsWithCode(admin.RoleCode)
	if err != nil {
		return response.KnownErrJSON(c, "err.urlGroup.read", err)
	}

	return response.SuccessInterface(c, urlGroups)
}
