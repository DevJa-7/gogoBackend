package v1

import (
	"fmt"
	"strconv"

	"../../api/response"
	"../../config"
	"../../model"
	"../../service/authService/permission"
	"../../service/roleService"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2/bson"
)

// InitRole inits role CRUD apis
// @Title Roles
// @Description Roles's router group.
func InitRole(parentRoute *echo.Group) {
	route := parentRoute.Group("/roles")
	route.Use(middleware.JWT([]byte(config.AuthTokenKey)))

	route.POST("", permission.AuthRequired(createRole))
	route.GET("/:id", permission.AuthRequired(readRole))
	route.PUT("/:id", permission.AuthRequired(updateRole))
	route.DELETE("/:id", permission.AuthRequired(deleteRole))

	route.GET("", permission.AuthRequired(readRoles))

	roleService.InitService()
}

// @Title createRole
// @Description Create a role.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string  true	"Bearer {token}"
// @Param   name       	form   	string  true	"Role name."
// @Param   code		form   	string 	true	"Role code."
// @Param   urlGroup	form   	array 	true	"URL groups for role."
// @Success 200 {object} model.Role             "Returns created role"
// @Failure 400 {object} response.BasicResponse "err.role.bind"
// @Failure 400 {object} response.BasicResponse "err.role.create"
// @Resource /roles
// @Router /roles [post]
func createRole(c echo.Context) error {
	role := &model.Role{}
	if err := c.Bind(role); err != nil {
		return response.KnownErrJSON(c, "err.role.bind", err)
	}
	fmt.Printf("%+v", role)
	// Create role
	role, err := roleService.CreateRole(role)
	if err != nil {
		return response.KnownErrJSON(c, "err.role.create", err)
	}
	return response.SuccessInterface(c, role)
}

// @Title readRole
// @Description Read a role.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"Role ID."
// @Success 200 {object} model.Role 		"Returns read role"
// @Failure 400 {object} response.BasicResponse "err.role.bind"
// @Failure 400 {object} response.BasicResponse "err.role.read"
// @Resource /roles
// @Router /roles/{id} [get]
func readRole(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	// Retrieve rider by id
	role, err := roleService.ReadRole(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.role.read", err)
	}
	return response.SuccessInterface(c, role)
}

// @Title updateRole
// @Description Update role of users.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"Role ID."
// @Param   name			form   	string  true	"Role Name"
// @Param   urlGroup      	form   	string  true	"Role urlGroup"
// @Success 200 {object} model.Role	 		"Returns updated role"
// @Failure 400 {object} response.BasicResponse "err.role.bind"
// @Failure 400 {object} response.BasicResponse "err.role.update"
// @Resource /roles
// @Router /roles/{id} [put]
func updateRole(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))
	role := &model.Role{}
	if err := c.Bind(role); err != nil {
		return response.KnownErrJSON(c, "err.role.bind", err)
	}

	// Update role
	role, err := roleService.UpdateRole(objid, role)
	if err != nil {
		return response.KnownErrJSON(c, "err.role.update", err)
	}
	return response.SuccessInterface(c, role)
}

// @Title deleteRole
// @Description Delete a role.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   id				path   	string  true	"User ID."
// @Success 200 {object} response.BasicResponse "User is deleted correctly."
// @Failure 400 {object} response.BasicResponse "err.user.bind"
// @Failure 400 {object} response.BasicResponse "err.user.delete"
// @Resource /users
// @Router /users/{id} [delete]
func deleteRole(c echo.Context) error {
	objid := bson.ObjectIdHex(c.Param("id"))

	// Remove rider with object id
	err := roleService.DeleteRole(objid)
	if err != nil {
		return response.KnownErrJSON(c, "err.role.delete", err)
	}
	return response.SuccessJSON(c, "Role is deleted correctly.")
}

// @Title readRoles
// @Description Read a roles.
// @Accept  json
// @Produce	json
// @Param   Authorization	header 	string	true	"Bearer {token}"
// @Param   query			form   	string  false	"Will search string."
// @Param   offset			form    int		false	"Offset for pagination."
// @Param   count 			form    int		false	"Count that will show per page."
// @Param   field			form    string  false	"Sort field."
// @Param   sort			form    int		false	"Sort direction. 0:default, 1:Ascending, -1:Descending"
// @Success 200 {object} response.BasicResponse "User is deleted correctly."
// @Failure 400 {object} response.BasicResponse "err.user.bind"
// @Failure 400 {object} response.BasicResponse "err.user.delete"
// @Resource /users
// @Router /users/{id} [delete]
func readRoles(c echo.Context) error {
	query := c.FormValue("query")
	offset, _ := strconv.Atoi(c.FormValue("offset"))
	count, _ := strconv.Atoi(c.FormValue("count"))
	field := c.FormValue("field")
	sort, _ := strconv.Atoi(c.FormValue("sort"))

	// Read roles with query
	roles, total, err := roleService.ReadRoles(query, offset, count, field, sort)
	if err != nil {
		return response.KnownErrJSON(c, "err.role.read", err)
	}

	return response.SuccessInterface(c, &model.ListForm{total, roles})
}
